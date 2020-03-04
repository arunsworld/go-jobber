package jobber

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Master manages remote workers that can perform instructions.
// Master exhibits behavior to perform these instructions and Shutdown when done.
type Master interface {
	PerformRemoteInstruction(instruction []byte) ([]byte, error)
	Shutdown()
}

// NewMaster creates a new Master that manages remote workers
func NewMaster() (Master, error) {
	regenerationCtx, stopRegeneration := context.WithCancel(context.Background())
	master := &master{
		cmd:              os.Args[0],
		regenerationCtx:  regenerationCtx,
		stopRegeneration: stopRegeneration,
	}
	if err := master.ensureCurrentRemoteWorker(); err != nil {
		return nil, fmt.Errorf("error creating Master: %v", err)
	}
	go master.regenerateForever()
	return master, nil
}

// master allows instructions to be run on a remote worker
type master struct {
	cmd string

	mu                  sync.Mutex
	currentRemoteWorker *remoteWorker
	regenerationCtx     context.Context
	stopRegeneration    context.CancelFunc
}

// PerformRemoteInstruction performs the instruction on a remote worker
// This method may be called concurrently by different goroutines
func (master *master) PerformRemoteInstruction(instruction []byte) ([]byte, error) {
	rw, err := master.remoteWorkerForUse()
	if err != nil {
		return []byte{}, fmt.Errorf("no remoteWorker: %v", err)
	}
	return rw.perform(&Instruction{
		Instruction: instruction,
	})
}

// Shutdown prepares for Master shutdown by issuing Shutdown instruction to the current remote worker
func (master *master) Shutdown() {
	master.stopRegeneration()

	master.mu.Lock()
	defer master.mu.Unlock()

	if master.currentRemoteWorker == nil || master.currentRemoteWorker.isWorkerDead() {
		log.Println("Shutdown: no remote worker alive...")
		return
	}
	if _, err := master.currentRemoteWorker.client.Shutdown(context.Background(), &Empty{}); err != nil {
		log.Printf("error response when shutting down remote worker: %v", err)
	}
}

func (master *master) regenerateForever() {
	for {
		select {
		case <-master.regenerationCtx.Done():
			log.Println("stopping regeneration...")
			return
		case <-time.After(time.Second * 5):
			if err := master.ensureCurrentRemoteWorker(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (master *master) remoteWorkerForUse() (*remoteWorker, error) {
	master.mu.Lock()
	defer master.mu.Unlock()

	if err := master.refreshCurrentRemoteWorkerIfNeeded(); err != nil {
		return nil, err
	}

	rw := master.currentRemoteWorker
	master.currentRemoteWorker = nil
	return rw, nil
}

func (master *master) ensureCurrentRemoteWorker() error {
	master.mu.Lock()
	defer master.mu.Unlock()

	if err := master.refreshCurrentRemoteWorkerIfNeeded(); err != nil {
		return err
	}
	return nil
}

// should only be called when protected within a lock since we're mutating master.currentRemoteWorker
func (master *master) refreshCurrentRemoteWorkerIfNeeded() error {
	select {
	case <-master.regenerationCtx.Done():
		return fmt.Errorf("regeneration stopped... cannot refresh remote worker")
	default:
	}

	if master.currentRemoteWorker == nil || master.currentRemoteWorker.isWorkerDead() {
		newRemoteWorker, err := newRemoteWorker(master.cmd)
		switch err != nil {
		case true:
			return fmt.Errorf("failed to generate remote worker: %v", err)
		default:
			master.currentRemoteWorker = newRemoteWorker
		}
	}
	return nil
}

func newRemoteWorker(workerCmd string) (*remoteWorker, error) {
	ports, err := getFreePorts(1)
	if err != nil {
		return nil, fmt.Errorf("in newRemoteWorker could not find a free port: %v", err)
	}
	port := ports[0]
	params := []string{
		"-worker",
		fmt.Sprintf("-port=%d", port),
	}
	cmd := exec.Command(workerCmd, params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("ERROR: Failure on worker : %v\n", err)
		}
	}()
	time.Sleep(time.Second * 1) // worker takes a bit of time to start up
	conn, err := getConnection(context.Background(), port)
	if err != nil {
		return nil, fmt.Errorf("in newRemoteWorker failed at getting a grpc connection: %v", err)
	}
	client, err := getValidatedClient(conn, 3)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to worker despite retries: %v", err)
	}
	remoteW := &remoteWorker{
		conn:   conn,
		client: client,
		port:   port,
	}
	return remoteW, nil
}

// remoteWorker is basically a wrapper around a grpc client that performs an instruction by calling
// the remote endpoint
type remoteWorker struct {
	conn   *grpc.ClientConn    // connection so it can be closed when we're done
	client JobberServiceClient // client that can be used to talk to the worker
	port   int                 // used as an ID
}

func (rw *remoteWorker) isWorkerDead() bool {
	_, err := rw.client.Hello(context.Background(), &Empty{})
	if err != nil {
		log.Printf("remoteWorker received error while heartbeating worker in isWorkerDead: %v", err)
		rw.conn.Close()
		return true
	}
	return false
}

func (rw *remoteWorker) perform(instruction *Instruction) ([]byte, error) {
	defer rw.conn.Close()

	stream, err := rw.client.Perform(context.Background(), instruction)
	if err != nil {
		return []byte{}, fmt.Errorf("remoteWorker got an error performing instruction: %v", err)
	}
	return processStream(stream)
}

func processStream(stream JobberService_PerformClient) ([]byte, error) {
	var result []byte
	var processingErr error
ReceiveLoop:
	for {
		resp, err := stream.Recv()
		switch {
		case err == io.EOF:
			break ReceiveLoop
		case err != nil:
			processingErr = fmt.Errorf("grpc error while processing stream: %v", err)
		default:
			switch resp.Status {
			case Response_StillProcessing:
			case Response_FinishedWithError:
				processingErr = errors.New(resp.Error)
			case Response_FinishedSuccessfully:
				result = resp.Message
			}
		}
	}
	return result, processingErr
}

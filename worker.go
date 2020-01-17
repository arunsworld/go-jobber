package jobber

import (
	context "context"
	fmt "fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arunsworld/nursery"
	grpc "google.golang.org/grpc"
)

// Job implementation contains the work that needs to be done based on provided instruction
type Job interface {
	Perform(ctx context.Context, instruction []byte) ([]byte, error)
}

// JobFunc is a function that implmenets Job so it can be called inline
type JobFunc func(ctx context.Context, instruction []byte) ([]byte, error)

// Perform simply calls JobFunc
func (jf JobFunc) Perform(ctx context.Context, instruction []byte) ([]byte, error) {
	return jf(ctx, instruction)
}

// NewWorker creates a new instance of a Worker that can be made to listen on a port
func NewWorker(job Job) *Worker {
	worker := &Worker{
		grpcServer: grpc.NewServer(),
		job:        job,
		stopReaper: make(chan struct{}, 1),
	}
	RegisterJobberServiceServer(worker.grpcServer, worker)
	return worker
}

// Worker acts as a GRPC server and executes the given Job when instructed
type Worker struct {
	grpcServer *grpc.Server
	job        Job
	port       int // port acts as an ID

	mu         sync.Mutex
	alreadyRun bool

	missedTicksThreshold time.Duration
	heartBeatTicker      *time.Ticker
	missedHeartbeats     int64
	stopReaper           chan struct{}
}

// Listen starts up the GRPC server listening on given port until instructions are received
func (worker *Worker) Listen(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	var tickerDuration time.Duration
	switch worker.missedTicksThreshold.Microseconds() == 0 {
	case true:
		tickerDuration = time.Second * 5 // it means we wait for 50 seconds of no heartbeats before considering worker as abandoned
	default:
		tickerDuration = worker.missedTicksThreshold / 10
	}
	worker.heartBeatTicker = time.NewTicker(tickerDuration)
	worker.port = port
	go worker.reapAbandonedWorker()
	log.Printf("starting grpc server on port %d", port)
	return worker.grpcServer.Serve(lis)
}

// Hello is nothing more than a heartbeater
func (worker *Worker) Hello(context.Context, *Empty) (*Empty, error) {
	atomic.StoreInt64(&worker.missedHeartbeats, 0)
	return &Empty{}, nil
}

// Perform executes the given Job
func (worker *Worker) Perform(instruction *Instruction, stream JobberService_PerformServer) error {
	worker.mu.Lock()
	if worker.alreadyRun {
		worker.mu.Unlock()
		return fmt.Errorf("Job is already running... cannot invoke Perform again")
	}
	worker.alreadyRun = true
	worker.mu.Unlock()

	worker.heartBeatTicker.Stop()
	worker.stopReaper <- struct{}{}

	var result []byte

	responseCtx, reponseCancel := context.WithCancel(context.Background())
	err := nursery.RunConcurrentlyWithContext(stream.Context(),
		// do the job unless context or stream context is completed
		func(ctx context.Context, errCh chan error) {
			jobW, jobWErr := worker.jobWrapper(ctx, worker.job, instruction)
			select {
			case response := <-jobW:
				result = response
			case err := <-jobWErr:
				errCh <- err
			}
			reponseCancel()
		},
		// send a response back to the server every second unless we have an error or we're done processing
		func(ctx context.Context, errCh chan error) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
					err := stream.Send(&Response{
						Status: Response_StillProcessing,
					})
					switch {
					case err == io.EOF:
						log.Printf("[Worker %d] io.EOF encountered while sending to stream", worker.port)
						return
					case err != nil:
						errCh <- fmt.Errorf("problem sending on stream in Perform: %v", err)
						return
					default:
					}
				case <-responseCtx.Done():
					return
				}
			}
		},
	)
	select {
	case <-stream.Context().Done():
		log.Printf("[Worker %d] not sending any final message because stream is already closed...", worker.port)
	default:
		switch err != nil {
		case true:
			log.Printf("[Worker %d] encountered error: %v", worker.port, err)
			err := stream.Send(&Response{
				Status: Response_FinishedWithError,
				Error:  err.Error(),
			})
			if err != nil {
				log.Printf("[Worker %d] error in final (error) send: %v", worker.port, err)
			}
		default:
			err := stream.Send(&Response{
				Status:  Response_FinishedSuccessfully,
				Message: result,
			})
			if err != nil {
				log.Printf("[Worker %d] error in final (success) send: %v", worker.port, err)
			}
		}
	}
	// shut server gracefully after a brief delay once the function returns
	defer func() {
		go func() {
			time.Sleep(time.Millisecond * 500)
			log.Printf("[Worker: %d] is going to stop, it's job is done", worker.port)
			worker.grpcServer.GracefulStop()
		}()
	}()
	return nil
}

// jobWrapper executes the Job in the background and provides completion & error channels to listen on
func (worker *Worker) jobWrapper(ctx context.Context, job Job, instruction *Instruction) (chan []byte, chan error) {
	resultCh, errCh := make(chan []byte, 1), make(chan error, 1)
	go func() {
		result, err := job.Perform(ctx, instruction.Instruction)
		switch err != nil {
		case true:
			errCh <- err
		default:
			resultCh <- result
		}
	}()
	return resultCh, errCh
}

func (worker *Worker) reapAbandonedWorker() {
	for {
		select {
		case <-worker.heartBeatTicker.C:
			switch atomic.LoadInt64(&worker.missedHeartbeats) >= 10 {
			case true:
				log.Printf("worker [port: %d]: we have missed 10 ticks... worker will stop since it's abandoned...", worker.port)
				worker.grpcServer.GracefulStop()
				return
			default:
				atomic.AddInt64(&worker.missedHeartbeats, 1)
			}
		case <-worker.stopReaper:
			return
		}
	}
}

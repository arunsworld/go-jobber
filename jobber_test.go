package jobber

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/arunsworld/nursery"
	"google.golang.org/grpc"
)

func TestMaster(t *testing.T) {
	t.Skip()

	m, err := NewMaster()
	if err != nil {
		log.Fatal(err)
	}
	result, err := m.PerformRemoteInstruction([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Finished successfully:", string(result))
}

func TestRemoteWorker(t *testing.T) {
	t.Skip()

	rw, err := newRemoteWorker(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	result, err := rw.perform(&Instruction{})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Finished successfully:", string(result))
}

func TestWorker(t *testing.T) {
	err := nursery.RunConcurrently(
		// start the worker
		func(ctx context.Context, errCh chan error) {
			w := NewWorker(JobFunc(func(ctx context.Context, instruction []byte) ([]byte, error) {
				log.Println("performing job...")
				time.Sleep(time.Second * 2)
				log.Println("finished performing job...")
				return []byte("response message"), nil
			}))
			if err := w.Listen(9856); err != nil {
				errCh <- err
			}
		},
		// ask it to perform a job with a client
		func(ctx context.Context, errCh chan error) {
			time.Sleep(time.Second)
			conn, err := grpc.Dial(":9856", grpc.WithInsecure())
			if err != nil {
				t.Fatal(err)
			}
			client := NewJobberServiceClient(conn)
			rw := remoteWorker{
				conn:   conn,
				client: client,
				port:   9856,
			}
			if rw.isWorkerDead() {
				errCh <- errors.New("worker is dead")
			}
			resp, err := rw.perform(&Instruction{})
			switch err != nil {
			case true:
				log.Println("Error:", err)
			default:
				log.Println("Finished successfully:", string(resp))
			}
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAbandonedWorkerGetsReaped(t *testing.T) {
	w := NewWorker(JobFunc(func(ctx context.Context, instruction []byte) ([]byte, error) {
		log.Println("performing job...")
		return []byte{}, nil
	}))
	// in a test we don't want to wait long so we'll make the threshold 1 sec
	w.(*worker).missedTicksThreshold = time.Second
	if err := w.Listen(9856); err != nil {
		t.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/arunsworld/go-jobber"
)

func main() {
	isWorker := flag.Bool("worker", false, "indicates whether we're running as a worker")
	port := flag.Int("port", 0, "port number to run on")
	flag.Parse()

	switch *isWorker {
	case true:
		w := jobber.NewWorker(jobber.JobFunc(func(ctx context.Context, instruction []byte) ([]byte, error) {
			time.Sleep(time.Second * 5)
			return []byte("job done"), nil
		}))
		if err := w.Listen(*port); err != nil {
			log.Fatal(err)
		}
	default:
		m, err := jobber.NewMaster()
		if err != nil {
			log.Fatal(err)
		}
		resp, err := m.PerformRemoteInstruction([]byte{})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("got response:", string(resp))
	}

}

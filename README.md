# go-jobber

Before application servers we had CGI - command line applications that web servers would execute on request. Due to performance issues for web applications they dissappeared but the advantage of CGI was complete process isolation. Once the task is done all resources such as memory are released back to the OS.

Serverless is really a modern form of CGI where upon request an entirely isolated process is kicked-off.

go-jobber implements a CGI/Serverless style of programming. The use-case is for batch jobs that take a lot of resources but once they are complete there is no reason to hold on to them. Go being a grabage collected language will always hold on to memory for quite a while.

The idea is for a master to run as a light-weight server process. It spins off a worker process that is kept ready for a request. When a request comes it is given to the worker which completes the request and dies. Communication between master and worker is via GRPC. The master heartbeats the worker to ensure it's alive and the same heartbeat is also used by the worker to know it's attached to the master. 

If a worker is abandoned it terminates itself. Similarly if a master is unable to reach it's worker it starts up a new one.

# Sample code

```
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
		m.Shutdown()
	}

}
```
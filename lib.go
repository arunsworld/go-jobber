package jobber

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func getFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		l, err := net.Listen("tcp", "[::]:0")
		if err != nil {
			return nil, fmt.Errorf("error in getFreePorts: %v", err)
		}
		l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

func getConnection(ctx context.Context, port int) (*grpc.ClientConn, error) {
	dialCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	url := fmt.Sprintf(":%d", port)
	conn, err := grpc.DialContext(dialCtx, url, grpc.WithInsecure())
	if err != nil {
		log.Printf("problem establishing connection with server: %v", err)
		return nil, err
	}
	return conn, nil
}

func getValidatedClient(conn *grpc.ClientConn, retries int) (JobberServiceClient, error) {
	client := NewJobberServiceClient(conn)
	if err := retryHello(client, retries); err != nil {
		return nil, err
	}
	return client, nil
}

func retryHello(client JobberServiceClient, retries int) error {
	if retries < 1 {
		retries = 1
	}
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.Hello(ctx, &Empty{})
		cancel()
		if err == nil {
			return nil
		}
		log.Printf("client unable to connect (retry %d): %v", i, err)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("client unable to connect to server, all retries exhausted")
}

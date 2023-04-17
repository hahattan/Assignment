package client

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hahattan/grpc/helloworld"
)

type Client struct {
	conn      *grpc.ClientConn
	timeout   time.Duration
	wg        *sync.WaitGroup
	freq      time.Duration
	ctx       context.Context
	metricsCh chan bool
}

func NewClient(ctx context.Context, wg *sync.WaitGroup, timeout time.Duration, freq time.Duration, conn *grpc.ClientConn, metricsCh chan bool) Client {
	return Client{
		conn:      conn,
		timeout:   timeout,
		wg:        wg,
		freq:      freq,
		ctx:       ctx,
		metricsCh: metricsCh,
	}
}

func (c Client) Run(name string) {
	defer c.wg.Done()

	client := pb.NewGreeterClient(c.conn)

	ticker := time.NewTicker(c.freq)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Client call from '%s' stopping ...", name)
			return
		case <-ticker.C:
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				c.metricsCh <- false
				if status.Code(err) == codes.Unavailable {
					log.Printf("[WARN] Server %s temporary unavailable: %v", c.conn.Target(), err)
				} else {
					log.Printf("[ERROR] Could not greet: %v", err)
					cancel()
					return
				}
			} else {
				c.metricsCh <- true
				//log.Printf("Greeting: %s", r.GetMessage())
			}

			cancel()
		}
	}
}

package client

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/hahattan/grpc/helloworld"
)

type mockGreeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (*mockGreeterServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is empty")
	}

	return &helloworld.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	helloworld.RegisterGreeterServer(server, &mockGreeterServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestClient_Run(t *testing.T) {
	validName := "chris"
	emptyName := ""

	tests := []struct {
		testCaseName string
		name         string
		expected     bool
	}{
		{"valid", validName, true},
		{"invalid - empty name", emptyName, false},
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	require.NoError(t, err)
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.testCaseName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}
			ch := make(chan bool, 1)

			wg.Add(1)
			c := NewClient(ctx, wg, time.Second, time.Second, conn, ch)
			go c.Run(tt.name)

			res := <-ch
			cancel()
			wg.Wait()
			assert.Equal(t, res, tt.expected)
		})
	}
}

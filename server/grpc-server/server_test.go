package server

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/hahattan/grpc/helloworld"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &Server{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestServer_SayHello(t *testing.T) {
	validName := "chris"
	emptyName := ""

	tests := []struct {
		testCaseName string
		name         string
		res          *pb.HelloReply
		errCode      codes.Code
		errMsg       string
	}{
		{"valid", validName, &pb.HelloReply{Message: "Hello " + validName}, codes.OK, ""},
		{"invalid - empty name", emptyName, nil, codes.InvalidArgument, "Name is empty"},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	for _, tt := range tests {
		t.Run(tt.testCaseName, func(t *testing.T) {
			request := &pb.HelloRequest{Name: tt.name}
			response, err := client.SayHello(ctx, request)
			if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, s.Code(), tt.errCode)
					assert.Equal(t, s.Message(), tt.errMsg)
				}
			} else if response != nil {
				assert.Equal(t, response.GetMessage(), tt.res.GetMessage())
			}
		})
	}
}

package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hahattan/grpc/helloworld"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is empty")
	}

	log.Printf("Received Greeting from: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

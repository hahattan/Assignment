/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/hahattan/grpc/helloworld"
	"github.com/hahattan/grpc/server/grpc-server"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server.Server{})
	log.Printf("server listening at %v", listener.Addr())

	go func() {
		err := s.Serve(listener)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	cancel()

	log.Println("server stopping...")
	s.Stop()
}

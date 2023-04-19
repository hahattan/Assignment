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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	client "github.com/hahattan/grpc/client/grpc-client"
	"github.com/hahattan/grpc/client/metrics"
)

const (
	defaultName            = "chris"
	defaultTimeout         = time.Second
	defaultMetricsInterval = 30 * time.Second
)

var (
	addr   = flag.String("addr", "localhost:50051", "the address to connect to")
	name   = flag.String("name", defaultName, "Name to greet")
	freq   = flag.Int("freq", 5000, "request frequency in millisecond")
	number = flag.Int("number", 1, "concurrent request number")
)

func main() {
	flag.Parse()

	// Create context that listens for the interrupt signal from the OS.
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ch := make(chan bool, *number)
	reporter := metrics.NewReporter(ctx, &wg, defaultMetricsInterval, ch)
	reporter.Run()

	for i := 1; i <= *number; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Set up a connection to the server.
			conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("[ERROR] did not connect: %v", err)
				return
			}
			defer conn.Close()

			c := client.NewClient(ctx, defaultTimeout, time.Duration(*freq)*time.Millisecond, conn, ch)
			c.Run(fmt.Sprintf("%s#%v", *name, i))
		}(i)
	}

	<-ctx.Done()
	close(ch)

	wg.Wait()
}

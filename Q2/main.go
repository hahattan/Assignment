package main

import (
	"context"
	"github.com/hahattan/assignment/Q2/node"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: main <num_nodes>")
	}

	n := os.Args[1]
	numNodes, err := strconv.Atoi(n)
	if err != nil {
		log.Fatalln("Invalid argument")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ready := make(chan any)
	nodes := make([]*node.Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = new(node.Node)
	}
	for i := 0; i < numNodes; i++ {
		nodes[i] = node.NewNode(i, nodes, ready)
	}

	close(ready)
	<-ctx.Done()
}

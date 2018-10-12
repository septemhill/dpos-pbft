package main

import (
	"context"
	"encoding/gob"
	"os"
	"os/signal"
	"syscall"
)

const (
	numberOfNodes = 20
)

var nodes []*Node

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	gob.Register(Block{})

	//Interrupt signal register
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- struct{}{}
	}()

	//Initialize nodes
	for i := 0; i < numberOfNodes; i++ {
		node := NewNode(ctx, int64(i), false)
		nodes = append(nodes, node)
	}

	//Initialize P2P Network
	for i := 0; i < numberOfNodes; i++ {
		nodes[i].Connect()
	}

	//Start forging
	for i := 0; i < numberOfNodes; i++ {
		nodes[i].Start()
	}

	<-done
	cancel()
}

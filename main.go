package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	numberOfNodes     int   = 21
	numberOfPeers     int   = 14
	listenPort        int64 = 1111
	numberOfDelegates int64 = 21
	slotTimeInterval  int64 = 5
)

func gobInterfaceRegister() {
	gob.Register(Block{})
	gob.Register(Transaction{})
	gob.Register(StageMessage{})
}

func init() {
	gobInterfaceRegister()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sysdone := make(chan struct{}, 1)
	sigCh := make(chan os.Signal, 1)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		sysdone <- struct{}{}
	}()

	nodes := make([]*Node, 0)

	fmt.Println("[Maximum FP Node]", maxFPNode)
	for i := 0; i < numberOfNodes; i++ {
		node := NewNode(ctx, int64(i))
		nodes = append(nodes, node)
	}

	for i := 0; i < numberOfNodes; i++ {
		nodes[i].Connect()
	}

	for i := 0; i < numberOfNodes; i++ {
		go nodes[i].StartForging()
	}

	//for i := 0; i < numberOfNodes; i++ {
	//	msg := BlockMessage(int64(i), *RandomGenerateBlock())
	//	nodes[i].Broadcast(msg)
	//	time.Sleep(time.Second * 2)
	//}

	<-sysdone
	cancel()
}

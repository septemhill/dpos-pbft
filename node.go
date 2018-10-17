package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	Mutex    sync.Mutex
	Id       int64
	Peers    map[int64]*Peer
	PeerIds  []int64
	Listener net.Listener
	Chain    *Blockchain
	Pbft     *Pbft
	LastSlot int64
}

func handleConnection(ctx context.Context, dec *gob.Decoder, node *Node) {
	for {
		var msg Message
		ReceiveMessage(&msg, dec)
		node.ProcessMessage(&msg)
		//fmt.Println("NodeId", node.Id, msg)
		time.Sleep(time.Millisecond * 100)
	}
}

func NewServer(ctx context.Context, node *Node, listenPort int64) net.Listener {
	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(listenPort+node.Id), 10))

	if err != nil {
		log.Println("NewServer Failed")
	}

	go func(ctx context.Context, listener net.Listener) {
		conns := make([]net.Conn, 0)
	END_LISTENER:
		for {
			conn, err := listener.Accept()

			if err != nil {
				log.Println("Accept Failed")
			}

			conns = append(conns, conn)
			dec := gob.NewDecoder(conn)

			go handleConnection(ctx, dec, node)

			select {
			case <-ctx.Done():
				for _, v := range conns {
					v.Close()
				}
				listener.Close()
				fmt.Println("End all connections and listener")
				break END_LISTENER
			default:
			}
		}
	}(ctx, listener)

	return listener
}

func NewNode(ctx context.Context, id int64) *Node {
	node := &Node{
		Id:      id,
		Peers:   make(map[int64]*Peer, 0),
		PeerIds: make([]int64, 0),
	}

	node.Listener = NewServer(ctx, node, listenPort)
	node.Chain = NewBlockchain(node)
	node.Pbft = NewPbft(node)
	fmt.Println("Node ", node.Id, " be created")

	return node
}

func (n *Node) Connect() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numberOfPeers; i++ {
		rand := rand.Int63n(int64(numberOfPeers))
		if rand != n.Id && n.Peers[rand] == nil {
			peer := NewPeer(rand, n.Id, listenPort+rand)
			n.Mutex.Lock()
			n.Peers[rand] = peer
			n.Mutex.Unlock()
			n.PeerIds = append(n.PeerIds, rand)
		}
	}
}

func (n *Node) StartForging() {
	for {
		currentSlot := GetSlotNumber(0)
		lastBlock := n.Chain.GetLastBlock()
		lastSlot := GetSlotNumber(GetTime(lastBlock.GetTimestamp()))

		if currentSlot == lastSlot {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if currentSlot == n.LastSlot {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		delegateId := currentSlot % numberOfDelegates

		if delegateId == n.Id {
			newBlock := n.Chain.CreateBlock()

			//n.Chain.AddBlock(newBlock)
			n.Broadcast(BlockMessage(n.Id, *newBlock))
			n.Pbft.AddBlock(newBlock, GetSlotNumber(GetTime(newBlock.GetTimestamp())))

			fmt.Println("[NODE", n.Id, " NewBlock]", newBlock)
			n.LastSlot = currentSlot
		}

		time.Sleep(time.Second * 1)
	}
}

func (n *Node) Broadcast(msg *Message) {
	for _, peer := range n.Peers {
		SendMessage(msg, peer.ConnEncoder, n.Id)
	}
}

func (n *Node) ProcessMessage(msg *Message) {
	switch msg.Type {
	case MessageTypeInit:
	case MessageTypeBlock:
		block := msg.Body.(Block)
		//fmt.Println("NodeId", n.Id, "receive block message:", block)
		if !n.Chain.HasBlock(block.GetHash()) && n.Chain.ValidateBlock(&block) {
			n.Broadcast(msg)
			n.Pbft.AddBlock(&block, GetSlotNumber(GetTime(block.GetTimestamp())))
		}
	default:
		n.Pbft.ProcessStageMessage(msg)
	}
}

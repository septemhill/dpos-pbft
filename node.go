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

//Node n
type Node struct {
	Mutex    sync.RWMutex
	ID       int64
	Peers    map[int64]*Peer
	PeerIds  []int64
	Listener net.Listener
	Chain    *Blockchain
	Pbft     *Pbft
	LastSlot int64
}

func handleConnection(ctx context.Context, conn net.Conn, dec *gob.Decoder, node *Node) {
	for {
		var msg Message
		ReceiveMessage(&msg, dec)
		node.ProcessMessage(&msg, conn)
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
}

func newServer(ctx context.Context, node *Node, listenPort int64) net.Listener {
	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(listenPort+node.ID), 10))

	if err != nil {
		log.Println("NewServer Failed")
	}

	go func(ctx context.Context, listener net.Listener) {
		conns := make([]net.Conn, 0)
	END_LISTENER:
		for {
			conn, err := listener.Accept()

			if err != nil {
				log.Println("Accept Failed", err)
			}

			conns = append(conns, conn)
			dec := gob.NewDecoder(conn)

			go handleConnection(ctx, conn, dec, node)

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

			//time.Sleep(time.Millisecond * 100)
		}
	}(ctx, listener)

	return listener
}

//NewNode create a new node
func NewNode(ctx context.Context, id int64) *Node {
	node := &Node{
		ID:      id,
		Peers:   make(map[int64]*Peer, 0),
		PeerIds: make([]int64, 0),
	}

	node.Listener = newServer(ctx, node, listenPort)
	node.Chain = NewBlockchain(node)
	node.Pbft = NewPbft(node)
	//fmt.Println("Node ", node.ID, " be created")

	return node
}

//Connect connect to peers
func (n *Node) Connect() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numberOfPeers; i++ {
		rand := rand.Int63n(int64(numberOfPeers))
		n.Mutex.RLock()
		_, ok := n.Peers[rand]
		n.Mutex.RUnlock()
		if rand != n.ID && !ok {
			//peer := NewPeer(rand, n.ID, listenPort+rand)
			peer := NewPeer(rand, listenPort+rand, n)
			n.Mutex.Lock()
			n.Peers[rand] = peer
			n.Mutex.Unlock()
			n.PeerIds = append(n.PeerIds, rand)
		}
	}
}

//StartForging start forging
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

		delegateID := currentSlot % numberOfDelegates

		if delegateID == n.ID {
			currentForger = n.ID
			newBlock := n.Chain.CreateBlock()

			n.Broadcast(BlockMessage(n.ID, *newBlock))
			n.Pbft.AddBlock(newBlock, GetSlotNumber(GetTime(newBlock.GetTimestamp())))

			fmt.Println("[NODE", n.ID, " NewBlock]", newBlock)
			n.LastSlot = currentSlot
		}

		time.Sleep(time.Second * 1)
	}
}

//Broadcast broadcast message to peers
func (n *Node) Broadcast(msg *Message) {
	for _, peer := range n.Peers {
		//fmt.Println("Forger ID:", currentForger, "NodeID", n.ID, "PeerID", peer.ID)
		SendMessage(msg, peer.ConnEncoder, n.ID)
	}
}

//ProcessMessage process message from message
func (n *Node) ProcessMessage(msg *Message, conn net.Conn) {
	switch msg.Type {
	case MessageTypeInit:
		peerID := msg.Body.(int64)
		n.Mutex.RLock()
		_, ok := n.Peers[peerID]
		n.Mutex.RUnlock()
		if !ok {
			n.Mutex.Lock()
			n.Peers[peerID] = &Peer{
				ID:          peerID,
				NodeID:      n.ID,
				Conn:        conn,
				ConnEncoder: gob.NewEncoder(conn),
			}
			n.Mutex.Unlock()
		}
	case MessageTypeBlock:
		block := msg.Body.(Block)
		if !n.Chain.HasBlock(block.GetHash()) && n.Chain.ValidateBlock(&block) {
			n.Broadcast(msg)
			n.Pbft.AddBlock(&block, GetSlotNumber(GetTime(block.GetTimestamp())))
		}
	default:
		n.Pbft.ProcessStageMessage(msg)
	}
}

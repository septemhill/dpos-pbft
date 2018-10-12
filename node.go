package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	numberOfConnectPeers = 5
	listenPort           = 1111
)

type Node struct {
	Id       int64
	LastSlot int64
	IsBad    bool
	PeerIds  []int64
	Peers    map[int64]*Peer
	Chain    *Blockchain
	Server   net.Listener
}

//NewServer create a listen port for peer nodes connection
func NewServer(ctx context.Context, id int64) net.Listener {
	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(listenPort+id), 10))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	go func(ctx context.Context) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				//fmt.Println("NOOOOOoooo")
			} else {

				go func(ctx context.Context, conn net.Conn) {
					for {
						var msg Message
						msg.Deserialize(conn)
						switch msg.Type {
						case MessageTypeInit:
							fmt.Println("Node ", id, " connect to Peer ", msg.Body)
						default:
							fmt.Println("[OTHER MSG TYPE]", msg)
							continue
						}
					}
				}(ctx, conn)

				//				var msg Message
				//				msg.Deserialize(conn)
				//				switch msg.Type {
				//				case MessageTypeInit:
				//					fmt.Println("Node ", id, " connect to Peer ", msg.Body)
				//				default:
				//					fmt.Println("[OTHER MSG TYPE]", msg)
				//					continue
				//				}
			}
		}
	}(ctx)

	return listener
}

//NewNode create a forging node
func NewNode(ctx context.Context, id int64, isBad bool) *Node {
	return &Node{
		Id:      id,
		IsBad:   isBad,
		PeerIds: make([]int64, 0),
		Peers:   make(map[int64]*Peer, 0),
		Chain:   NewBlockchain(),
		Server:  NewServer(ctx, id),
	}
}

//Connect make current node connect to peers
func (node *Node) Connect() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numberOfConnectPeers; i++ {
		rand := int64(rand.Intn(10))
		if rand != node.Id && node.Peers[rand] != nil {
			continue
		}

		peer := NewPeer(rand, listenPort+rand, node.Id)

		node.Peers[rand] = peer
		node.PeerIds = append(node.PeerIds, rand)
	}
}

//PrintBlockchain print current node chain records
func (node *Node) PrintBlockchain() {

}

func (node *Node) CreateBlock() Block {
	lastBlock := node.Chain.GetLastBlock()
	newBlock := Block{
		Version:       1,
		Height:        lastBlock.GetHeight() + 1,
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: lastBlock.GetPrevBlockHash(),
		Forger:        strconv.FormatInt(node.Id, 10),
	}

	newBlock.CalculateBlockHash()
	newBlock.CalculateMerkleHash()

	return newBlock
}

func (node *Node) AddBlock(block Block) {

}

//Start start forging
func (node *Node) Start() {
	//node.Chain.
	currentSlot := GetSlotNumber(0)
	lastBlock := node.Chain.GetLastBlock()
	lastSlot := GetSlotNumber(GetTime(lastBlock.GetTimestamp()))

	if currentSlot == lastSlot {
		//Sleep()
	}
	if currentSlot == node.LastSlot {
		//Sleep()
	}

	delegateId := currentSlot % int64(Delegates)

	if node.Id == delegateId {
		var msg Message
		block := node.CreateBlock()
		node.Broadcast(*msg.BlockMessage(block))
		node.AddBlock(block)
		node.LastSlot = currentSlot
	}
}

//Stop stop forging
func (node *Node) Stop() {

}

//Broadcast send message to peer nodes
func (node *Node) Broadcast(msg Message) {
	for _, peer := range node.Peers {
		//fmt.Println(peer)
		msg.Serialize(peer.Conn)
	}
}

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

//Peer p
type Peer struct {
	ID          int64
	NodeID      int64
	Conn        net.Conn
	ConnEncoder *gob.Encoder
}

func handlePeerConnection(conn net.Conn, dec *gob.Decoder, node *Node) {
	for {
		var msg Message
		rand.Seed(time.Now().UnixNano())
		//fmt.Println("NodeID", node.ID, "get message")
		ReceiveMessage(&msg, dec)
		node.ProcessMessage(&msg, conn)
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
}

//NewPeer create a new peer
func NewPeer(peerID, port int64, node *Node) *Peer {
	conn, err := net.Dial("tcp", ":"+strconv.FormatInt(port, 10))

	if err != nil {
		log.Println("Dial Failed")
	}

	peer := &Peer{
		ID:          peerID,
		NodeID:      node.ID,
		Conn:        conn,
		ConnEncoder: gob.NewEncoder(conn),
	}

	go handlePeerConnection(conn, gob.NewDecoder(conn), node)

	fmt.Println("Node ", node.ID, " connect to peer ", peerID)
	SendMessage(InitMessage(node.ID), peer.ConnEncoder, node.ID)

	return peer
}

//NewPeer create a new peer
//func NewPeer(peerID, nodeID, port int64) *Peer {
//	conn, err := net.Dial("tcp", ":"+strconv.FormatInt(port, 10))
//
//	if err != nil {
//		log.Println("Dial Failed")
//	}
//
//	peer := &Peer{
//		ID:          peerID,
//		NodeID:      nodeID,
//		Conn:        conn,
//		ConnEncoder: gob.NewEncoder(conn),
//	}
//
//	fmt.Println("Node ", nodeID, " connect to peer ", peerID)
//	SendMessage(InitMessage(nodeID), peer.ConnEncoder, nodeID)
//
//	return peer
//}

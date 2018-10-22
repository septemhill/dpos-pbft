package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

//Peer p
type Peer struct {
	ID          int64
	NodeID      int64
	Conn        net.Conn
	ConnEncoder *gob.Encoder
}

//NewPeer create a new peer
func NewPeer(ctx context.Context, peerID, port int64, node *Node) *Peer {
	conn, err := net.Dial("tcp", ":"+strconv.FormatInt(port, 10))

	if err != nil {
		log.Println("Dial Failed", err)
	}

	peer := &Peer{
		ID:          peerID,
		NodeID:      node.ID,
		Conn:        conn,
		ConnEncoder: gob.NewEncoder(conn),
	}

	go handleConnection(ctx, conn, gob.NewDecoder(conn), node)

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

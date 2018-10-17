package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

//Peer p
type Peer struct {
	Id          int64
	NodeId      int64
	Conn        net.Conn
	ConnEncoder *gob.Encoder
}

//NewPeer create a new peer
func NewPeer(peerId, nodeId, port int64) *Peer {
	conn, err := net.Dial("tcp", ":"+strconv.FormatInt(port, 10))

	if err != nil {
		log.Println("Dial Failed")
	}

	peer := &Peer{
		Id:          peerId,
		NodeId:      nodeId,
		Conn:        conn,
		ConnEncoder: gob.NewEncoder(conn),
	}

	fmt.Println("Node ", nodeId, " connect to peer ", peerId)
	SendMessage(InitMessage(nodeId), peer.ConnEncoder, nodeId)

	return peer
}

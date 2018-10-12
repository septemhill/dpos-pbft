package main

import (
	"bytes"
	"encoding/gob"
	"net"
	"strconv"
)

type Peer struct {
	Id     int64
	Port   int64
	NodeId int64
	Conn   net.Conn
}

//NewPeer create a peer node
func NewPeer(peerId, port, nodeId int64) *Peer {
	conn, err := net.Dial("tcp", ":"+strconv.FormatInt(int64(port), 10))

	if err != nil {
		return nil
	}

	peer := &Peer{
		Id:     peerId,
		Port:   port,
		NodeId: nodeId,
		Conn:   conn,
	}

	var msg Message
	msg.InitMessage(nodeId)
	msg.Serialize(conn)

	return peer
}

func (p *Peer) Send(msg Message) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(msg)
	p.Conn.Write(buf.Bytes())
}

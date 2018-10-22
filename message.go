package main

import (
	"encoding/gob"
	"log"
)

//Message Type List
const (
	MessageTypeInit = iota
	MessageTypeBlock
	MessageTypePrepare
	MessageTypeCommit
)

//Message m
type Message struct {
	Type      int
	Body      interface{}
	RoutePath []int64
}

//StageMessage s
type StageMessage struct {
	Height int64
	Hash   string
	Signer string
}

//InitMessage wrap up initialize message
func InitMessage(nodeID int64) *Message {
	m := &Message{Type: MessageTypeInit /*, RoutePath: make([]int64, 0)*/}
	m.Body = nodeID
	//m.RoutePath = append(m.RoutePath, nodeID)
	return m
}

//BlockMessage wrap up block message
func BlockMessage( /*nodeID int64, */ block Block) *Message {
	m := &Message{Type: MessageTypeBlock /*, RoutePath: make([]int64, 0)*/}
	m.Body = block
	//m.RoutePath = append(m.RoutePath, nodeID)
	return m
}

//PrepareMessage wrap up prepare message
func PrepareMessage( /*nodeID int64, */ stage StageMessage) *Message {
	m := &Message{Type: MessageTypePrepare /*, RoutePath: make([]int64, 0)*/}
	m.Body = stage
	//m.RoutePath = append(m.RoutePath, nodeID)
	return m
}

//CommitMessage wrap up commit message
func CommitMessage( /*nodeID int64, */ stage StageMessage) *Message {
	m := &Message{Type: MessageTypeCommit /*, RoutePath: make([]int64, 0)*/}
	m.Body = stage
	//m.RoutePath = append(m.RoutePath, nodeID)
	return m
}

//SendMessage serialize message and send data by socket
func SendMessage(msg *Message, enc *gob.Encoder, nodeID int64) error {
	//time.Sleep(time.Millisecond * 100)
	//Trace routing path (DEBUG)
	//if msg.RoutePath[len(msg.RoutePath)-1] != nodeID {
	//	msg.RoutePath = append(msg.RoutePath, nodeID)
	//}

	err := enc.Encode(msg)
	if err != nil {
		log.Println("[Send Message]", err)
	}

	return err
}

//ReceiveMessage deserialize message and receive data from socket
func ReceiveMessage(msg *Message, dec *gob.Decoder) error {
	err := dec.Decode(msg)
	if err != nil {
		log.Println("[Receive Message]", err)
	}

	return err
}

package main

import (
	"encoding/gob"
	"log"
)

const (
	MessageTypeInit = iota
	MessageTypeBlock
	MessageTypePrepare
	MessageTypeCommit
)

//Message m
type Message struct {
	Type int
	Body interface{}
}

//StageMessage s
type StageMessage struct {
	Height int64
	Hash   string
	Signer string
}

//InitMessage wrap up initialize message
func InitMessage(nodeId int64) *Message {
	m := &Message{Type: MessageTypeInit}
	m.Body = nodeId
	return m
}

//BlockMessage wrap up block message
func BlockMessage(block Block) *Message {
	m := &Message{Type: MessageTypeBlock}
	m.Body = block
	return m
}

//PrepareMessage wrap up prepare message
func PrepareMessage(stage StageMessage) *Message {
	m := &Message{Type: MessageTypePrepare}
	m.Body = stage
	return m
}

//CommitMessage wrap up commit message
func CommitMessage(stage StageMessage) *Message {
	m := &Message{Type: MessageTypeCommit}
	m.Body = stage
	return m
}

//SendMessage serialize message and send data by socket
func SendMessage(msg *Message, enc *gob.Encoder, nodeId int64) error {
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

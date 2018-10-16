package main

import (
	"encoding/gob"
	"log"
	"strconv"
)

const (
	MessageTypeInit = iota
	MessageTypeBlock
	MessageTypePrepare
	MessageTypeCommit
)

type Message struct {
	Type      int
	Body      interface{}
	RoutePath []int64
}

type StageMessage struct {
	Height int64
	Hash   string
	Signer string
}

func InitMessage(nodeId int64) *Message {
	m := &Message{Type: MessageTypeInit, RoutePath: make([]int64, 0)}
	m.Body = "Message from" + strconv.FormatInt(nodeId, 10)
	m.RoutePath = append(m.RoutePath, nodeId)
	return m
}

func BlockMessage(nodeId int64, block Block) *Message {
	m := &Message{Type: MessageTypeBlock, RoutePath: make([]int64, 0)}
	m.Body = block
	m.RoutePath = append(m.RoutePath, nodeId)
	return m
}

func PrepareMessage(nodeId int64, stage StageMessage) *Message {
	m := &Message{Type: MessageTypePrepare, RoutePath: make([]int64, 0)}
	m.Body = stage
	m.RoutePath = append(m.RoutePath, nodeId)
	return m
}

func CommitMessage(nodeId int64, stage StageMessage) *Message {
	m := &Message{Type: MessageTypeCommit, RoutePath: make([]int64, 0)}
	m.Body = stage
	m.RoutePath = append(m.RoutePath, nodeId)
	return m
}

func SendMessage(msg *Message, enc *gob.Encoder, nodeId int64) error {
	//Trace routing path (DEBUG)
	if msg.RoutePath[len(msg.RoutePath)-1] != nodeId {
		msg.RoutePath = append(msg.RoutePath, nodeId)
	}

	err := enc.Encode(msg)
	if err != nil {
		log.Println("[Send Message]", err)
	}

	return err
}

func ReceiveMessage(msg *Message, dec *gob.Decoder) error {
	err := dec.Decode(msg)
	if err != nil {
		log.Println("[Receive Message]", err)
	}

	return err
}

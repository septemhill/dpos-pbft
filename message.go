package main

import (
	"encoding/gob"
	"fmt"
	"io"
)

const (
	MessageTypeInit = iota
	MessageTypeBlock
	MessageTypePrepare
	MessageTypeCommit
)

type Message struct {
	Type int
	Body interface{}
}

type InitMsg struct {
	Id int
}

type BlockMsg struct {
	Block
}

type PrepareMsg struct {
	Height int64
	Hash   string
	Signer int64
}

type CommitMsg struct {
	Height int64
	Hash   string
	Signer int64
}

func (m *Message) Serialize(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(*m); err != nil {
		fmt.Println(err)
		return fmt.Errorf("Message Serialize Failed")
	}
	return nil
}

func (m *Message) Deserialize(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(m); err != nil {
		return fmt.Errorf("Message Deserialize Failed")
	}
	return nil
}

func (m *Message) InitMessage(id int64) *Message {
	m.Type = MessageTypeInit
	m.Body = id
	return m
}

func (m *Message) BlockMessage(block Block) *Message {
	m.Type = MessageTypeBlock
	m.Body = block
	return m
}

func (m *Message) PrepareMessage() *Message {
	m.Type = MessageTypePrepare
	//m.Body =
	return m
}

func (m *Message) CommitMessage() *Message {
	m.Type = MessageTypeCommit
	//m.Body =
	return m
}

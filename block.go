package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
)

type Block struct {
	Version       int
	Height        int64
	Timestamp     int64
	Forger        string
	PrevBlockHash string
	Hash          string
	MerkleRoot    string
	Transtions    []Transaction
}

func (b *Block) Serialize(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("Block Serialize Failed")
	}
	return nil
}

func (b *Block) Deserialize(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(b); err != nil {
		return fmt.Errorf("Block Deserialize Failed")
	}
	return nil
}

func (b *Block) AddTransaction() {}

func (b *Block) GetPrevBlockHash() string {
	return b.PrevBlockHash
}

func (b *Block) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *Block) GetHeight() int64 {
	return b.Height
}

func (b *Block) GetHash() string {
	return b.Hash
}

func (b *Block) CalculateMerkleHash() {
	b.MerkleRoot = "SeptemMerkleHash"
	//return "SeptemMerkleHash"
}

func (b *Block) CalculateBlockHash() {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Version)
	enc.Encode(b.Height)
	enc.Encode(b.Timestamp)
	enc.Encode(b.Forger)
	enc.Encode(b.PrevBlockHash)
	enc.Encode(b.Transtions)

	hash := sha256.Sum256(buf.Bytes())
	b.Hash = hex.EncodeToString(hash[:])
	fmt.Println("[CalculateBlockHash]", b.Hash)
}

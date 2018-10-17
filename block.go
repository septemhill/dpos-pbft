package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"time"
)

//Block B
type Block struct {
	Version       int64
	Height        int64
	Timestamp     int64
	Forger        string
	PrevBlockHash string
	Hash          string
	MerkleRoot    string
	Transactions  []Transaction
}

//NewGenesisBlock g
func NewGenesisBlock() *Block {
	b := &Block{
		Version:       1,
		Height:        0,
		Timestamp:     time.Now().Unix(),
		Forger:        "Septem",
		PrevBlockHash: "0000000000000000000000000000000000000000000000000000000000000000",
		MerkleRoot:    "0000000000000000000000000000000000000000000000000000000000000000",
		Transactions:  make([]Transaction, 0),
	}

	b.CalculateHash()

	return b
}

//GetHash return hash of block
func (b *Block) GetHash() string {
	return b.Hash
}

//GetPrevBlockHash return previous block hash of block
func (b *Block) GetPrevBlockHash() string {
	return b.PrevBlockHash
}

//GetHeight return height of block
func (b *Block) GetHeight() int64 {
	return b.Height
}

//GetTransactions return transactions of block
func (b *Block) GetTransactions() []Transaction {
	return b.Transactions
}

//GetTimestamp return timestamp of block
func (b *Block) GetTimestamp() int64 {
	return b.Timestamp
}

//GetForger return forger of block
func (b *Block) GetForger() string {
	return b.Forger
}

//CalculateMerkleRoot calculte block merkle root hash value
func (b *Block) CalculateMerkleRoot() {

}

//CalculateHash calculate block hash value
func (b *Block) CalculateHash() {
	buff := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buff)
	b.CalculateMerkleRoot()

	enc.Encode(b.Version)
	enc.Encode(b.Height)
	enc.Encode(b.Timestamp)
	enc.Encode(b.Forger)
	enc.Encode(b.PrevBlockHash)
	enc.Encode(b.MerkleRoot)
	enc.Encode(b.Transactions)

	hash := sha256.Sum256(buff.Bytes())
	b.Hash = hex.EncodeToString(hash[:])
}

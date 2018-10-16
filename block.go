package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"time"
)

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

func (b *Block) GetHash() string {
	return b.Hash
}

func (b *Block) GetPrevBlockHash() string {
	return b.PrevBlockHash
}

func (b *Block) GetHeight() int64 {
	return b.Height
}

func (b *Block) GetTransactions() []Transaction {
	return b.Transactions
}

func (b *Block) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *Block) CalculateMerkleRoot() {

}

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

//func RandomGenerateBlock() *Block {
//	rand.Seed(time.Now().UnixNano())
//
//	block := &Block{
//		Version:       rand.Int63n(2000),
//		Height:        rand.Int63n(1000),
//		Timestamp:     rand.Int63n(6000),
//		Forger:        strconv.FormatInt(rand.Int63n(9000000000), 10),
//		PrevBlockHash: strconv.FormatInt(rand.Int63n(9000000000), 10),
//		Hash:          strconv.FormatInt(rand.Int63n(9000000000), 10),
//		MerkleRoot:    strconv.FormatInt(rand.Int63n(9000000000), 10),
//		Transactions:  make([]Transaction, 0),
//	}
//
//	return block
//}

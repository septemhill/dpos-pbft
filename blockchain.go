package main

import (
	"time"
)

type Blockchain struct {
	Node   Node
	Blocks []Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []Block{
			{
				Version:       1,
				Height:        0,
				Timestamp:     time.Now().Unix(),
				Forger:        "Septem",
				PrevBlockHash: "1234567890ABCDEFG",
				MerkleRoot:    "HIJKLMNOPQRSTUVXX",
				Hash:          "YZ1234567890ABCDE",
			},
		},
	}
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

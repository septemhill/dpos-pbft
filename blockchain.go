package main

import (
	"strconv"
	"time"
)

type Blockchain struct {
	Node     *Node
	Pbft     *Pbft
	Blocks   []*Block
	BlockMap map[string]struct{}
}

func NewBlockchain(node *Node) *Blockchain {
	genesisBlock := NewGenesisBlock()

	bc := &Blockchain{
		Node:     node,
		Pbft:     node.Pbft,
		Blocks:   make([]*Block, 0),
		BlockMap: make(map[string]struct{}, 0),
	}

	bc.Blocks = append(bc.Blocks, genesisBlock)
	bc.BlockMap[genesisBlock.GetHash()] = struct{}{}

	return bc
}

func (bc *Blockchain) CreateBlock() *Block {
	lastBlock := bc.GetLastBlock()
	b := &Block{
		Version:       1,
		Height:        lastBlock.GetHeight() + 1,
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: lastBlock.GetHash(),
		Forger:        strconv.FormatInt(bc.Node.Id, 10),
		Transactions:  make([]Transaction, 0),
	}

	b.CalculateHash()

	return b
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) HasBlock(hash string) bool {
	_, ok := bc.BlockMap[hash]
	return ok
}

func (bc *Blockchain) ValidateBlock(block *Block) bool {
	lastBlock := bc.GetLastBlock()
	return block.GetHeight() == lastBlock.GetHeight()+1 &&
		block.GetPrevBlockHash() == lastBlock.GetHash()
}

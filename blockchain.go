package main

import (
	"strconv"
	"sync"
	"time"
)

type Blockchain struct {
	Mutex    sync.Mutex
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
	bc.Mutex.Lock()
	bc.BlockMap[genesisBlock.GetHash()] = struct{}{}
	bc.Mutex.Unlock()

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
	bc.Mutex.Lock()
	_, ok := bc.BlockMap[hash]
	bc.Mutex.Unlock()
	return ok
}

func (bc *Blockchain) ValidateBlock(block *Block) bool {
	lastBlock := bc.GetLastBlock()
	return block.GetHeight() == lastBlock.GetHeight()+1 &&
		block.GetPrevBlockHash() == lastBlock.GetHash()
}

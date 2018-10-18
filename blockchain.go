package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

//Blockchain s
type Blockchain struct {
	Mutex    sync.RWMutex
	Node     *Node
	Pbft     *Pbft
	Blocks   []*Block
	BlockMap map[string]struct{}
}

//NewBlockchain create a new blockchain
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

//CreateBlock create a block
func (bc *Blockchain) CreateBlock() *Block {
	lastBlock := bc.GetLastBlock()
	b := &Block{
		Version:       1,
		Height:        lastBlock.GetHeight() + 1,
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: lastBlock.GetHash(),
		Forger:        strconv.FormatInt(bc.Node.ID, 10),
		Transactions:  make([]Transaction, 0),
	}

	b.CalculateHash()

	return b
}

func (bc *Blockchain) printBlockchain() {
	fmt.Printf("NodeId %d ", bc.Node.ID)
	for i := 0; i < len(bc.Blocks); i++ {
		fmt.Printf("%s -> ", bc.Blocks[i].GetForger())
	}
	fmt.Println()
}

//CommitBlock add block into chain
func (bc *Blockchain) CommitBlock(block *Block) {
	bc.Blocks = append(bc.Blocks, block)
	bc.Mutex.Lock()
	bc.BlockMap[block.GetHash()] = struct{}{}
	bc.Mutex.Unlock()
	bc.printBlockchain()
}

//GetLastBlock get last block of the chain
func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

//HasBlock check the block whether in the chain or not
func (bc *Blockchain) HasBlock(hash string) bool {
	bc.Mutex.RLock()
	_, ok := bc.BlockMap[hash]
	bc.Mutex.RUnlock()
	return ok
}

//ValidateBlock valiate  block
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	lastBlock := bc.GetLastBlock()
	return block.GetHeight() == lastBlock.GetHeight()+1 &&
		block.GetPrevBlockHash() == lastBlock.GetHash()
}

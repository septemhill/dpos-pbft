package main

import (
	"fmt"
	"log"
	"strconv"
)

const (
	PBFTStateNone = iota
	PBFTStatePrepare
	PBFTStateCommit
)

type ConsensusInfo struct {
	Height      int64
	Hash        string
	VotesNumber int64
	Votes       map[string]struct{}
}

type Pbft struct {
	State            int64
	Node             *Node
	PendingBlocks    map[string]*Block
	CurrentSlot      int64
	PrepareInfo      *ConsensusInfo
	CommitInfos      []*ConsensusInfo
	PreparehashCache map[string]struct{}
}

func NewConsensusInfo() *ConsensusInfo {
	return &ConsensusInfo{
		Votes: make(map[string]struct{}, 0),
	}
}

func NewPbft(node *Node) *Pbft {
	pbft := &Pbft{
		State:            PBFTStateNone,
		Node:             node,
		PendingBlocks:    make(map[string]*Block, 0),
		CurrentSlot:      0,
		PrepareInfo:      NewConsensusInfo(),
		CommitInfos:      make([]*ConsensusInfo, 0),
		PreparehashCache: make(map[string]struct{}, 0),
	}

	return pbft
}

func (p *Pbft) AddBlock(block *Block, slotNumber int64) {
	hash := block.GetHash()
	p.PendingBlocks[hash] = block

	if slotNumber > p.CurrentSlot {
		p.ClearState()
	}

	if p.State == PBFTStateNone {
		p.CurrentSlot = slotNumber
		p.State = PBFTStatePrepare
		p.PrepareInfo.Height = block.GetHeight()
		p.PrepareInfo.Hash = block.GetHash()
		p.PrepareInfo.VotesNumber = 1
		p.PrepareInfo.Votes[strconv.FormatInt(p.Node.Id, 10)] = struct{}{}

		stageMsg := StageMessage{
			Height: block.GetHeight(),
			Hash:   block.GetHash(),
			Signer: strconv.FormatInt(p.Node.Id, 10),
		}
		p.Node.Broadcast(PrepareMessage(p.Node.Id, stageMsg))
	}
}

func (p *Pbft) ClearState() {

}

func (p *Pbft) handlePrepareMessage(msg StageMessage) {
	cacheKey := fmt.Sprintf("%s:%d:%s", msg.Hash, msg.Height, msg.Signer)
	fmt.Println("[CACHE KEY]", cacheKey)
}

func (p *Pbft) handleCommitMessage(msg StageMessage) {
}

func (p *Pbft) ProcessStageMessage(msg *Message) {
	switch msg.Type {
	case MessageTypePrepare:
		//fmt.Println("[Prepare]", msg.Body)
		p.handlePrepareMessage(msg.Body.(StageMessage))
	case MessageTypeCommit:
	default:
		log.Println("ProcessStageMessage cannot find match message type")
	}
}

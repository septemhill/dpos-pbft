package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

var (
	maxFPNode = int64(math.Floor(float64((numberOfDelegates - 1) / 3)))
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
	CommitInfos      map[string]*ConsensusInfo
	PreparehashCache map[string]struct{}
}

func NewConsensusInfo() *ConsensusInfo {
	return &ConsensusInfo{
		VotesNumber: 1,
		Votes:       make(map[string]struct{}, 0),
	}
}

func NewPbft(node *Node) *Pbft {
	pbft := &Pbft{
		State:            PBFTStateNone,
		Node:             node,
		PendingBlocks:    make(map[string]*Block, 0),
		CurrentSlot:      0,
		PrepareInfo:      NewConsensusInfo(),
		CommitInfos:      make(map[string]*ConsensusInfo, 0),
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

func (p *Pbft) handlePrepareMessage(msg *Message) {
	fmt.Printf("NodeId %d receive prepare message: %s\n", p.Node.Id, msg.Body.(StageMessage).Hash)
	stageMsg := msg.Body.(StageMessage)
	cacheKey := fmt.Sprintf("%s:%d:%s", stageMsg.Hash, stageMsg.Height, stageMsg.Signer)

	if _, ok := p.PreparehashCache[cacheKey]; !ok {
		p.PreparehashCache[cacheKey] = struct{}{}
		p.Node.Broadcast(msg)
	} else {
		return
	}

	_, voted := p.PrepareInfo.Votes[stageMsg.Signer]

	if p.State == PBFTStatePrepare && stageMsg.Hash == p.PrepareInfo.Hash &&
		stageMsg.Height == p.PrepareInfo.Height && !voted {
		p.PrepareInfo.Votes[stageMsg.Signer] = struct{}{}
		p.PrepareInfo.VotesNumber++

		if p.PrepareInfo.VotesNumber > maxFPNode {
			p.State = PBFTStateCommit
			commitInfo := NewConsensusInfo()
			commitInfo.Hash = p.PrepareInfo.Hash
			commitInfo.Height = p.PrepareInfo.Height
			commitInfo.Votes[strconv.FormatInt(p.Node.Id, 10)] = struct{}{}
			p.CommitInfos[commitInfo.Hash] = commitInfo

			stageMsg := StageMessage{
				Height: p.PrepareInfo.Height,
				Hash:   p.PrepareInfo.Hash,
				Signer: strconv.FormatInt(p.Node.Id, 10),
			}

			p.Node.Broadcast(CommitMessage(p.Node.Id, stageMsg))
		}
	}

	//fmt.Println("[CACHE KEY]", cacheKey)
}

func (p *Pbft) handleCommitMessage(msg *Message) {
	fmt.Println("[Commit Message]", msg.Body)
}

func (p *Pbft) ProcessStageMessage(msg *Message) {
	switch msg.Type {
	case MessageTypePrepare:
		p.handlePrepareMessage(msg)
	case MessageTypeCommit:
		p.handleCommitMessage(msg)
	default:
		log.Println("ProcessStageMessage cannot find match message type")
	}
}

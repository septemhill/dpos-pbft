package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

//PBFT Status List
const (
	PBFTStateNone = iota
	PBFTStatePrepare
	PBFTStateCommit
)

//ConsensusInfo c
type ConsensusInfo struct {
	Height      int64
	Hash        string
	VotesNumber int64
	Votes       map[string]struct{}
}

//Pbft p
type Pbft struct {
	MsgLock          sync.RWMutex
	Mutex            sync.RWMutex
	State            int64
	Node             *Node
	PendingBlocks    map[string]*Block
	CurrentSlot      int64
	PrepareInfo      *ConsensusInfo
	CommitInfos      map[string]*ConsensusInfo
	PrepareHashCache map[string]struct{}
	CommitHashCache  map[string]struct{}
	Chain            *Blockchain
}

func newConsensusInfo() *ConsensusInfo {
	return &ConsensusInfo{
		VotesNumber: 1,
		Votes:       make(map[string]struct{}, 0),
	}
}

//NewPbft create a new pbft instance
func NewPbft(node *Node) *Pbft {
	pbft := &Pbft{
		State:            PBFTStateNone,
		Node:             node,
		PendingBlocks:    make(map[string]*Block, 0),
		CurrentSlot:      0,
		PrepareInfo:      newConsensusInfo(),
		CommitInfos:      make(map[string]*ConsensusInfo, 0),
		PrepareHashCache: make(map[string]struct{}, 0),
		CommitHashCache:  make(map[string]struct{}, 0),
		Chain:            node.Chain,
	}

	return pbft
}

//AddBlock made block into prepare state
func (p *Pbft) AddBlock(block *Block, slotNumber int64) {
	hash := block.GetHash()
	p.Mutex.Lock()
	p.PendingBlocks[hash] = block
	p.Mutex.Unlock()

	if slotNumber > p.CurrentSlot {
		p.ClearState()
	}

	if p.State == PBFTStateNone {
		p.CurrentSlot = slotNumber
		p.State = PBFTStatePrepare
		p.PrepareInfo.Height = block.GetHeight()
		p.PrepareInfo.Hash = block.GetHash()
		p.PrepareInfo.VotesNumber = 1
		p.Mutex.Lock()
		p.PrepareInfo.Votes[strconv.FormatInt(p.Node.ID, 10)] = struct{}{}
		p.Mutex.Unlock()
		//fmt.Println("node", p.Node.ID, "change state to prepare")

		stageMsg := StageMessage{
			Height: block.GetHeight(),
			Hash:   block.GetHash(),
			Signer: strconv.FormatInt(p.Node.ID, 10),
		}

		p.Node.Broadcast(PrepareMessage( /*p.Node.ID, */ stageMsg))
	}
}

//ClearState clear state
func (p *Pbft) ClearState() {
	p.State = PBFTStateNone
	p.Mutex.Lock()
	p.PrepareInfo = newConsensusInfo()
	p.CommitInfos = make(map[string]*ConsensusInfo)
	p.PendingBlocks = make(map[string]*Block)
	p.Mutex.Unlock()
}

func (p *Pbft) handlePrepareMessage(msg *Message) {
	stageMsg := msg.Body.(StageMessage)
	cacheKey := fmt.Sprintf("%s:%d:%s", stageMsg.Hash, stageMsg.Height, stageMsg.Signer)

	p.Mutex.RLock()
	_, ok := p.PrepareHashCache[cacheKey]
	p.Mutex.RUnlock()

	if !ok {
		p.Mutex.Lock()
		p.PrepareHashCache[cacheKey] = struct{}{}
		p.Mutex.Unlock()
		p.Node.Broadcast(msg)
	} else {
		return
	}

	p.Mutex.RLock()
	_, voted := p.PrepareInfo.Votes[stageMsg.Signer]
	p.Mutex.RUnlock()

	if p.State == PBFTStatePrepare && stageMsg.Hash == p.PrepareInfo.Hash &&
		stageMsg.Height == p.PrepareInfo.Height && !voted {
		p.Mutex.Lock()
		p.PrepareInfo.Votes[stageMsg.Signer] = struct{}{}
		p.Mutex.Unlock()
		p.PrepareInfo.VotesNumber++
		//fmt.Println("pbft", p.Node.ID, "prepare votes", p.PrepareInfo.VotesNumber)

		if p.PrepareInfo.VotesNumber > maxFPNode {
			//fmt.Println("node", p.Node.ID, "change state to commit")
			p.State = PBFTStateCommit
			commitInfo := newConsensusInfo()
			commitInfo.Hash = p.PrepareInfo.Hash
			commitInfo.Height = p.PrepareInfo.Height
			commitInfo.Votes[strconv.FormatInt(p.Node.ID, 10)] = struct{}{}
			p.Mutex.Lock()
			p.CommitInfos[commitInfo.Hash] = commitInfo
			p.Mutex.Unlock()

			stageMsg := StageMessage{
				Height: p.PrepareInfo.Height,
				Hash:   p.PrepareInfo.Hash,
				Signer: strconv.FormatInt(p.Node.ID, 10),
			}

			p.Node.Broadcast(CommitMessage( /*p.Node.ID, */ stageMsg))
		}
	}
}

func (p *Pbft) handleCommitMessage(msg *Message) {
	stageMsg := msg.Body.(StageMessage)
	cacheKey := fmt.Sprintf("%s:%d:%s", stageMsg.Hash, stageMsg.Height, stageMsg.Signer)

	p.Mutex.RLock()
	_, ok := p.CommitHashCache[cacheKey]
	p.Mutex.RUnlock()

	if !ok {
		p.Mutex.Lock()
		p.CommitHashCache[cacheKey] = struct{}{}
		p.Mutex.Unlock()
		p.Node.Broadcast(msg)
	} else {
		return
	}

	p.Mutex.RLock()
	commitInfo := p.CommitInfos[stageMsg.Hash]
	p.Mutex.RUnlock()

	if commitInfo != nil {
		if _, ok := commitInfo.Votes[stageMsg.Signer]; !ok {
			p.Mutex.Lock()
			commitInfo.Votes[stageMsg.Signer] = struct{}{}
			p.Mutex.Unlock()
			commitInfo.VotesNumber++
			//fmt.Println("pbft", p.Node.ID, "commit votes", commitInfo.VotesNumber)
			if commitInfo.VotesNumber > 2*maxFPNode {
				p.Mutex.RLock()
				_, ok := p.PendingBlocks[stageMsg.Hash]

				if ok {
					p.Chain.CommitBlock(p.PendingBlocks[stageMsg.Hash])
				}
				p.Mutex.RUnlock()
				p.ClearState()
			}
		}
	} else {
		commitInfo := newConsensusInfo()
		commitInfo.Hash = stageMsg.Hash
		commitInfo.Height = stageMsg.Height
		commitInfo.Votes[strconv.FormatInt(p.Node.ID, 10)] = struct{}{}
		p.Mutex.Lock()
		p.CommitInfos[stageMsg.Hash] = commitInfo
		p.Mutex.Unlock()
	}
}

//ProcessStageMessage process prepare and commit message
func (p *Pbft) ProcessStageMessage(msg *Message) {
	switch msg.Type {
	case MessageTypePrepare:
		p.MsgLock.Lock()
		p.handlePrepareMessage(msg)
		p.MsgLock.Unlock()
	case MessageTypeCommit:
		p.MsgLock.Lock()
		p.handleCommitMessage(msg)
		p.MsgLock.Unlock()
	default:
		log.Println("ProcessStageMessage cannot find match message type")
	}
}

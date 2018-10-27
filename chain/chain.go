package chain

import (
	"blockchain-experiements/block"
	"blockchain-experiements/consensus"
	"time"
)

type Blockchain struct {
	Blocks []*block.Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks ) - 1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *block.Block {
	return NewBlock("Genesis block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{NewGenesisBlock()}}
}

func NewBlock(data string, prevBlockHash []byte) *block.Block {
	block := &block.Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0 }

	pow := consensus.NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

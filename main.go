package main

import (
	"blockchain-experiements/chain"
	"fmt"
)

func main() {
	bc := chain.NewBlockchain()

	bc.AddBlock("send 1 from 123 to 231")
	bc.AddBlock("send 1 from 231 to 123")

	for _, block := range bc.Blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}

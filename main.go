package main

import (
	"blockchain-experiements/chain"
	"fmt"
	"math/rand"
)

func main() {
	bc, err := chain.NewBlockchain()
	defer bc.DB.Close()

	count := rand.Intn(10)
	err = bc.AddBlock("send " + string(count) + " from 123 to 231")

	i := bc.Iterator()

	if err != nil {
		fmt.Printf("Error " + string(err.Error()))
		return
	}

	for {
		b, _:= i.Next()

		fmt.Printf("Previous hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

}

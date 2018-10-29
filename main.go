package main

import (
	"blockchain-experiements/chain"
	"blockchain-experiements/tx"
	"fmt"
)

func getBalance(bc *chain.Blockchain, address string) int {
	outs := bc.FindUTXOs(address)

	var balance int

	fmt.Printf("Outs len: %d\n", len(outs))

	for _, out := range outs {
		balance += out.Value
	}

	return balance
}

func main() {
	bc, err := chain.CreateBlockchain("Test")
	defer (func() {
		_ = bc.DB.Close()
	})()

	fmt.Printf("Balance of Test: %d\n", getBalance(bc, "Test"))
	fmt.Printf("Balance of Test2: %d\n", getBalance(bc, "Test2"))

	tr := bc.NewUTXOTransaction("Test", "Test2", 1)
	_ = bc.AddBlock([]*tx.Transaction{tr})

	i := bc.Iterator()

	if err != nil {
		fmt.Printf("Error " + string(err.Error()))
		return
	}

	for {
		b, _:= i.Next()

		fmt.Printf("Previous hash: %x\n", b.PrevBlockHash)
		fmt.Printf("TxHash: %x\n", b.HashTransactions())
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

}

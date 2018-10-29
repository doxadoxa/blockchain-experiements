package tx

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

const CoinbaseReward = 1

type TXInput struct {
	Txid 		[]byte
	Vout 		int
	ScriptSig 	string
}

type TXOutput struct {
	Value 			int
	ScriptPubKey 	string
}

type Transaction struct {
	ID 		[]byte
	Vin 	[]TXInput
	Vout 	[]TXOutput
}

func (t *Transaction) SetID() {
	var prevTxsHashes [][]byte

	for _, in := range t.Vin {
		prevTxsHashes = append(prevTxsHashes, in.Txid)
	}

	hash := sha256.Sum256(bytes.Join(prevTxsHashes, []byte{}))

	t.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.ID) == 0
}

func (in *TXInput) CanUnlockOutputWith(key string) bool {
	return in.ScriptSig == key
}

func (out *TXOutput) CanBeUnlockedWith(key string) bool {
	return out.ScriptPubKey == key
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Rewards to %s", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{CoinbaseReward, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}

	tx.SetID()

	return &tx
}
package block

import (
	"blockchain-experiements/tx"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Timestamp		int64
	Transactions 	[]*tx.Transaction
	PrevBlockHash	[]byte
	Hash			[]byte
	Nonce 			uint
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	return result.Bytes(), err
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block, err
}

func (b *Block) HashTransactions() []byte {
	var txHashes 	[][]byte
	var txHash 		[32]byte

	for _, transaction := range b.Transactions {
		txHashes = append(txHashes, transaction.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
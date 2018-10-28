package block

import (
	"bytes"
	"encoding/gob"
)

type Block struct {
	Timestamp		int64
	Data			[]byte
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
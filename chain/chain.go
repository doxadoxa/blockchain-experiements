package chain

import (
	"blockchain-experiements/block"
	"blockchain-experiements/consensus"
	"github.com/boltdb/bolt"
	"time"
)

const dbFile string = "blocks.dat"
const blocksBucket string = "blocks"
const tailPointer string = "l"

type Blockchain struct {
	tail 	[]byte
	db		*bolt.DB
}

type BlockchainIterator struct {
	currentHash		[]byte
	db				*bolt.DB
}

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	var err error

	err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(tailPointer))

		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serialized, serializeErr := newBlock.Serialize()

		if serializeErr != nil { return nil }

		err = b.Put(newBlock.Hash, serialized)
		err = b.Put([]byte(tailPointer), newBlock.Hash)
		bc.tail = newBlock.Hash

		return nil
	})

	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tail, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() (*block.Block, error) {
	var b *block.Block
	var err error

	err = i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		b, err = block.DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = b.PrevBlockHash

	return b, err
}

func NewGenesisBlock() *block.Block {
	return NewBlock("Genesis block", []byte{})
}

func NewBlockchain() (*Blockchain, error) {
	var tail []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			genesis := NewGenesisBlock()
			bucket, _ := bucket.CreateBucket([]byte(blocksBucket))

			serialized, _ := genesis.Serialize()

			err = bucket.Put(genesis.Hash, serialized)
			err = bucket.Put([]byte(tailPointer), genesis.Hash)
			tail = genesis.Hash
		} else {
			tail = bucket.Get([]byte(tailPointer))
		}

		return nil
	})

	return &Blockchain{tail, db}, err
}

func NewBlock(data string, prevBlockHash []byte) *block.Block {
	b := &block.Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0 }

	pow := consensus.NewProofOfWork(b)
	nonce, hash := pow.Run()

	b.Hash = hash
	b.Nonce = nonce

	return b
}

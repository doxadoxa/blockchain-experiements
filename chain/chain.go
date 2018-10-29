package chain

import (
	"blockchain-experiements/block"
	"blockchain-experiements/consensus"
	t "blockchain-experiements/tx"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const dbFile string = "blocks.dat"
const blocksBucket string = "blocks"
const tailPointer string = "l"
const genesisBlockData string = "Genesis start here"

type Blockchain struct {
	tail []byte
	DB   *bolt.DB
}

type BlockchainIterator struct {
	currentHash		[]byte
	db				*bolt.DB
}

func (bc *Blockchain) AddBlock(transactions []*t.Transaction) error {
	var lastHash []byte
	var err error

	err = bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(tailPointer))

		return nil
	})

	newBlock := NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
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
	bci := &BlockchainIterator{bc.tail, bc.DB}

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

func (blockchain *Blockchain) FindUnspentTransactions(address string) []t.Transaction {
	var unspentTXs []t.Transaction
	spentTXOs := make(map[string][]int)

	i := blockchain.Iterator()

	for {
		bl, _ := i.Next()

		for _, transaction := range bl.Transactions {
			txId := hex.EncodeToString(transaction.ID)

		Outputs:
			for outNum, outTx := range transaction.Vout {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outNum {
							continue Outputs
						}
					}
				}

				if outTx.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *transaction)
				}
			}

			if transaction.IsCoinbase() == false {
				for _, spentIn := range transaction.Vin {
					if spentIn.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(spentIn.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], spentIn.Vout)
					}
				}
			}
		}

		if len(bl.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXOs(address string) []t.TXOutput {
	var unspentTXOs []t.TXOutput
	unspentTXs := bc.FindUnspentTransactions(address)

	for _, transaction := range unspentTXs {
		for _, out := range transaction.Vout {
			if out.CanBeUnlockedWith(address) {
				unspentTXOs = append(unspentTXOs, out)
			}
		}
	}

	return unspentTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

	FindTX:
		for _, tx := range unspentTXs {
			txID := hex.EncodeToString(tx.ID)

			for outNum, out := range tx.Vout {
				if out.CanBeUnlockedWith(address) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outNum)

					if accumulated > amount {
						break FindTX
					}
				}
			}
		}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) NewUTXOTransaction(from string, to string, amount int) *t.Transaction {
	var inputs []t.TXInput
	var outputs []t.TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds.")
	}

	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := t.TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, t.TXOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, t.TXOutput{acc - amount, from})
	}

	transaction := t.Transaction{nil, inputs, outputs}
	transaction.SetID()

	return &transaction
}

func NewGenesisBlock(coinbase *t.Transaction) *block.Block {
	return NewBlock([]*t.Transaction{coinbase}, []byte{})
}

func CreateBlockchain(address string) (*Blockchain, error) {
	var tail []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(blocksBucket))

		if bucket == nil {
			cbtx := t.NewCoinbaseTX(address, genesisBlockData)
			genesis := NewGenesisBlock(cbtx)
			bucket, _ := transaction.CreateBucket([]byte(blocksBucket))

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

func NewBlock(transactions []*t.Transaction, prevBlockHash []byte) *block.Block {
	coinbase := t.NewCoinbaseTX("Test", "")
	transactions = append([]*t.Transaction{coinbase}, transactions...)

	b := &block.Block{ time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0 }

	pow := consensus.NewProofOfWork(b)
	nonce, hash := pow.Run()

	b.Hash = hash
	b.Nonce = nonce

	return b
}

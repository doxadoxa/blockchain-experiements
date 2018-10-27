package consensus

import (
	"blockchain-experiements/block"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 16
const maxNonce = math.MaxUint64

type ProofOfWork struct {
	block *block.Block
	target *big.Int
}

func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	fmt.Printf("Target POW: %x\n", target.Bytes())

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce uint) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			[]byte(fmt.Sprintf("%x", pow.block.Timestamp)),
			[]byte(fmt.Sprintf("%x", targetBits)),
			[]byte(fmt.Sprintf("%x", nonce)),
		},
		[]byte{})

	return data
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) Run() (uint, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce uint = 0

	fmt.Printf("Maining block with data \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == - 1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}
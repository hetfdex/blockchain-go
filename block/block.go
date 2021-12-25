package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/hetfdex/blockchain-go/transaction"
)

const (
	difficulty = 10
)

type ProofOfWork interface {
	Validate() bool
	HashTransactions() []byte
}

type Block struct {
	PrevHash     []byte                    `json:"prev_hash"`
	Hash         []byte                    `json:"hash"`
	Nonce        uint64                    `json:"nonce"`
	Target       *big.Int                  `json:"target"`
	Transactions []transaction.Transaction `json:"transactions"`
}

func New(prevHash []byte, transactions []transaction.Transaction) Block {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-difficulty))

	b := Block{
		PrevHash:     prevHash,
		Hash:         []byte{},
		Nonce:        0,
		Target:       target,
		Transactions: transactions,
	}

	b.prove()

	return b
}

func NewGenesis(to string) Block {
	return New(
		[]byte{},
		[]transaction.Transaction{
			transaction.NewGenesis(to),
		},
	)
}

func (b *Block) Validate() bool {
	var intHash big.Int

	data := b.init(b.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(b.Target) == -1
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) prove() {
	var intHash big.Int
	var hash [32]byte

	var nonce uint64

	for nonce < math.MaxInt64 {
		data := b.init(nonce)

		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(b.Target) == -1 {
			break
		}
		nonce++
	}
	b.Nonce = nonce
	b.Hash = hash[:]
}

func (b *Block) init(nonce uint64) []byte {
	return bytes.Join(
		[][]byte{
			b.PrevHash,
			toHex(nonce),
			toHex(difficulty),
			b.HashTransactions(),
		},
		[]byte{},
	)
}

func toHex(i uint64) []byte {
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.BigEndian, i)

	return buffer.Bytes()
}

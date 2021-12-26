package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
	"time"

	"github.com/hetfdex/blockchain-go/transaction"
)

const (
	difficulty          = 10
	genesisPreviousHash = "genesis_previous_hash"
)

type Block struct {
	PreviousHash    []byte
	Hash            []byte
	Nonce           uint64
	TargetDificulty *big.Int
	Transactions    []transaction.Transaction
	CreatedAt       time.Time
}

func New(previousHash []byte, transactions []transaction.Transaction) Block {
	b := Block{
		PreviousHash:    previousHash,
		Hash:            []byte{},
		Nonce:           0,
		TargetDificulty: makeTargetDifficulty(),
		Transactions:    transactions,
		CreatedAt:       time.Now().UTC(),
	}

	b.mine()

	return b
}

func Genesis() Block {
	return New([]byte(genesisPreviousHash), []transaction.Transaction{transaction.Genesis()})
}

func (b *Block) Validate() bool {
	var intHash big.Int

	data := b.makeHashData(b.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(b.TargetDificulty) == -1
}

func makeTargetDifficulty() *big.Int {
	targetDificulty := big.NewInt(1)

	return targetDificulty.Lsh(targetDificulty, uint(256-difficulty))
}

func (b *Block) mine() {
	var intHash big.Int
	var hash [32]byte

	var nonce uint64

	for nonce < math.MaxInt64 {
		data := b.makeHashData(nonce)

		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(b.TargetDificulty) == -1 {
			break
		}
		nonce++
	}
	b.Nonce = nonce
	b.Hash = hash[:]
}

func (b *Block) makeHashData(nonce uint64) []byte {
	return bytes.Join(
		[][]byte{
			b.PreviousHash,
			toHex(nonce),
			toHex(difficulty),
			b.hashTransactions(),
			[]byte(b.CreatedAt.String()),
		},
		[]byte{},
	)
}

func toHex(i uint64) []byte {
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.BigEndian, i)

	return buffer.Bytes()
}

func (b *Block) hashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

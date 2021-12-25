package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

const (
	genesisBlockData = "genesis_block_data"
	difficulty       = 10
)

type ProofOfWork interface {
	Validate() bool
}

type Block struct {
	Data     []byte   `json:"data"`
	PrevHash []byte   `json:"prev_hash"`
	Hash     []byte   `json:"hash"`
	Nonce    uint64   `json:"nonce"`
	Target   *big.Int `json:"target"`
}

func New(data string, prevHash []byte) Block {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-difficulty))

	b := Block{
		Data:     []byte(data),
		PrevHash: prevHash,
		Hash:     []byte{},
		Nonce:    0,
		Target:   target,
	}

	b.prove()

	return b
}

func NewGenesis() Block {
	return New(genesisBlockData, []byte{})
}

func (b *Block) Validate() bool {
	var intHash big.Int

	data := b.init(b.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(b.Target) == -1
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
			b.Data,
			toHex(nonce),
			toHex(difficulty),
		},
		[]byte{},
	)
}

func toHex(i uint64) []byte {
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.BigEndian, i)

	return buffer.Bytes()
}

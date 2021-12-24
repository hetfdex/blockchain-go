package proofofwork

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/hetfdex/blockchain-go/block"
)

const (
	difficulty = 10
)

type ProofOfWork interface {
	Init(nonce uint64) []byte
	Prove() (uint64, []byte)
	Validate() bool
}

type proofOfWork struct {
	Block  *block.Block
	Target *big.Int
}

func New(block *block.Block) ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-difficulty))

	return &proofOfWork{
		Block:  block,
		Target: target,
	}
}

func (p *proofOfWork) Init(nonce uint64) []byte {
	return bytes.Join(
		[][]byte{
			p.Block.PrevHash,
			p.Block.Data,
			toHex(int64(nonce)),
			toHex(int64(difficulty)),
		},
		[]byte{},
	)
}

func (p *proofOfWork) Prove() (uint64, []byte) {
	var intHash big.Int
	var hash [32]byte

	var nonce uint64

	for nonce < math.MaxInt64 {
		data := p.Init(nonce)

		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(p.Target) == -1 {
			break
		}
		nonce++
	}
	return nonce, hash[:]
}

func (p *proofOfWork) Validate() bool {
	var intHash big.Int

	data := p.Init(p.Block.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(p.Target) == -1
}

func toHex(i int64) []byte {
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.BigEndian, i)

	return buffer.Bytes()
}

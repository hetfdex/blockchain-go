package proofofwork

import (
	"testing"

	"github.com/hetfdex/blockchain-go/block"
	"github.com/stretchr/testify/assert"
)

var b = &block.Block{
	Data:     []byte("test_data"),
	PrevHash: []byte("test_prev_hash"),
}

func TestInit(t *testing.T) {
	n := uint64(666)

	pow := New(b)

	res := pow.Init(n)

	assert.NotNil(t, res)
}

func TestProve(t *testing.T) {
	pow := New(b)

	n, h := pow.Prove()

	assert.NotNil(t, n)
	assert.Equal(t, 32, len(h))
}

func TestValidate_False(t *testing.T) {
	pow := New(b)

	_, h := pow.Prove()

	b.Nonce = 666
	b.Hash = h

	assert.False(t, pow.Validate())
}

func TestValidate_True(t *testing.T) {
	pow := New(b)

	n, h := pow.Prove()

	b.Nonce = n
	b.Hash = h

	assert.True(t, pow.Validate())
}

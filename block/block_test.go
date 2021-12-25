package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	data = "test_data"
)

var (
	prevHash = []byte("test_prev_hash")
)

func TestNewBlock(t *testing.T) {
	b := New(data, prevHash)

	assert.Equal(t, []byte(data), b.Data)
	assert.Equal(t, prevHash, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
}

func TestNewGenesisBlock(t *testing.T) {
	b := NewGenesis()

	assert.Equal(t, []byte(genesisBlockData), b.Data)
	assert.Equal(t, []byte{}, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
}

func TestValidate_False(t *testing.T) {
	b := New(data, prevHash)

	b.Nonce = 666

	assert.False(t, b.Validate())
}

func TestValidate_True(t *testing.T) {
	b := New(data, prevHash)

	assert.True(t, b.Validate())
}

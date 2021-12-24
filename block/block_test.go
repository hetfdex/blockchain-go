package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
	data := "test_data"
	prevHash := []byte("test_prev_hash")

	b := New(data, prevHash)

	assert.Equal(t, []byte(data), b.Data)
	assert.Equal(t, prevHash, b.PrevHash)
	assert.Equal(t, uint64(0), b.Nonce)
	assert.Equal(t, 0, len(b.Hash))
}

func TestNewGenesisBlock(t *testing.T) {
	b := NewGenesis()

	assert.Equal(t, []byte(genesisBlockData), b.Data)
	assert.Equal(t, []byte{}, b.PrevHash)
	assert.Equal(t, uint64(0), b.Nonce)
	assert.Equal(t, 0, len(b.Hash))
}

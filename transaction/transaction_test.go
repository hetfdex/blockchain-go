package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	to = "hetfdex"
)

func TestNew(t *testing.T) {
	inputs := []TxInput{}
	outputs := []TxOutput{}

	tx := New(inputs, outputs)

	assert.NotNil(t, tx)
	assert.Equal(t, inputs, tx.Inputs)
	assert.Equal(t, outputs, tx.Outputs)
	assert.Equal(t, 32, len(tx.ID))
}

func TestNewGenesis(t *testing.T) {
	tx := NewGenesis(to)

	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(tx.Inputs))
	assert.Equal(t, 1, len(tx.Outputs))
	assert.Equal(t, 32, len(tx.ID))
	assert.Equal(t, uint64(0), tx.Inputs[0].OutputIndex)
	assert.Equal(t, genesisData, tx.Inputs[0].Signature)
	assert.Equal(t, uint64(reward), tx.Outputs[0].Value)
	assert.Equal(t, to, tx.Outputs[0].PublicKey)
}

func TestValidGenesisTransaction_True(t *testing.T) {
	tx := NewGenesis(to)

	assert.True(t, tx.ValidGenesisTransaction())
}

func TestValidGenesisTransaction_False(t *testing.T) {
	tx := NewGenesis(to)

	tx.ID = []byte("test")

	assert.False(t, tx.ValidGenesisTransaction())

}

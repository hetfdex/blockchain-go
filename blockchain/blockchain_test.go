package blockchain

import (
	"testing"

	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	b := block.Genesis()

	bc := New()

	assert.NotNil(t, bc)
	assert.Equal(t, 1, len(bc.Blocks))
	assert.Equal(t, b.PreviousHash, bc.Blocks[0].PreviousHash)
	assert.Equal(t, b.Transactions, bc.Blocks[0].Transactions)
}

func TestAdd(t *testing.T) {
	txs := []transaction.Transaction{{ID: []byte("test_id")}}

	bc := New()

	bc.Add(txs)

	assert.NotNil(t, bc)
	assert.Equal(t, 2, len(bc.Blocks))
	assert.Equal(t, bc.Blocks[0].Hash, bc.Blocks[1].PreviousHash)
	assert.Equal(t, txs, bc.Blocks[1].Transactions)
}

package block

import (
	"testing"

	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/assert"
)

const (
	data = "test_data"
)

var (
	prevHash = []byte("test_prev_hash")
)

func TestNewBlock(t *testing.T) {
	tx := transaction.Transaction{}

	b := New(data, prevHash, []transaction.Transaction{tx})

	assert.Equal(t, []byte(data), b.Data)
	assert.Equal(t, prevHash, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
}

func TestNewGenesisBlock(t *testing.T) {
	to := "hetfdex"

	b := NewGenesis(data, to)

	tx := transaction.NewGenesisTransaction(data, to)

	assert.Equal(t, []byte(genesisBlockData), b.Data)
	assert.Equal(t, []byte{}, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])

}

func TestValidate_False(t *testing.T) {
	b := New(data, prevHash, []transaction.Transaction{})

	b.Nonce = 666

	assert.False(t, b.Validate())
}

func TestValidate_True(t *testing.T) {
	b := New(data, prevHash, []transaction.Transaction{})

	assert.True(t, b.Validate())
}

func TestHashTransactions(t *testing.T) {
	tx := transaction.Transaction{}

	b := New(data, prevHash, []transaction.Transaction{tx})

	res := b.HashTransactions()

	assert.Equal(t, 32, len(res))
}

package block

import (
	"testing"
	"time"

	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/assert"
)

var (
	previousHash = []byte("test_previous_hash")
)

func TestNew(t *testing.T) {
	tx := transaction.Transaction{ID: []byte("test_id")}

	b := New(previousHash, []transaction.Transaction{tx})

	tm := time.Now().UTC()

	assert.NotNil(t, b)
	assert.Equal(t, previousHash, b.PreviousHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce > 0)
	assert.Equal(t, 31, len(b.TargetDificulty.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
	assert.True(t, b.CreatedAt.Before(tm))
}

func TestGenesis(t *testing.T) {
	tx := transaction.Genesis()

	b := Genesis()

	tm := time.Now().UTC()

	assert.NotNil(t, b)
	assert.Equal(t, []byte(genesisPreviousHash), b.PreviousHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.TargetDificulty.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
	assert.True(t, b.CreatedAt.Before(tm))
}

func TestValidate_False(t *testing.T) {
	b := New(previousHash, []transaction.Transaction{})

	b.Nonce = 666

	assert.NotNil(t, b)
	assert.False(t, b.Validate())
}

func TestValidate_True(t *testing.T) {
	b := New(previousHash, []transaction.Transaction{})

	assert.NotNil(t, b)
	assert.True(t, b.Validate())
}

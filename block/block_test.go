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

	assert.True(t, b.CreatedAt.Before(tm))
	assert.Equal(t, previousHash, b.PreviousHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce > 0)
	assert.Equal(t, 31, len(b.TargetDificulty.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
}

func TestGenesis(t *testing.T) {
	tx := transaction.Genesis()

	b := Genesis()

	tm := time.Now().UTC()

	assert.True(t, b.CreatedAt.Before(tm))
	assert.Equal(t, []byte(genesisPreviousHash), b.PreviousHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.TargetDificulty.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
}

func TestValidHash_False(t *testing.T) {
	b := New(previousHash, []transaction.Transaction{})

	b.Nonce = 666

	assert.False(t, b.ValidHash())
}

func TestValidHash_True(t *testing.T) {
	b := New(previousHash, []transaction.Transaction{})

	assert.True(t, b.ValidHash())
}

func TestValidGenesis_False_PreviousHash(t *testing.T) {
	b := Genesis()

	b.PreviousHash = []byte("fake_test_previous_hash")

	assert.False(t, b.ValidGenesis())
}

func TestValidGenesis_False_Tx_Len(t *testing.T) {
	b := Genesis()

	b.Transactions = append(b.Transactions, transaction.Transaction{ID: []byte("fake_test_id")})

	assert.False(t, b.ValidGenesis())
}

func TestValidGenesis_False_Tx_Genesis(t *testing.T) {
	b := Genesis()

	b.Transactions[0] = transaction.Transaction{ID: []byte("fake_test_id")}

	assert.False(t, b.ValidGenesis())
}

func TestValidGenesis_True(t *testing.T) {
	b := Genesis()

	assert.True(t, b.ValidGenesis())
}

package block

import (
	"testing"

	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/assert"
)

var (
	prevHash = []byte("test_prev_hash")
)

func TestNewBlock(t *testing.T) {
	tx := transaction.Transaction{}

	b := New(prevHash, []transaction.Transaction{tx})

	assert.Equal(t, prevHash, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])
}

func TestNewGenesisBlock(t *testing.T) {
	to := "hetfdex"

	b := NewGenesis(to)

	tx := transaction.NewGenesis(to)

	assert.Equal(t, []byte{}, b.PrevHash)
	assert.Equal(t, 32, len(b.Hash))
	assert.True(t, b.Nonce != 0)
	assert.Equal(t, 31, len(b.Target.Bytes()))
	assert.Equal(t, tx, b.Transactions[0])

}

func TestValidate_False(t *testing.T) {
	b := New(prevHash, []transaction.Transaction{})

	b.Nonce = 666

	assert.False(t, b.Validate())
}

func TestValidate_True(t *testing.T) {
	b := New(prevHash, []transaction.Transaction{})

	assert.True(t, b.Validate())
}

func TestHashTransactions(t *testing.T) {
	tx := transaction.Transaction{}

	b := New(prevHash, []transaction.Transaction{tx})

	res := b.HashTransactions()

	assert.Equal(t, 32, len(res))
}

//Needs proper testing
func TestFindUnspentTransactions(t *testing.T) {
	from := "from"
	to := "to"

	prevHash := []byte("prev_hash")

	tx1 := transaction.New(
		[]transaction.TxInput{
			{
				ID:          []byte("626c61"),
				OutputIndex: 0,
				Signature:   to,
			},
			{
				ID:          []byte("626c65"),
				OutputIndex: 1,
				Signature:   from,
			},
		},
		[]transaction.TxOutput{
			{
				Value:     10,
				PublicKey: from,
			},
			{
				Value:     20,
				PublicKey: to,
			},
		})

	tx2 := transaction.New(
		[]transaction.TxInput{
			{
				ID:          []byte("626c65"),
				OutputIndex: 0,
				Signature:   from,
			},
			{
				ID:          []byte("626c61"),
				OutputIndex: 1,
				Signature:   to,
			},
		},
		[]transaction.TxOutput{
			{
				Value:     50,
				PublicKey: to,
			},
			{
				Value:     30,
				PublicKey: from,
			},
		})

	transactions := []transaction.Transaction{tx1, tx2}

	b := New(prevHash, transactions)

	res := b.FindUnspentTransactions(to)

	assert.Equal(t, transactions, res)
}

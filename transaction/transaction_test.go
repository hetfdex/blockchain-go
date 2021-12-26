package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	txIns = []TransactionInput{{
		ID:          []byte("test_id"),
		OutputIndex: 0,
		Signature:   "test_signature",
	}}

	txOuts = []TransactionOutput{{
		Value:     666,
		PublicKey: "test_public_key",
	}}
)

func TestNew(t *testing.T) {
	tx := New(txIns, txOuts)

	assert.NotNil(t, tx)
	assert.Equal(t, txIns, tx.TxInputs)
	assert.Equal(t, txOuts, tx.TxOutputs)
	assert.Equal(t, 32, len(tx.ID))
}

func TestGenesis(t *testing.T) {
	tx := Genesis()

	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(tx.TxInputs))
	assert.Equal(t, 1, len(tx.TxOutputs))
	assert.Equal(t, 32, len(tx.ID))
	assert.Equal(t, uint64(0), tx.TxInputs[0].OutputIndex)
	assert.Equal(t, genesisSender, tx.TxInputs[0].Signature)
	assert.Equal(t, uint64(rewardValue), tx.TxOutputs[0].Value)
	assert.Equal(t, genesisReceiver, tx.TxOutputs[0].PublicKey)
}

func TestEqual_False_Tx_ID(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	txB.ID = []byte("fake_id")

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxOuts_Len(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	txB.TxOutputs = append(txB.TxOutputs, TransactionOutput{
		Value:     1985,
		PublicKey: "another_test_pubic_key",
	})

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxIn_Id(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	cp := make([]TransactionInput, len(txIns))

	copy(cp, txIns)

	cp[0].ID = []byte("fake_test_id")

	txB.TxInputs = cp

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxOut_OutputIndex(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	cp := make([]TransactionInput, len(txIns))

	copy(cp, txIns)

	cp[0].OutputIndex = 666

	txB.TxInputs = cp

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxOut_Signature(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	cp := make([]TransactionInput, len(txIns))

	copy(cp, txIns)

	cp[0].Signature = "fake_signature"

	txB.TxInputs = cp

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxIns_Len(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	txB.TxInputs = append(txB.TxInputs, TransactionInput{
		ID:          []byte("fake_test_id"),
		OutputIndex: 1,
		Signature:   "fake_signature",
	})

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxOut_Value(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	cp := make([]TransactionOutput, len(txOuts))

	copy(cp, txOuts)

	cp[0].Value = 0

	txB.TxOutputs = cp

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_False_TxOut_PublicKey(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	cp := make([]TransactionOutput, len(txOuts))

	copy(cp, txOuts)

	cp[0].PublicKey = "fake_public_key"

	txB.TxOutputs = cp

	res := Equal(txA, txB)

	assert.False(t, res)
}

func TestEqual_True(t *testing.T) {
	txA := New(txIns, txOuts)
	txB := New(txIns, txOuts)

	res := Equal(txA, txB)

	assert.True(t, res)
}

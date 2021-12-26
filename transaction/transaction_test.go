package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	inputs := []TransactionInput{{ID: []byte("test_id")}}
	outputs := []TransactionOutput{{PublicKey: "test_public_key"}}

	tx := New(inputs, outputs)

	assert.NotNil(t, tx)
	assert.Equal(t, inputs, tx.TxInputs)
	assert.Equal(t, outputs, tx.TxOutputs)
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

package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidGenesisTransaction_True(t *testing.T) {
	tx := NewGenesisTransaction("data", "hetfdex")

	assert.True(t, tx.ValidGenesisTransaction())
}

func TestValidGenesisTransaction_False(t *testing.T) {
	tx := NewGenesisTransaction("data", "hetfdex")

	tx.ID = []byte("test")

	assert.False(t, tx.ValidGenesisTransaction())

}

func TestSetId_Ok(t *testing.T) {
	tx := Transaction{}

	err := tx.SetID()

	assert.Nil(t, err)
}

func TestValidSignature_True(t *testing.T) {
	signature := "test"

	in := TxInput{
		Signature: signature,
	}

	assert.True(t, in.ValidSignature(signature))
}

func TestValidSignature_False(t *testing.T) {
	signature := "test"

	in := TxInput{
		Signature: "signature",
	}

	assert.False(t, in.ValidSignature(signature))
}

func TestValidPublicKey_True(t *testing.T) {
	publicKey := "test"

	out := TxOutput{
		PublicKey: publicKey,
	}

	assert.True(t, out.ValidPublicKey(publicKey))
}

func TestValidPublicKey_False(t *testing.T) {
	publicKey := "test"

	out := TxOutput{
		PublicKey: "publicKey",
	}

	assert.False(t, out.ValidPublicKey(publicKey))
}

package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

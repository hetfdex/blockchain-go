package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

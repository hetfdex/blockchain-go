package wallet

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	kp := []byte("tbd")

	w := New()

	assert.Equal(t, uint64(startingBalance), w.Balance)
	assert.Equal(t, kp, w.Keypair)
	assert.Equal(t, kp, w.PublicKey)
}

//func TestSign_Err(t *testing.T) {}

func TestSign_Ok(t *testing.T) {
	outputMap := make(map[string]uint64)

	expectedRes, err := json.Marshal(outputMap)

	if err != nil {
		t.Fatal(err)
	}

	w := New()

	res, err := w.Sign(outputMap)

	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
}

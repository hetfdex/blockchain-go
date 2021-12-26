package transaction

import (
	"fmt"
	"testing"
	"time"

	"github.com/hetfdex/blockchain-go/wallet"
	"github.com/stretchr/testify/assert"
)

const (
	amount = uint64(666)
)

var (
	w = wallet.New()

	recipientAddress = []byte("test_recipient_address")
)

//func TestNew_Err(t *testing.T) {}

func TestNew_Ok(t *testing.T) {
	res, err := New(w, recipientAddress, amount)

	tm := time.Now().UTC()

	senderSignature, er := w.Sign(res.OutputMap)

	if er != nil {
		t.Fatal(er)
	}

	assert.Nil(t, err)
	assert.Equal(t, 32, len(res.ID))
	assert.Equal(t, amount, res.OutputMap[string(recipientAddress)])
	assert.Equal(t, w.Balance-amount, res.OutputMap[string(w.PublicKey)])
	assert.True(t, res.TxInput.CreatedAt.Before(tm))
	assert.Equal(t, w.PublicKey, res.TxInput.SenderAddress)
	assert.Equal(t, w.Balance, res.TxInput.Balance)
	assert.Equal(t, senderSignature, res.TxInput.SenderSignature)
}

func TestGenesis(t *testing.T) {
	res := MinerReward(w)

	tm := time.Now().UTC()

	assert.Equal(t, 32, len(res.ID))
	assert.Equal(t, uint64(minerRewardAmount), res.OutputMap[string(w.PublicKey)])
	assert.True(t, res.TxInput.CreatedAt.Before(tm))
	assert.Equal(t, []byte(minerRewardSender), res.TxInput.SenderAddress)
}

func TestUpdate_ErrUpdate(t *testing.T) {
	tx, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	err := tx.Update(w, recipientAddress, amount)

	assert.EqualError(t, err, errUpdate.Error())
}

//func TestUpdate_ErrSign(t *testing.T) {}

func TestUpdate_AnotherRecipientAddress(t *testing.T) {
	otherRecipientAddress := []byte("another_test_recipient_address")

	updateAmount := amount / 2

	res, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	finalSenderAmount := res.OutputMap[string(w.PublicKey)] - updateAmount

	err := res.Update(w, otherRecipientAddress, updateAmount)

	senderSignature, er := w.Sign(res.OutputMap)

	if er != nil {
		t.Fatal(er)
	}

	tm := time.Now().UTC()

	assert.Nil(t, err)
	assert.Equal(t, updateAmount, res.OutputMap[string(otherRecipientAddress)])
	assert.Equal(t, finalSenderAmount, res.OutputMap[string(w.PublicKey)])
	assert.True(t, res.TxInput.CreatedAt.Before(tm))
	assert.Equal(t, w.PublicKey, res.TxInput.SenderAddress)
	assert.Equal(t, w.Balance, res.TxInput.Balance)
	assert.Equal(t, senderSignature, res.TxInput.SenderSignature)

}

func TestUpdate_SameRecipientAddress(t *testing.T) {
	updateAmount := amount / 2

	res, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	finalSenderAmount := res.OutputMap[string(w.PublicKey)] - updateAmount
	finalRecipientAmount := res.OutputMap[string(recipientAddress)] + updateAmount

	err := res.Update(w, recipientAddress, updateAmount)

	senderSignature, er := w.Sign(res.OutputMap)

	if er != nil {
		t.Fatal(er)
	}

	tm := time.Now().UTC()

	fmt.Println(res.OutputMap[string(recipientAddress)])
	fmt.Println(res.OutputMap[string(recipientAddress)] + updateAmount)

	assert.Nil(t, err)
	assert.Equal(t, finalRecipientAmount, res.OutputMap[string(recipientAddress)])
	assert.Equal(t, finalSenderAmount, res.OutputMap[string(w.PublicKey)])
	assert.True(t, res.TxInput.CreatedAt.Before(tm))
	assert.Equal(t, w.PublicKey, res.TxInput.SenderAddress)
	assert.Equal(t, w.Balance, res.TxInput.Balance)
	assert.Equal(t, senderSignature, res.TxInput.SenderSignature)
}

func TestValid_False_OutputTotal(t *testing.T) {
	tx, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	tx.OutputMap[string(recipientAddress)] = amount * 2

	assert.False(t, tx.Valid())
}

//func TestValid_False_Verify(t *testing.T) {}

func TestValid_True(t *testing.T) {
	tx, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	assert.True(t, tx.Valid())
}

var ()

func TestEqual_False_TxID(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.ID = []byte("fake_test_id")

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxOutLen(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.OutputMap = make(map[string]uint64)

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxOutEqualAmount(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	outpuMap := make(map[string]uint64, len(txA.OutputMap))

	for add := range txA.OutputMap {
		outpuMap[add] = txA.OutputMap[add] * 2
	}

	txB.OutputMap = outpuMap

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxInCreatedAt(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.TxInput.CreatedAt = time.Now().UTC()

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxInSenderAddress(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.TxInput.SenderAddress = []byte("fake_sender_address")

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxInBalance(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.TxInput.Balance = uint64(666)

	assert.False(t, Equal(txA, txB))
}

func TestEqual_False_TxInSenderSignature(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	txB.TxInput.SenderSignature = []byte("fake_sender_signature")

	assert.False(t, Equal(txA, txB))
}

func TestEqual_True(t *testing.T) {
	txA, er := New(w, recipientAddress, amount)

	if er != nil {
		t.Fatal(er)
	}

	txB := txA

	assert.True(t, Equal(txA, txB))
}

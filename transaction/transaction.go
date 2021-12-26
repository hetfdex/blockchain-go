package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
)

const (
	rewardValue     = 50
	genesisSender   = "genesis_sender"
	genesisReceiver = "genesis_receiver"
)

type Transaction struct {
	ID []byte
	//CreatedAt
	TxInputs  []TransactionInput
	TxOutputs []TransactionOutput
}

type TransactionInput struct {
	ID          []byte
	OutputIndex uint64
	Signature   string
}

type TransactionOutput struct {
	Value     uint64
	PublicKey string
}

func New(txInputs []TransactionInput, txOutputs []TransactionOutput) Transaction {
	tx := Transaction{
		TxInputs:  txInputs,
		TxOutputs: txOutputs,
	}

	tx.setID()

	return tx
}

func Genesis() Transaction {
	inputs := []TransactionInput{
		{
			OutputIndex: 0,
			Signature:   genesisSender,
		},
	}

	outputs := []TransactionOutput{
		{
			Value:     rewardValue,
			PublicKey: genesisReceiver,
		},
	}

	return New(inputs, outputs)
}

func (t *Transaction) setID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := json.NewEncoder(&encoded)

	err := encoder.Encode(t)

	if err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())

	t.ID = hash[:]

	return nil
}

func Equal(a Transaction, b Transaction) bool {
	if !bytes.Equal(a.ID, b.ID) {
		return false
	}

	txOutsA := a.TxOutputs
	txOutsB := b.TxOutputs

	if len(txOutsA) != len(txOutsB) {
		return false
	}

	for i, txOutA := range txOutsA {
		txOutB := txOutsB[i]

		if txOutA.Value != txOutB.Value {
			return false
		}

		if txOutA.PublicKey != txOutB.PublicKey {
			return false
		}
	}

	txInsA := a.TxInputs
	txInsB := b.TxInputs

	if len(txInsA) != len(txInsB) {
		return false
	}

	for i, txInA := range txInsA {
		txInB := txInsB[i]

		if !bytes.Equal(txInA.ID, txInB.ID) {
			return false
		}

		if txInA.OutputIndex != txInB.OutputIndex {
			return false
		}

		if txInA.Signature != txInB.Signature {
			return false
		}

	}
	return true
}

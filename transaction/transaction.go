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
	ID        []byte
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

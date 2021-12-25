package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
)

const (
	reward = 100

	genesisData = "genesis_data"
)

type TransactionValidator interface {
	ValidGenesisTransaction() bool
}

type Transaction struct {
	ID      []byte     `json:"id"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

func New(inputs []TxInput, outputs []TxOutput) Transaction {
	tx := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	tx.setID()

	return tx
}

func NewGenesis(to string) Transaction {
	inputs := []TxInput{
		{
			OutputIndex: 0,
			Signature:   genesisData,
		},
	}

	outputs := []TxOutput{
		{
			Value:     reward,
			PublicKey: to,
		},
	}

	return New(inputs, outputs)
}

func (t *Transaction) ValidGenesisTransaction() bool {
	if len(t.ID) == 32 && len(t.Inputs) == 1 && len(t.Outputs) == 1 {
		input := t.Inputs[0]
		output := t.Outputs[0]

		if input.ID == nil && input.OutputIndex == 0 {
			return output.Value == reward
		}
	}
	return false
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

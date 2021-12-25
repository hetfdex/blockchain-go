package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

const (
	reward      = 100
	defaultData = "transaction for %s"
)

type TxOutput struct {
	Value     uint64 `json:"value"`
	PublicKey string `json:"public_key"`
}

type TxInput struct {
	ID          []byte `json:"id"`
	OutputIndex uint64 `json:"output_index"`
	Signature   string `json:"signature"`
}

type Transaction struct {
	ID      []byte     `json:"id"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

func NewGenesisTransaction(data string, to string) Transaction {
	if data == "" {
		data = fmt.Sprintf(defaultData, to)
	}
	txInput := TxInput{
		ID:          nil,
		OutputIndex: 0,
		Signature:   data,
	}

	txOutput := TxOutput{
		Value:     reward,
		PublicKey: to,
	}

	return Transaction{
		ID: nil,
		Inputs: []TxInput{
			txInput,
		},
		Outputs: []TxOutput{
			txOutput,
		},
	}
}

func (t *Transaction) SetID() error {
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

func (t *Transaction) ValidGenesisTransaction() bool {
	if t.ID == nil && len(t.Inputs) == 1 && len(t.Outputs) == 1 {
		input := t.Inputs[0]
		output := t.Outputs[0]

		if input.ID == nil && input.OutputIndex == 0 {
			return output.Value == reward
		}
	}
	return false
}

func (i *TxInput) ValidSignature(data string) bool {
	return i.Signature == data
}

func (o *TxOutput) ValidPublicKey(data string) bool {
	return o.PublicKey == data
}

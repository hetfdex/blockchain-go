package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/transaction"
)

const (
	latestBlockHashKey = "latest_block_hash_key"
)

var (
	errNotEnoughBalance = errors.New("not enough balance")
)

type Blockchain interface {
	Set(block.Block) error
	Get([]byte) (block.Block, error)
	GetLatest() (block.Block, error)
	NewTransaction(string, string, uint64) (transaction.Transaction, error)
	FindUnspentTxOutputs(string) ([]transaction.TxOutput, error)
}

type blockchain struct {
	wrapper     badgerwrapper.BadgerWrapper
	latestBlock block.Block
}

func New(wrapper badgerwrapper.BadgerWrapper) *blockchain {
	return &blockchain{
		wrapper: wrapper,
	}
}

func (bc *blockchain) Set(block block.Block) error {
	jsonBlock, err := json.Marshal(block)

	if err != nil {
		return err
	}

	err = bc.wrapper.Set(block.Hash, jsonBlock)

	if err != nil {
		return err
	}

	err = bc.wrapper.Set([]byte(latestBlockHashKey), block.Hash)

	if err != nil {
		return err
	}
	bc.latestBlock = block

	return nil
}

func (bc *blockchain) Get(hash []byte) (block.Block, error) {
	var b block.Block

	jsonBlock, err := bc.wrapper.Get(hash)

	if err != nil {
		return b, err
	}

	err = json.Unmarshal(jsonBlock, &b)

	if err != nil {
		return b, err
	}
	return b, nil
}

func (bc *blockchain) GetLatest() (block.Block, error) {
	if bc.latestBlock.Hash != nil && len(bc.latestBlock.Hash) > 0 {
		return bc.latestBlock, nil
	}

	latestBlockHash, err := bc.wrapper.Get([]byte(latestBlockHashKey))

	if err != nil {
		return bc.latestBlock, err
	}

	latestBlock, err := bc.Get(latestBlockHash)

	if err != nil {
		return bc.latestBlock, err
	}
	bc.latestBlock = latestBlock

	return latestBlock, nil
}

func (bc *blockchain) NewTransaction(from string, to string, amount uint64) (transaction.Transaction, error) {
	var res transaction.Transaction
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	balance, validOutputs, err := bc.findSpendableOutputs(from, amount)

	if err != nil {
		return res, err
	}

	if balance < amount {
		return res, errNotEnoughBalance
	}

	for i, outs := range validOutputs {
		id, err := hex.DecodeString(i)

		if err != nil {
			return res, err
		}

		for _, out := range outs {
			input := transaction.TxInput{
				ID:          id,
				OutputIndex: out,
				Signature:   from,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, transaction.TxOutput{
		Value:     amount,
		PublicKey: to,
	})

	if balance > amount {
		outputs = append(outputs, transaction.TxOutput{
			Value:     balance - amount,
			PublicKey: from,
		})
	}

	return transaction.New(inputs, outputs), nil
}

func (bc *blockchain) FindUnspentTxOutputs(address string) ([]transaction.TxOutput, error) {
	unspentTxOuts := []transaction.TxOutput{}

	unspentTxs, err := bc.findUnspentTransactions(address)

	if err != nil {
		return nil, err
	}

	for _, tx := range unspentTxs {
		for _, out := range tx.Outputs {
			if out.ValidPublicKey(address) {
				unspentTxOuts = append(unspentTxOuts, out)
			}
		}
	}

	return unspentTxOuts, nil
}

func (bc *blockchain) findSpendableOutputs(address string, amount uint64) (uint64, map[string][]uint64, error) {
	unspentOuts := make(map[string][]uint64)

	unspentTxs, err := bc.findUnspentTransactions(address)

	if err != nil {
		return 0, nil, err
	}

	var funds uint64

Work:
	for _, tx := range unspentTxs {
		id := hex.EncodeToString(tx.ID)
		for i, out := range tx.Outputs {
			if out.ValidPublicKey(address) && funds < amount {
				funds += out.Value

				unspentOuts[id] = append(unspentOuts[id], uint64(i))

				if funds >= amount {
					break Work
				}
			}
		}
	}
	return funds, unspentOuts, nil
}

func (bc *blockchain) findUnspentTransactions(address string) ([]transaction.Transaction, error) {
	unspentTxs := []transaction.Transaction{}

	latestBlock, err := bc.GetLatest()

	if err != nil {
		return nil, err
	}

	unspentTxs = append(unspentTxs, latestBlock.FindUnspentTransactions(address)...)

	for {
		block, err := bc.Get(latestBlock.PrevHash)

		if err != nil {
			return nil, err
		}

		unspentTxs = append(unspentTxs, block.FindUnspentTransactions(address)...)

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs, nil
}

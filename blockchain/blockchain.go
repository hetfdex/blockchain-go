package blockchain

import (
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/transaction"
)

type blockchain struct {
	Blocks []block.Block
}

func New() *blockchain {
	return &blockchain{
		Blocks: []block.Block{
			block.Genesis(),
		},
	}
}

func (bc *blockchain) Add(transactions []transaction.Transaction) {
	previousHash := bc.Blocks[len(bc.Blocks)-1].Hash

	b := block.New(previousHash, transactions)

	bc.Blocks = append(bc.Blocks, b)
}

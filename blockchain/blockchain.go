package blockchain

import (
	"bytes"

	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/hetfdex/blockchain-go/wallet"
)

type Blockchain struct {
	Blocks []block.Block
}

func New() *Blockchain {
	return &Blockchain{
		Blocks: []block.Block{
			block.Genesis(),
		},
	}
}

/*createTransaction ({recipient, amount, chain}) {
  if(chain) {
    this.balance = Wallet.calculateBalance({chain, address: this.publicKey});
  }

  if (amount > this.balance) {
    throw new Error("Creation Fail: Amount exceeds Balance");
  }

  return new Transaction({senderWallet: this, recipient, amount});
}*/

func (bc *Blockchain) Add(transactions []transaction.Transaction) {
	previousHash := bc.Blocks[len(bc.Blocks)-1].Hash

	b := block.New(previousHash, transactions)

	bc.Blocks = append(bc.Blocks, b)
}

func (bc *Blockchain) Valid() bool {
	if !bc.Blocks[0].ValidGenesis() {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		if !bytes.Equal(currentBlock.PreviousHash, previousBlock.Hash) {
			return false
		}

		if !currentBlock.Valid() {
			return false
		}

		for j := 0; j < len(currentBlock.Transactions); j++ {
			tx := currentBlock.Transactions[i]

			if !tx.Valid() {
				return false
			}
		}
	}

	return true
}

func (bc *Blockchain) Balance(address []byte) uint64 {
	var outputsTotal uint64

	hasTransacted := false

	for i := len(bc.Blocks) - 1; i > 0; i-- {
		b := bc.Blocks[i]

		for _, tx := range b.Transactions {
			if bytes.Equal(tx.TxInput.SenderAddress, address) {
				hasTransacted = true
			}
			addressOutput := tx.OutputMap[string(address)]

			outputsTotal += addressOutput
		}
		if hasTransacted {
			break
		}
	}
	if hasTransacted {
		return outputsTotal
	}
	return wallet.StartingBalance + outputsTotal
}

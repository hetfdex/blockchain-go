package orchestrator

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/blockchain"
	"github.com/hetfdex/blockchain-go/transaction"
)

const (
	dbPath     = "./tmp/blocks"
	dbManifest = "./tmp/blocks/MANIFEST"
)

func InitDb() (*badger.DB, badgerwrapper.BadgerWrapper, error) {
	opts := badger.DefaultOptions(dbPath).WithLoggingLevel(badger.WARNING)

	db, err := badger.Open(opts)

	if err != nil {
		return nil, nil, err
	}
	return db, badgerwrapper.New(db), nil
}

func InitBlockchain(wrapper badgerwrapper.BadgerWrapper) (blockchain.Blockchain, error) {
	bc := blockchain.New(wrapper)

	_, err := bc.GetLatest()

	if err != nil {
		if err != badger.ErrKeyNotFound {
			return nil, err
		}
		err = bc.Set(block.NewGenesis("hetfdex"))

		if err != nil {
			return nil, err
		}
	}

	return bc, nil
}
func AddBlock(bc blockchain.Blockchain, transactions []transaction.Transaction) error {
	previousBlock, err := bc.GetLatest()

	if err != nil {
		return err
	}

	b := block.New(previousBlock.Hash, transactions)

	err = bc.Set(b)

	if err != nil {
		return err
	}
	return nil
}

func PrintBlocks(bc blockchain.Blockchain) error {
	block, err := bc.GetLatest()

	if err != nil {
		return err
	}

	printBlock(block)

	for {
		block, err = bc.Get(block.PrevHash)

		if err != nil {
			return err
		}

		printBlock(block)

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return nil
}

func printBlock(block block.Block) {
	fmt.Printf("Prev Hash: %x\n", block.PrevHash)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Valid: %s\n", strconv.FormatBool(block.Validate()))
	fmt.Println()
}

package orchestrator

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/blockchain"
	"github.com/hetfdex/blockchain-go/proofofwork"
)

const (
	path = "./tmp/blocks"
)

func InitDb() (*badger.DB, badgerwrapper.BadgerWrapper, error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.WARNING)

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
		if err == badgerwrapper.ErrBlockNotFound {
			err := makeNewBlockchain(bc)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return bc, nil
}

func makeNewBlockchain(bc blockchain.Blockchain) error {
	genesisBlock := block.NewGenesis()

	pow := proofofwork.New(&genesisBlock)

	nonce, hash := pow.Prove()

	genesisBlock.Hash = hash[:]
	genesisBlock.Nonce = nonce

	err := bc.Set(genesisBlock)

	if err != nil {
		return err
	}
	return nil
}

func AddBlock(bc blockchain.Blockchain, data string) error {
	previousBlock, err := bc.GetLatest()

	if err != nil {
		return err
	}

	b := block.New(data, previousBlock.Hash)

	pow := proofofwork.New(&b)

	nonce, hash := pow.Prove()

	b.Hash = hash[:]
	b.Nonce = nonce

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
	pow := proofofwork.New(&block)

	fmt.Printf("Data: %s\n", block.Data)
	fmt.Printf("Prev Hash: %x\n", block.PrevHash)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Valid: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println()
}

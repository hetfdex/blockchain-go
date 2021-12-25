package blockchain

import (
	"encoding/json"

	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
)

const (
	genesisData        = "genesis data"
	latestBlockHashKey = "latest_block_hash_key"
)

type Blockchain interface {
	Set(block block.Block) error
	Get(hash []byte) (block.Block, error)
	GetLatest() (block.Block, error)
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

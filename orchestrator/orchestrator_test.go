package orchestrator

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/blockchain"
	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	key = []byte("key")

	errExpected = errors.New("error")

	prevHash = []byte("test_prev_hash")

	genesisBlock = block.NewGenesis("hetfdex")
)

func TestInitDb(t *testing.T) {
	_ = os.RemoveAll("./tmp")

	res, wrapper, err := InitDb()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, wrapper)

	_ = os.RemoveAll("./tmp")

	_ = res.Close()
}

func TestInitBlockchain_ErrGetLatest(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, errExpected)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestInitBlockchain_ErrKeyNotFound_ErrSet(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, badger.ErrKeyNotFound)
	wrapper.On("Set", mock.Anything, mock.Anything).Return(errExpected)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestInitBlockchain_Restored(t *testing.T) {
	b := block.New(prevHash, []transaction.Transaction{})

	value, err := json.Marshal(b)

	if err != nil {
		t.Fatal(err)
	}

	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", []byte("latest_block_hash_key")).Return(key, nil)
	wrapper.On("Get", key).Return(value, nil)

	res, err := InitBlockchain(&wrapper)

	latestBlock, er := res.GetLatest()

	if er != nil {
		t.Fatal(er)
	}

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, b, latestBlock)
}

func TestInitBlockchain_Created(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, badger.ErrKeyNotFound)
	wrapper.On("Set", mock.Anything, mock.Anything).Return(nil)

	res, err := InitBlockchain(&wrapper)

	latestBlock, er := res.GetLatest()

	if er != nil {
		t.Fatal(er)
	}

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, genesisBlock, latestBlock)
}

func TestAddBlock_ErrGetLatest(t *testing.T) {
	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(block.Block{}, errExpected)

	err := AddBlock(&bc, []transaction.Transaction{})

	assert.EqualError(t, err, errExpected.Error())
}

func TestAddBlock_ErrSet(t *testing.T) {
	b := block.New(prevHash, []transaction.Transaction{})

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Set", mock.Anything).Return(errExpected)

	err := AddBlock(&bc, []transaction.Transaction{})

	assert.EqualError(t, err, errExpected.Error())
}

func TestAddBlock_Ok(t *testing.T) {
	b := block.New(prevHash, []transaction.Transaction{})

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Set", mock.Anything).Return(nil)

	err := AddBlock(&bc, []transaction.Transaction{})

	assert.Nil(t, err)
}

func TestPrintBlocks_ErrGetLatest(t *testing.T) {
	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(block.Block{}, errExpected)

	err := PrintBlocks(&bc)

	assert.EqualError(t, err, errExpected.Error())
}

func TestPrintBlocks_ErrGet(t *testing.T) {
	b := block.New(prevHash, []transaction.Transaction{})

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Get", mock.Anything).Return(block.Block{}, errExpected)

	err := PrintBlocks(&bc)

	assert.EqualError(t, err, errExpected.Error())
}

func TestPrintBlocks_Ok(t *testing.T) {
	b2 := block.New(genesisBlock.Hash, []transaction.Transaction{})
	b3 := block.New(b2.Hash, []transaction.Transaction{})

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b3, nil)
	bc.On("Get", b3.PrevHash).Return(b2, nil)
	bc.On("Get", b2.PrevHash).Return(genesisBlock, nil)

	err := PrintBlocks(&bc)

	assert.Nil(t, err)
}

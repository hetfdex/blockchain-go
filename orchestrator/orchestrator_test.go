package orchestrator

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	data = "test_data"
)

var (
	errExpected = errors.New("error")
	prevHash    = []byte("test_prev_hash")
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

func TestInitBlockchain_MakeNew_Err(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, badgerwrapper.ErrBlockNotFound)
	wrapper.On("Set", mock.Anything, mock.Anything).Return(errExpected)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestInitBlockchain_Err(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, errExpected)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestInitBlockchain_MakeNew(t *testing.T) {
	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, badgerwrapper.ErrBlockNotFound)
	wrapper.On("Set", mock.Anything, mock.Anything).Return(nil)

	res, err := InitBlockchain(&wrapper)

	latestBlock, er := res.GetLatest()

	if er != nil {
		t.Fatal(er)
	}

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, block.NewGenesis().Data, latestBlock.Data)
}

func TestInitBlockchain_Restored(t *testing.T) {
	b := block.New(data, prevHash)

	value, err := json.Marshal(b)

	if err != nil {
		t.Fatal(err)
	}

	key := []byte("key")

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
	assert.Equal(t, b.Data, latestBlock.Data)
}

func TestAddBlock_ErrGetLatest(t *testing.T) {
	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(block.Block{}, errExpected)

	err := AddBlock(&bc, "test")

	assert.EqualError(t, err, errExpected.Error())
}

func TestAddBlock_ErrSet(t *testing.T) {
	b := block.New(data, prevHash)

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Set", mock.Anything).Return(errExpected)

	err := AddBlock(&bc, "test")

	assert.EqualError(t, err, errExpected.Error())
}

func TestAddBlock_Ok(t *testing.T) {
	b := block.New(data, prevHash)

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Set", mock.Anything).Return(nil)

	err := AddBlock(&bc, "test")

	assert.Nil(t, err)
}

func TestPrintBlocks_ErrGetLatest(t *testing.T) {
	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(block.Block{}, errExpected)

	err := PrintBlocks(&bc)

	assert.EqualError(t, err, errExpected.Error())
}

func TestPrintBlocks_ErrGet(t *testing.T) {
	b := block.New(data, prevHash)

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b, nil)
	bc.On("Get", mock.Anything).Return(block.Block{}, errExpected)

	err := PrintBlocks(&bc)

	assert.EqualError(t, err, errExpected.Error())
}

func TestPrintBlocks_Ok(t *testing.T) {
	b1 := block.NewGenesis()
	b2 := block.New(data, b1.Hash)
	b3 := block.New(data, b2.Hash)

	bc := blockchain.BlockchainMock{}

	bc.On("GetLatest").Return(b3, nil)
	bc.On("Get", b3.PrevHash).Return(b2, nil)
	bc.On("Get", b2.PrevHash).Return(b1, nil)

	err := PrintBlocks(&bc)

	assert.Nil(t, err)
}

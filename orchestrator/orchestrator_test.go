package orchestrator

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	expectedErr := errors.New("error")

	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, badgerwrapper.ErrBlockNotFound)
	wrapper.On("Set", mock.Anything, mock.Anything).Return(expectedErr)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, expectedErr.Error())

}

func TestInitBlockchain_Err(t *testing.T) {
	expectedErr := errors.New("error")

	wrapper := badgerwrapper.BadgerWrapperMock{}

	wrapper.On("Get", mock.Anything).Return([]byte{}, expectedErr)

	res, err := InitBlockchain(&wrapper)

	assert.Nil(t, res)
	assert.EqualError(t, err, expectedErr.Error())

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
	b := &block.Block{
		Data:     []byte("data"),
		PrevHash: []byte("previous_hash"),
		Hash:     []byte("hash"),
		Nonce:    0,
	}

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

func TestAddBlock(t *testing.T) {
	//getlatest err
	//set err
	//ok
}

func TestPrintBlocks(t *testing.T) {
	//getlatest err
	//get err
	//ok
}

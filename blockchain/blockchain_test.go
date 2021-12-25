package blockchain

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/hetfdex/blockchain-go/badgerwrapper"
	"github.com/hetfdex/blockchain-go/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	b = block.Block{
		PrevHash: []byte("previous_hash"),
		Hash:     []byte("hash"),
		Nonce:    0,
	}
	errExpected = errors.New("error")

	key = []byte("key")
)

func TestSet_ErrBlockHash(t *testing.T) {

	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Set", b.Hash, mock.Anything).Return(errExpected)

	bc := New(&wrapperMock)

	err := bc.Set(b)

	assert.EqualError(t, err, errExpected.Error())
}

func TestSet_ErrLatestBlockHashKey(t *testing.T) {
	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Set", b.Hash, mock.Anything).Return(nil)
	wrapperMock.On("Set", []byte(latestBlockHashKey), mock.Anything).Return(errExpected)

	bc := New(&wrapperMock)

	err := bc.Set(b)

	assert.EqualError(t, err, errExpected.Error())
}

func TestSet_Ok(t *testing.T) {
	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Set", mock.Anything, mock.Anything).Return(nil)

	bc := New(&wrapperMock)

	err := bc.Set(b)

	latestBlock, er := bc.GetLatest()

	if er != nil {
		t.Fatal(er)
	}

	assert.Nil(t, err)
	assert.Equal(t, b, latestBlock)
}

func TestGet_ErrGet(t *testing.T) {
	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", key).Return([]byte{}, errExpected)

	bc := New(&wrapperMock)

	res, err := bc.Get(key)

	assert.Equal(t, block.Block{}, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestGet_ErrUnmarshall(t *testing.T) {
	value := []byte("value")

	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", key).Return(value, nil)

	bc := New(&wrapperMock)

	res, err := bc.Get(key)

	assert.Equal(t, block.Block{}, res)
	assert.NotNil(t, err)
}

func TestGet_OK(t *testing.T) {
	value, err := json.Marshal(b)

	if err != nil {
		t.Fatal(err)
	}

	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", key).Return(value, nil)

	bc := New(&wrapperMock)

	res, err := bc.Get(key)

	assert.Nil(t, err)
	assert.Equal(t, b, res)
}

func TestGetLatest_OkFromLatestBlock(t *testing.T) {
	bc := blockchain{
		latestBlock: b,
	}

	res, err := bc.GetLatest()

	assert.Nil(t, err)
	assert.Equal(t, b, res)
}

func TestGetLatest_ErrLatestBlockHashKey(t *testing.T) {
	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", []byte(latestBlockHashKey)).Return([]byte{}, errExpected)

	bc := New(&wrapperMock)

	res, err := bc.GetLatest()

	assert.Equal(t, block.Block{}, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestGetLatest_ErrLatestBlockHash(t *testing.T) {
	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", []byte(latestBlockHashKey)).Return(key, nil)
	wrapperMock.On("Get", key).Return([]byte{}, errExpected)

	bc := New(&wrapperMock)

	res, err := bc.GetLatest()

	assert.Equal(t, block.Block{}, res)
	assert.EqualError(t, err, errExpected.Error())
}

func TestGetLatest_OkFromDb(t *testing.T) {
	value, err := json.Marshal(b)

	if err != nil {
		t.Fatal(err)
	}

	wrapperMock := badgerwrapper.BadgerWrapperMock{}

	wrapperMock.On("Get", []byte(latestBlockHashKey)).Return(key, nil)
	wrapperMock.On("Get", key).Return(value, nil)

	bc := New(&wrapperMock)

	res, err := bc.GetLatest()

	assert.Nil(t, err)
	assert.Equal(t, b, res)
}

package blockchain

import (
	"github.com/hetfdex/blockchain-go/block"
	"github.com/hetfdex/blockchain-go/transaction"
	"github.com/stretchr/testify/mock"
)

type BlockchainMock struct {
	mock.Mock
}

func (m *BlockchainMock) Set(block block.Block) error {
	args := m.Called(block)

	return args.Error(0)
}

func (m *BlockchainMock) Get(hash []byte) (block.Block, error) {
	args := m.Called(hash)

	return args.Get(0).(block.Block), args.Error(1)
}

func (m *BlockchainMock) GetLatest() (block.Block, error) {
	args := m.Called()

	return args.Get(0).(block.Block), args.Error(1)
}

func (m *BlockchainMock) NewTransaction(from string, to string, amount uint64) (transaction.Transaction, error) {
	args := m.Called(from, to, amount)

	return args.Get(0).(transaction.Transaction), args.Error(1)
}
func (m *BlockchainMock) FindUnspentTxOutputs(address string) ([]transaction.TxOutput, error) {
	args := m.Called(address)

	return args.Get(0).([]transaction.TxOutput), args.Error(1)
}

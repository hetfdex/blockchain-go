package blockchain

import (
	"github.com/hetfdex/blockchain-go/block"
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

package badgerwrapper

import "github.com/stretchr/testify/mock"

type BadgerWrapperMock struct {
	mock.Mock
}

func (m *BadgerWrapperMock) Set(key []byte, value []byte) error {
	args := m.Called(key, value)

	return args.Error(0)

}

func (m *BadgerWrapperMock) Get(key []byte) ([]byte, error) {
	args := m.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

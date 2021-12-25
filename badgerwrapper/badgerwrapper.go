package badgerwrapper

import (
	"errors"

	"github.com/dgraph-io/badger/v3"
)

var (
	ErrBlockNotFound = errors.New("block not found")
)

type BadgerWrapper interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
}

type badgerWrapper struct {
	db *badger.DB
}

func New(db *badger.DB) *badgerWrapper {
	return &badgerWrapper{
		db: db,
	}
}

func (w *badgerWrapper) Set(key []byte, value []byte) error {
	err := w.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return err
}

func (w *badgerWrapper) Get(key []byte) ([]byte, error) {
	var res []byte

	err := w.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)

		if err != nil {
			return err
		}

		res, err = item.ValueCopy(nil)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, handleErrKeyNotFound(err)
	}

	return res, nil
}

func handleErrKeyNotFound(err error) error {
	if err == badger.ErrKeyNotFound {
		return ErrBlockNotFound
	}
	return err
}

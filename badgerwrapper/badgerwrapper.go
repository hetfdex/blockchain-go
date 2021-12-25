package badgerwrapper

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerWrapper interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
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
		return nil, err
	}

	return res, nil
}

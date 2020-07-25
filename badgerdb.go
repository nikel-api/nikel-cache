package cache

import (
	"github.com/dgraph-io/badger/v2"
)

// BadgerDB struct
type BadgerDB struct {
	DB *badger.DB
}

// NewBadgerDB initializes BadgerDB database
func NewBadgerDB(path string) (*BadgerDB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	bdb := BadgerDB{
		DB: db,
	}
	return &bdb, err
}

// Get cache value
func (bdb *BadgerDB) Get(key string) ([]byte, error) {
	var data []byte
	err := bdb.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	return data, err
}

// Set cache value
func (bdb *BadgerDB) Set(key string, value []byte) error {
	return bdb.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

// Remove cache value
func (bdb *BadgerDB) Remove(key string) error {
	return bdb.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

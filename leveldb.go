package cache

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// LevelDB struct
type LevelDB struct {
	db *leveldb.DB
}

// NewLevelDB initializes LevelDB database
func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	ldb := LevelDB{
		db: db,
	}
	return &ldb, err
}

// Get cache value
func (ldb *LevelDB) Get(key string) ([]byte, error) {
	return ldb.db.Get([]byte(key), nil)
}

// Set cache value
func (ldb *LevelDB) Set(key string, value []byte) error {
	return ldb.db.Put([]byte(key), value, nil)
}

// Remove cache value
func (ldb *LevelDB) Remove(key string) error {
	return ldb.db.Delete([]byte(key), nil)
}

// Update cache value
func (ldb *LevelDB) Update(key string, value []byte) error {
	batch := new(leveldb.Batch)
	batch.Delete([]byte(key))
	batch.Put([]byte(key), value)
	return ldb.db.Write(batch, nil)
}

// Keys returns all keys
func (ldb *LevelDB) Keys() []string {
	var cumul []string
	iter := ldb.db.NewIterator(util.BytesPrefix([]byte(KeyPrefix)), nil)
	for iter.Next() {
		cumul = append(cumul, string(iter.Key()))
	}
	iter.Release()
	return cumul
}

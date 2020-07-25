package cache

import (
	"github.com/syndtr/goleveldb/leveldb"
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
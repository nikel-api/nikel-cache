package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestBadgerDB tests BadgerDB implementation
func TestBadgerDB(t *testing.T) {
	bdb, err := NewBadgerDB("BadgerCache")
	assert.NoError(t, err)

	err = bdb.Set("foo", []byte{1, 2, 3})
	assert.NoError(t, err)

	res, err := bdb.Get("foo")
	assert.NoError(t, err)
	assert.Len(t, res, 3)

	err = bdb.Remove("foo")
	assert.NoError(t, err)

	err = bdb.DB.Close()
	assert.NoError(t, err)

	err = os.RemoveAll("BadgerCache")
	assert.NoError(t, err)
}

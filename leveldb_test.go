package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestLDB tests LevelDB implementation
func TestLDB(t *testing.T) {
	ldb, err := NewLevelDB("LDBCache")
	assert.NoError(t, err)

	err = ldb.Set("foo", []byte{1, 2, 3})
	assert.NoError(t, err)

	res, err := ldb.Get("foo")
	assert.NoError(t, err)
	assert.Len(t, res, 3)

	err = ldb.Remove("foo")
	assert.NoError(t, err)

	err = ldb.DB.Close()
	assert.NoError(t, err)

	err = os.RemoveAll("LDBCache")
	assert.NoError(t, err)
}

package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestIMGet tests Get
func TestIMGet(t *testing.T) {
	im := InMemory{
		hash: map[string][]byte{
			"foo": {1, 2, 3},
		},
	}

	res, err := im.Get("foo")
	assert.NoError(t, err)
	assert.Len(t, res, 3)

	_, err = im.Get("bar")
	assert.Equal(t, err, errNotFound)
}

// TestIMSet tests Set
func TestIMSet(t *testing.T) {
	im := InMemory{hash: map[string][]byte{}}

	err := im.Set("foo", []byte{1, 2, 3})
	assert.NoError(t, err)

	res, ok := im.hash["foo"]
	assert.True(t, ok)
	assert.Len(t, res, 3)

	err = im.Set("foo", []byte{1, 2, 3})
	assert.Equal(t, err, errAlreadyExists)
}

// TestIMRemove tests Remove
func TestIMRemove(t *testing.T) {
	im := InMemory{
		hash: map[string][]byte{
			"foo": {1, 2, 3},
		},
	}

	err := im.Remove("foo")
	assert.NoError(t, err)

	err = im.Remove("bar")
	assert.Equal(t, err, errNotFound)
	assert.Empty(t, im.hash)
}

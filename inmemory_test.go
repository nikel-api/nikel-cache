package cache

import (
	"reflect"
	"testing"
)

// TestGet tests Get
func TestGet(t *testing.T) {
	im := InMemory{
		hash: map[string][]byte{
			"key": {1, 2, 3},
		},
	}
	// im.Set("key", []byte{1, 2, 3})
	res, err := im.Get("key")
	expect(t, err, nil)
	expect(t, len(res), 3)

	_, err = im.Get("__key__")
	expect(t, err, errNotFound)
}

// TestIMSet tests Set
func TestIMSet(t *testing.T) {
	im := InMemory{hash: map[string][]byte{}}

	err := im.Set("key", []byte{1, 2, 3})
	expect(t, err, nil)

	res, ok := im.hash["key"]
	expect(t, ok, true)
	expect(t, len(res), 3)

	err = im.Set("key", []byte{1, 2, 3})
	expect(t, err, errAlreadyExists)
}

// TestIMRemove tests Remove
func TestIMRemove(t *testing.T) {
	im := InMemory{
		hash: map[string][]byte{
			"key": {1, 2, 3},
		},
	}

	err := im.Remove("key")
	expect(t, err, nil)

	err = im.Remove("__key__")
	expect(t, err, errNotFound)

	expect(t, len(im.hash), 0)
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

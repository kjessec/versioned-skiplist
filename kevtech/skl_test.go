package kevtech

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var latestVersion uint64 = 24
var numRecords uint64 = 100

func TestSkipList(t *testing.T) {
	// 512mb arena
	skl := NewSkipList(512 * 1024 * 1024)

	for v := uint64(1); v <= latestVersion; v++ {
		for i := uint64(0); i < numRecords; i++ {
			kv := toBigEndian(i * 2)
			val := []byte(fmt.Sprintf("key %d at version %d", i*2, v))
			skl.Insert(kv, val, v)
		}
	}

	testBasic(t, skl)
	testIteraor(t, skl)
}

func testBasic(t *testing.T, skl *SkipList) {
	t.Run("get without versions", func(t *testing.T) {
		assert.Equal(t,
			[]byte(nil),
			skl.Get(toBigEndian(5)),
			"key 5 is not in the set",
		)
		assert.Equal(t,
			[]byte(fmt.Sprintf("key %d at version %d", 6, latestVersion)),
			skl.Get(toBigEndian(6)),
			"key 6 exists",
		)
		assert.Nil(t,
			skl.Get(toBigEndian(292)),
			"out of bounds",
		)
	})

	t.Run("get with versions", func(t *testing.T) {
		version := uint64(11)
		assert.Equal(t,
			[]byte(nil),
			skl.GetVersion(toBigEndian(5), version),
			"key 5 is not in the set",
		)
		assert.Equal(t,
			[]byte(fmt.Sprintf("key %d at version %d", 6, version)),
			skl.GetVersion(toBigEndian(6), version),
			"key 6 at version 11 exists",
		)
		assert.Nil(t,
			skl.GetVersion(toBigEndian(292), version),
			"out of bounds",
		)
	})
}

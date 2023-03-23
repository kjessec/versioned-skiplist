package kevtech

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testIteraor(t *testing.T, skl *SkipList) {
	t.Run("iterate without versions", func(t *testing.T) {
		// search range from 13-73, however the database is filled with even numbers.
		// so the actual range is 14-72. This is to test whether .First() method works as
		// intended (along with the internal .getNearestNodePtr() method)
		var start, end = toBigEndian(13), toBigEndian(73)
		it, err := skl.Iterate(start, end)
		assert.NoError(t, err)

		var results [][2][]byte
		for ; it.Valid(); it.Next() {
			results = append(results, [2][]byte{it.Key(), it.Value()})
		}

		assert.Equal(t, 30, len(results))
		assert.Equal(t, toBigEndian(14), results[0][0])
		assert.Equal(t, toBigEndian(72), results[len(results)-1][0])
	})

	t.Run("iterate with versions", func(t *testing.T) {
		t.Fail()
	})
}

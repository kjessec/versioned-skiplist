package kevtech

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	// 512mb arena
	skl := NewSkipList(512 * 1024 * 1024)

	for v := uint64(0); v < 10; v++ {
		for i := 0; i < 10; i++ {
			kv := toBigEndian(uint64(i))
			val := []byte(fmt.Sprintf("key %d at version %d", i, i))
			fmt.Print(fmt.Sprintf("inserting at version %d; %x %x len=%d\n", v, string(kv), string(val), len(val)))

			skl.Insert(kv, val, v)
		}
	}

	fmt.Print(skl.Search(toBigEndian(5)))
	//fmt.Print(skl.SearchWithValue(toBigEndian(5), 4))
}

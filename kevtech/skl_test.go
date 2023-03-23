package kevtech

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	// 512mb arena
	skl := NewSkipList(512 * 1024 * 1024)

	for v := uint64(1); v <= 1; v++ {
		for i := 0; i < 5; i++ {
			kv := toBigEndian(uint64(i))
			val := []byte(fmt.Sprintf("key %d at version %d", i, v))
			skl.Insert(kv, val, v)
		}
	}

	fmt.Println(string(skl.Get(toBigEndian(5))))
	fmt.Println(string(skl.GetVersion(toBigEndian(5), 4)))
}

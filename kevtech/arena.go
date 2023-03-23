package kevtech

import (
	"fmt"
	"unsafe"
)

type Arena []byte

func NewArena(size uint64) *Arena {
	a := &Arena{}
	*a = make([]byte, size)

	return a
}

func (a *Arena) Alloc(size int) unsafe.Pointer {
	if len(*a) < int(size) {
		panic("arena too small")
	}
	fmt.Printf("alloc %d\n", size)

	alloc := (*a)[:size]
	*a = (*a)[size:]

	return unsafe.Pointer(&alloc[0])
}

func (a *Arena) Append(value []byte) (head unsafe.Pointer) {
	ptr := a.Alloc(len(value))
	copy((*(*[1 << 30]byte)(ptr))[:], value)

	return ptr
}

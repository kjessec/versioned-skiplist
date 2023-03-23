package kevtech

import (
	"unsafe"
)

type Arena []byte

const (
	sizeLen = 8
)

func NewArena(size uint64) *Arena {
	a := &Arena{}
	*a = make([]byte, size)

	return a
}

func (a *Arena) Alloc(size int) unsafe.Pointer {
	if len(*a) < int(size) {
		panic("arena too small")
	}

	alloc := (*a)[:size]
	*a = (*a)[size:]

	return unsafe.Pointer(&alloc[0])
}

func (a *Arena) Append(value []byte) (head unsafe.Pointer) {
	ptr := a.Alloc(8 + len(value))
	copy((*(*[1 << 30]byte)(ptr))[:], toBigEndian(uint64(len(value))))
	copy((*(*[1 << 30]byte)(ptr))[8:], value)
	return ptr
}

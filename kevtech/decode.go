package kevtech

import (
	"unsafe"
)

func decode(ptr unsafe.Pointer) []byte {
	size := *(*[8]byte)(ptr)
	sizeUint64 := fromBigEndian(size)

	return (*(*[1 << 30]byte)(ptr))[8 : 8+sizeUint64]
}

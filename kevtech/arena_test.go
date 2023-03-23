package kevtech

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestArena(t *testing.T) {
	ar := NewArena(1024)
	head := &(([]byte)(*ar))[0]
	alloc := ar.Alloc(0x100)
	alloc2 := ar.Alloc(0x100)

	assert.Equal(t, unsafe.Pointer(head), alloc)
	assert.Equal(t, unsafe.Add(unsafe.Pointer(head), 0x100), alloc2)

	appended := []byte("hello world")

	ar.Append(appended)
	alloc3 := ar.Alloc(0x100)
	assert.Equal(t, unsafe.Add(unsafe.Pointer(head), 0x200+len(appended)), alloc3)
}

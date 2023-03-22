package kevtech

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"unsafe"
)

const (
	maxLevel = 32                         // maximum level of the skiplist
	p        = 0.25                       // probability of increasing level
	nodeSize = int(unsafe.Sizeof(node{})) // size of a node in bytes
)

type node struct {
	key            []byte
	isDeleted      bool
	value          unsafe.Pointer
	valueAtVersion map[uint64]unsafe.Pointer
	next           []*node
}

type SkipList struct {
	head     *node
	level    int
	arena    []byte
	valarena []byte
}

// NewSkipList creates a new skip list.
func NewSkipList(arenaSize uint64) *SkipList {
	s := &SkipList{
		head:     &node{next: make([]*node, maxLevel)},
		level:    1,
		arena:    make([]byte, arenaSize),
		valarena: make([]byte, arenaSize*8),
	}
	return s
}

// randomLevel generates a random level for a new node.
func (s *SkipList) randomLevel() int {
	level := 1
	for level < maxLevel && rand.Float64() < p {
		level++
	}
	return level
}

func (s *SkipList) upsertNode(node *node, version uint64, value []byte) {
	// no need to handle valueAtVersion; should already be created when newNode is created
	if _, ok := node.valueAtVersion[version]; ok {
		panic(fmt.Sprintf("value already exists for version %d", value))
	}

	valsize := len(value) + 8
	valarena := s.valarena[:valsize]
	copy(valarena, toBigEndian(uint64(len(value))))
	copy(valarena[8:], value)

	// debug
	fmt.Printf("upsertNode: valarena=%v\n", valarena[:valsize])

	valptr := &valarena[0]

	node.valueAtVersion[version] = unsafe.Pointer(valptr)
	node.value = unsafe.Pointer(valptr)

	// update tail
	s.valarena = s.valarena[valsize:]
}

// insertNode inserts a new node with the given key and value.
func (s *SkipList) Insert(key []byte, value []byte, version uint64) {
	update := make([]*node, maxLevel)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		// Follow next pointers until we reach a node whose next pointer points to a node with a larger key.
		for x.next[i] != nil && bytes.Compare(x.next[i].key, key) < 0 {
			x = x.next[i]
		}

		// If a node with the same key already exists, update its value and return the node.
		if x.next[i] != nil && bytes.Equal(x.next[i].key, key) {
			s.upsertNode(x.next[i], version, value)
			return
		}

		update[i] = x
	}
	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.head
		}
		s.level = level
	}
	nodeBytes := s.arena[:nodeSize]
	s.arena = s.arena[nodeSize:]
	newNode := (*node)(unsafe.Pointer(&nodeBytes[0]))
	newNode.valueAtVersion = make(map[uint64]unsafe.Pointer)

	s.upsertNode(newNode, version, value)
	newNode.key = key
	newNode.next = make([]*node, level)
	for i := 0; i < level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

// Search returns the value associated with the given key, or nil if the key is not found.
func (s *SkipList) Search(key []byte) []byte {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.next[i] != nil && bytes.Compare(x.next[i].key, key) <= 0 {
			if bytes.Equal(x.next[i].key, key) {
				valptr := x.next[i].value
				valsize := *(*uint64)(valptr)
				valval := *(*[]byte)(unsafe.Pointer(valsize + 1))
				valval := (*(*[1 << 30]byte)(valptr))[8:]

				fmt.Println(valsize, valval)

				fmt.Printf("get value %v\n", *(*[]byte)(x.next[i].value))
				return *(*[]byte)(x.next[i].value)
			}
			x = x.next[i]
		}
	}
	return nil
}

func (s *SkipList) SearchWithValue(key []byte, version uint64) []byte {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.next[i] != nil && bytes.Compare(x.next[i].key, key) <= 0 {
			if bytes.Equal(x.next[i].key, key) {
				if val, ok := x.next[i].valueAtVersion[version]; ok {
					return (*(*[1 << 30]byte)(val))[8:]
				}
			}
			x = x.next[i]
		}
	}
	return nil
}

func (s *SkipList) GetHeadAtLevel(level int) *node {
	return s.head.next[level]
}

// Visualize returns a string representation of the skiplist.
func (s *SkipList) Visualize(deskey func([]byte) uint64) string {
	var levels []string
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		var nodes []string
		var prev uint64 = 0
		nodes = append(nodes, fmt.Sprintf("H%d", i))
		for x.next[i] != nil {
			kk := deskey(x.next[i].key)

			for j := 0; j < int(kk-prev-1); j++ {
				nodes = append(nodes, fmt.Sprintf("    "))
			}

			nodes = append(nodes, fmt.Sprintf("|%2d ", kk))
			prev = kk
			x = x.next[i]
		}
		level := strings.Join(nodes, "")
		levels = append(levels, level)
		x = s.head
	}
	return strings.Join(levels, "\n=======================\n")
}

func toBigEndian(n uint64) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[7-i] = byte(n >> (i * 8))
	}
	return b
}

func fromBigEndian(b []byte) uint64 {
	var n uint64
	for i := 0; i < 8; i++ {
		n |= uint64(b[7-i]) << (i * 8)
	}
	return n
}

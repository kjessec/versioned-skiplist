package kevtech

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"unsafe"
)

const (
	maxLevel      = 32                         // maximum level of the skiplist
	p             = 0.25                       // probability of increasing level
	nodeSize      = int(unsafe.Sizeof(node{})) // size of a node in bytes
	versionedSize = int(unsafe.Sizeof(versioned{}))
)

type node struct {
	key       []byte
	value     []byte
	version   uint64
	isDeleted bool

	next []*node
}

type versioned struct {
	version uint64
	valptr  []byte
}

type SkipList struct {
	head        *node
	lastVersion uint64
	level       int

	nodeArena  *Arena
	keyArena   *Arena
	valueArena *Arena
}

// NewSkipList creates a new skip list.
func NewSkipList(arenaSize uint64) *SkipList {
	s := &SkipList{
		head:  &node{next: make([]*node, maxLevel)},
		level: 1,

		nodeArena:   NewArena(arenaSize),
		keyArena:    NewArena(arenaSize),
		valueArena:  NewArena(arenaSize * 8),
		lastVersion: 0,
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

// version is NEVER reset; no need to handle the edge case where
// version being inserted is less than the current version of this specific key
func (s *SkipList) upsertNode(node *node, version uint64, key, value []byte) {
	if s.lastVersion > version {
		panic(fmt.Sprintf("version %d is less than current version %d", version, s.lastVersion))
	}

	// key
	keyptr := s.keyArena.Append(key)
	valueptr := s.valueArena.Append(value)

	node.key = *(*[]byte)(keyptr)
	node.value = *(*[]byte)(valueptr)
	node.version = version
	node.isDeleted = false

	fmt.Println(*(*[]byte)(keyptr))

	// handle versioning link
}

// insertNode inserts a new node with the given key and value.
func (s *SkipList) Insert(key []byte, value []byte, version uint64) {
	update := make([]*node, maxLevel)
	x := s.head

	for i := s.level - 1; i >= 0; i-- {
		next := x.next[i]
		if next == nil {
			update[i] = x
			continue
		}

		nextKey := x.next[i].key

		// Follow next pointers until we reach a node whose next pointer points to a node with a larger key.
		for x.next[i] != nil && bytes.Compare(nextKey, key) < 0 {
			x = x.next[i]
		}

		// If a node with the same key already exists, update its value and return the node.
		if x.next[i] != nil && bytes.Equal(nextKey, key) {
			s.upsertNode(x.next[i], version, key, value)
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

	// assign new node
	newNode := (*node)(s.nodeArena.Alloc(nodeSize))
	s.upsertNode(newNode, version, key, value)

	newNode.next = make([]*node, level)
	for i := 0; i < level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

// Get returns the value associated with the given key, or nil if the key is not found.
func (s *SkipList) Get(key []byte) []byte {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.next[i] != nil && bytes.Compare(x.next[i].key, key) <= 0 {
			if bytes.Equal(x.next[i].key, key) {
				valval := x.next[i].value
				fmt.Printf("???? %v\n", valval)
				return valval
			}
			x = x.next[i]
		}
	}
	return nil
}

func (s *SkipList) getNodePtr(key []byte) *node {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.next[i] != nil && bytes.Compare(x.next[i].key, key) <= 0 {
			if bytes.Equal(x.next[i].key, key) {
				return x.next[i]
			}
			x = x.next[i]
		}
	}
	return x
}

func (s *SkipList) GetVersion(key []byte, version uint64) []byte {
	return nil
	//node := s.getNodePtr(key)
	//
	////valueAtVersion, ok := node.valueAtVersion[version]
	//if !ok {
	//	return nil
	//}
	//
	//valptr := (*[]byte)(valueAtVersion)
	//valsize := fromBigEndian(*(valptr))
	//valval := (*valptr)[8 : 8+valsize]

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

func arenaAlloc(arena []byte, size int) []byte {
	if len(arena) < size {
		panic("arena too small")
	}
	alloc := arena[:size]
	arena = arena[size:]
	return alloc
}

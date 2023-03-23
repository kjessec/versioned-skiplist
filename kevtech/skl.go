package kevtech

import (
	"bytes"
	"fmt"
	"math/rand"
	"unsafe"
)

const (
	maxLevel = 32                         // maximum level of the skiplist
	p        = 0.25                       // probability of increasing level
	nodeSize = int(unsafe.Sizeof(node{})) // size of a node in bytes
)

type node struct {
	key       unsafe.Pointer
	lastValue unsafe.Pointer
	vValues   []unsafe.Pointer
	next      []*node
	vNext     [][]*node
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

	node.key = keyptr
	node.lastValue = valueptr

	ivptr := (*internalVersionValue)(s.valueArena.Alloc(sizeIv))
	ivptr.version = version
	ivptr.value = valueptr

	node.vValues = append(node.vValues, unsafe.Pointer(ivptr))

	// handle version links!
	//for _, v := range node.vValues {
	//	iv := (*internalVersionValue)(v)
	//}

}

// insertNode inserts a new node with the given key and lastValue.
func (s *SkipList) Insert(key []byte, value []byte, version uint64) {
	update := make([]*node, maxLevel)
	x := s.head

	for i := s.level - 1; i >= 0; i-- {
		// Follow next pointers until we reach a node whose next pointer points to a node with a larger key.
		for x.next[i] != nil && bytes.Compare(decode(x.next[i].key), key) < 0 {
			x = x.next[i]
		}

		// If a node with the same key already exists, update its lastValue and return the node.
		if x.next[i] != nil && bytes.Equal(decode(x.next[i].key), key) {
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

	// assign new node ahd handle initialization
	newNode := (*node)(s.nodeArena.Alloc(nodeSize))
	newNode.vValues = make([]unsafe.Pointer, 0)
	newNode.next = make([]*node, level)
	for i := 0; i < level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
	s.upsertNode(newNode, version, key, value)
}

func (s *SkipList) getNodePtr(key []byte) *node {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.next[i] != nil {
			nkey := decode(x.next[i].key)
			cmp := bytes.Compare(nkey, key)
			if cmp == -1 {
				x = x.next[i]
			} else if cmp == 0 {
				return x.next[i]
			} else {
				break
			}
		}
	}

	// not found; this is out-of-bounds error
	return nil
}

// Get returns the lastValue associated with the given key, or nil if the key is not found.
func (s *SkipList) Get(key []byte) []byte {
	node := s.getNodePtr(key)
	if node == nil {
		return nil
	}

	return decode(s.getNodePtr(key).lastValue)
}

func (s *SkipList) GetVersion(key []byte, version uint64) []byte {
	node := s.getNodePtr(key)
	if node == nil {
		return nil
	}

	verptr := findNearestVersion(node.vValues, version)
	iv := (*internalVersionValue)(verptr)

	return decode(iv.value)

}

func toBigEndian(n uint64) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[7-i] = byte(n >> (i * 8))
	}
	return b
}

func fromBigEndian(b [8]byte) uint64 {
	var n uint64
	for i := 0; i < 8; i++ {
		n |= uint64(b[7-i]) << (i * 8)
	}
	return n
}

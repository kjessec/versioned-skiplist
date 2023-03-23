package kevtech

import (
	"bytes"
)

type iterator struct {
	skl   *SkipList
	start []byte
	end   []byte

	cursor  *node
	version uint64
}

func (s *SkipList) getNearestNodePtr(key []byte) *node {
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

	return x.next[0]
}

func (s *SkipList) Iterate(start, end []byte) (*iterator, error) {
	return s.IterateWithVersion(start, end, s.lastVersion)
}

func (s *SkipList) IterateWithVersion(start, end []byte, version uint64) (*iterator, error) {
	it := &iterator{
		skl:   s,
		start: start,
		end:   end,

		cursor:  nil,
		version: s.lastVersion,
	}

	// load up first cursor
	it.First()

	return it, nil
}

func (it *iterator) First() {
	it.cursor = it.skl.getNearestNodePtr(it.start)
}

func (it *iterator) Next() {
	it.cursor = it.cursor.next[0]
}

func (it *iterator) Valid() bool {
	return it.cursor != nil && bytes.Compare(decode(it.cursor.key), it.end) == -1
}

func (it *iterator) Key() []byte {
	return decode(it.cursor.key)
}

func (it *iterator) Value() []byte {
	return decode(it.cursor.lastValue)
}

func (it *iterator) Close() {
	*it = *(*iterator)(nil)
}

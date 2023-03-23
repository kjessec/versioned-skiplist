package kevtech

type iterator struct {
	skl   *SkipList
	start []byte
	end   []byte

	cursor  *node
	version uint64
}

func (s *SkipList) Iterate(start, end []byte) (*iterator, error) {
	return s.IterateWithVersion(start, end, s.version)
}

func (s *SkipList) IterateWithVersion(start, end []byte, version uint64) (*iterator, error) {
	it := &iterator{
		skl:   s,
		start: start,
		end:   end,

		cursor:  nil,
		version: s.version,
	}

	it.Next()

	return it, nil
}

func (it *iterator) First() {
	it.cursor = it.skl.getNodePtr(it.start)
}

func (it *iterator) Next() {
	it.cursor = it.cursor.next[0]
}

func (it *iterator) Valid() bool {
	return it.cursor != nil
}

func (it *iterator) Key() []byte {
	return it.cursor.valueAtVersion[it.version]
}

func (it *iterator) Value() []byte {
	return nil
}

func (it *iterator) Close() {

}

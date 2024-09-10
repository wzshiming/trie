package trie

import (
	"fmt"
	"io"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

// MatchWithReader returns most matching handler and prefix bytes data to use for the given reader.
func (t *Trie[T]) MatchWithReader(r io.Reader) (handler T, prefix []byte, err error) {
	if t.Size() == 0 {
		return handler, nil, ErrNotFound
	}
	parent := t.Mapping()
	off := 0
	prefix = make([]byte, t.Depth())
	for {
		i, err := r.Read(prefix[off:])
		if err != nil {
			return handler, nil, err
		}
		if i == 0 {
			break
		}

		data, next, ok := parent.Get(prefix[off : off+i])
		if ok && data != nil {
			handler = data
		}

		off += i
		if next == nil {
			break
		}
		parent = next
	}
	if handler == nil {
		return handler, prefix[:off], ErrNotFound
	}
	return handler, prefix[:off], nil
}

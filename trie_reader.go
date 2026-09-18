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
	var zero T
	off := 0
	prefix = make([]byte, t.Depth())
	found := false
	for {
		count, err := r.Read(prefix[off:])
		off += count
		if count > 0 {
			data, _, ok, more := t.mapping.get(nil, prefix[:off], zero, false)
			if ok {
				handler = data
				found = true
			}
			if !more {
				break
			}
		}
		if err != nil {
			return handler, prefix[:off], err
		}
		if count == 0 {
			break
		}
	}
	if !found {
		return handler, prefix[:off], ErrNotFound
	}
	return handler, prefix[:off], nil
}

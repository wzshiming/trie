package trie

import (
	"fmt"
	"io"
	"strings"
)

// Mapping is a getter of the trie.
type Mapping[T any] struct {
	array [byteLength]*node[T]
}

func (m *Mapping[T]) String() string {
	var buf strings.Builder
	buf.WriteString("\n")
	m.deepString(&buf, nil, 1)
	return buf.String()
}

func (m *Mapping[T]) deepString(w io.Writer, key []byte, deep int) {
	isNil := true
	for i, v := range m.array {
		if v == nil {
			continue
		}
		isNil = false
		deepin(w, deep)
		fmt.Fprintf(w, "[%d] %q:\n", i, []byte{byte(i)})
		v.deepString(w, append(key, byte(i)), deep+1)
	}
	if isNil {
		deepin(w, deep)
		fmt.Fprintf(w, "nil\n")
	}
}

func (m *Mapping[T]) walk(buf []byte, f func(k []byte, v T)) {
	for k, v := range m.array {
		if v == nil {
			continue
		}
		tmp := buf
		tmp = append(tmp, byte(k))
		if len(v.zip) != 0 {
			tmp = append(tmp, v.zip...)
		}

		if v.has {
			f(tmp, v.data)
		}

		if v.mapping != nil {
			v.mapping.walk(tmp, f)
		}
	}
}

func (m *Mapping[T]) put(key []byte, val T) (finish bool) {
	if len(key) == 0 {
		return false
	}

	car := key[0]
	cdr := key[1:]
	if len(cdr) == 0 {
		cdr = nil
	}

	child := m.array[car]
	if child == nil {
		m.array[car] = &node[T]{
			zip:  cdr,
			data: val,
			has:  true,
		}
		child = m.array[car]
		return true
	}

	if len(child.zip) != 0 {
		var diff int
		if len(cdr) != 0 {
			diff = bytesDiff(child.zip, cdr)
			if diff == -1 {
				child.data = val
				child.has = true
				return true
			}
			cdr = cdr[diff:]
		}
		child.split(diff)
		if len(cdr) == 0 {
			child.data = val
			child.has = true
			return true
		}
	}

	if child.mapping == nil {
		child.mapping = &Mapping[T]{}
	}
	return child.mapping.put(cdr, val)
}

// Get returns the val in the trie for a key.
func (m *Mapping[T]) Get(key []byte) (val T, current *Mapping[T], finish bool) {
	return m.get(nil, key, val, finish)
}

func (m *Mapping[T]) get(prev *Mapping[T], key []byte, defaulted T, has bool) (val T, current *Mapping[T], finish bool) {
	if len(key) == 0 {
		return defaulted, prev, has
	}
	car := key[0]
	cdr := key[1:]

	child := m.array[car]
	if child == nil {
		return defaulted, prev, has
	}

	if len(child.zip) != 0 {
		var diff int
		if len(cdr) != 0 {
			diff = bytesDiff(child.zip, cdr)
			if diff == -1 {
				if child.has {
					return child.data, m, true
				}
				return defaulted, prev, has
			}
		}

		if len(child.zip) > diff {
			return defaulted, prev, has
		}

		cdr = cdr[diff:]
	}

	if len(cdr) == 0 {
		if child.has {
			return child.data, m, true
		}
		return defaulted, prev, has
	}

	if child.mapping == nil {
		if child.has {
			return child.data, child.mapping, true
		}
		return defaulted, prev, has
	}

	if child.has {
		return child.mapping.get(m, cdr, child.data, true)
	}
	return child.mapping.get(prev, cdr, defaulted, has)
}

func bytesDiff(a, b []byte) int {
	min := len(a)
	max := len(b)
	if min > max {
		min, max = max, min
	}
	for i := 0; i != min; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if min == max {
		return -1
	}
	return min
}

func deepin(w io.Writer, deep int) {
	for i := 0; i != deep; i++ {
		fmt.Fprint(w, "-")
	}
}

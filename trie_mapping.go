package trie

import (
	"fmt"
	"io"
	"strings"
)

type mapping struct {
	array [byteLength]*node
}

func (m *mapping) String() string {
	var buf strings.Builder
	buf.WriteString("\n")
	m.deepString(&buf, nil, 1)
	return buf.String()
}

func (m *mapping) deepString(w io.Writer, key []byte, deep int) {
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

func (m *mapping) walk(buf []byte, f func(k, v []byte)) {
	for k, v := range m.array {
		if v == nil {
			continue
		}
		buf := append(append(buf, byte(k)), v.zip...)

		if v.data != nil {
			f(buf, v.data)
		}
		v.mapping.walk(buf, f)
	}
}

func (m *mapping) put(key, val []byte) (finish bool) {
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
		child = &node{
			zip:  cdr,
			data: val,
		}
		m.array[car] = child
		return true
	}

	if len(child.zip) != 0 {
		var diff int
		if len(cdr) != 0 {
			diff = bytesDiff(child.zip, cdr)
			if diff == -1 {
				child.data = val
				return true
			}
			cdr = cdr[diff:]
		}
		child.split(diff)
		if len(cdr) == 0 {
			child.data = val
			return true
		}
	}

	return child.mapping.put(cdr, val)
}

func (m *mapping) get(key []byte, defaulted []byte) (val []byte, finish bool) {
	if len(key) == 0 {
		return defaulted, false
	}
	car := key[0]
	cdr := key[1:]

	child := m.array[car]
	if child == nil {
		return defaulted, false
	}

	if len(child.zip) != 0 {
		var diff int
		if len(cdr) != 0 {
			diff = bytesDiff(child.zip, cdr)
			if diff == -1 {
				return child.data, true
			}
		}

		if len(child.zip) > diff {
			return defaulted, false
		}

		cdr = cdr[diff:]
	}

	if len(cdr) == 0 {
		return child.data, child.data != nil
	}

	if child.data != nil {
		defaulted = child.data
	}
	return child.mapping.get(cdr, defaulted)
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

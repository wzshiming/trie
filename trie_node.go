package trie

import (
	"fmt"
	"io"
	"strings"
)

type node[T any] struct {
	mapping *Mapping[T]
	data    T
	has     bool
	zip     []byte
}

func (n *node[T]) String() string {
	var buf strings.Builder
	buf.WriteString("\n")
	n.deepString(&buf, nil, 1)
	return buf.String()
}

func (n *node[T]) deepString(w io.Writer, key []byte, deep int) {
	key = append(key, n.zip...)
	deepin(w, deep)
	fmt.Fprintf(w, "key:  %v %q\n", key, key)
	deepin(w, deep)
	fmt.Fprintf(w, "zip:  %v %q\n", n.zip, n.zip)
	deepin(w, deep)
	if n.has {
		fmt.Fprintf(w, "data: %v\n", n.data)
	} else {
		fmt.Fprintf(w, "data: <empty>\n")
	}
	deepin(w, deep)
	if n.mapping != nil {
		fmt.Fprintf(w, "Mapping:\n")
		n.mapping.deepString(w, key, deep+1)
	} else {
		fmt.Fprintf(w, "Mapping: <empty>\n")
	}
}

func (n *node[T]) split(off int) {
	if len(n.zip) <= off {
		return
	}
	key := n.zip[off:]

	if off == 0 {
		n.zip = nil
	} else {
		n.zip = n.zip[:off]
	}

	car := key[0]
	cdr := key[1:]
	if len(cdr) == 0 {
		cdr = nil
	}

	var m Mapping[T]
	m.array[car] = &node[T]{
		mapping: n.mapping,
		zip:     cdr,
		data:    n.data,
		has:     n.has,
	}
	n.mapping = &m

	// empty data
	var t T
	n.data = t
	n.has = false
}

func (n *node[T]) join() {
	var child *node[T]
	var car byte
	for i, v := range n.mapping.array {
		if v == nil {
			continue
		}
		car = byte(i)
		child = v
		break
	}
	if child == nil {
		return
	}

	n.zip = append(n.zip, car)
	n.zip = append(n.zip, child.zip...)
	n.data = child.data
	n.has = child.has
	n.mapping = child.mapping
}

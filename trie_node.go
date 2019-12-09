package trie

import (
	"fmt"
	"io"
	"strings"
)

type node struct {
	mapping mapping
	data    []byte
	zip     []byte
}

func (n *node) String() string {
	var buf strings.Builder
	buf.WriteString("\n")
	n.deepString(&buf, nil, 1)
	return buf.String()
}

func (n *node) deepString(w io.Writer, key []byte, deep int) {
	key = append(key, n.zip...)
	deepin(w, deep)
	fmt.Fprintf(w, "key:  %v %q\n", key, key)
	deepin(w, deep)
	fmt.Fprintf(w, "zip:  %v %q\n", n.zip, n.zip)
	deepin(w, deep)
	fmt.Fprintf(w, "data: %v %q\n", n.data, n.data)
	deepin(w, deep)
	fmt.Fprintf(w, "mapping:\n")
	n.mapping.deepString(w, key, deep+1)
}

func (n *node) split(off int) {
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

	var array [byteLength]*node
	array[car] = &node{
		mapping: n.mapping,
		zip:     cdr,
		data:    n.data,
	}
	n.mapping.array = array
	n.data = nil
}

func (n *node) join() {
	var child *node
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
	n.mapping = child.mapping
}

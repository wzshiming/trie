package trie

import (
	"reflect"
	"testing"
)

func Test_node_split(t *testing.T) {
	got := &node[[]byte]{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
		has:  true,
	}
	got.split(3)
	want1 := &node[[]byte]{
		zip:  []byte{1, 2, 3},
		data: nil,
		mapping: &mapping[[]byte]{
			array: [byteLength]*node[[]byte]{
				nil, nil, nil, nil,
				/* 4 */ {zip: []byte{5}, data: []byte{1}, has: true},
			},
		},
	}
	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.split(2)
	want2 := &node[[]byte]{
		zip:  []byte{1, 2},
		data: nil,
		mapping: &mapping[[]byte]{
			array: [byteLength]*node[[]byte]{
				nil, nil, nil,
				/* 3 */ {zip: nil, data: nil, mapping: &mapping[[]byte]{
					array: [byteLength]*node[[]byte]{
						nil, nil, nil, nil,
						/* 4 */ {zip: []byte{5}, data: []byte{1}, has: true},
					},
				}},
			},
		},
	}
	if !reflect.DeepEqual(got, want2) {
		t.Errorf("put() = %s, want %s", got, want2)
	}
}

func Test_node_split_and_join(t *testing.T) {
	got := &node[[]byte]{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
		has:  true,
	}
	want := &node[[]byte]{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
		has:  true,
	}

	for i := 0; i != 5; i++ {
		got.split(i)
		got.join()
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("split() = %s, want %s", got, want)
		}
	}
}

func Test_node_String(t *testing.T) {
	want := &node[[]byte]{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
	}
	t.Log(want.String())
}

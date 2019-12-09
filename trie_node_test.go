package trie

import (
	"reflect"
	"testing"
)

func Test_node_split(t *testing.T) {
	got := &node{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
	}
	got.split(3)
	want1 := &node{
		zip:  []byte{1, 2, 3},
		data: nil,
		mapping: mapping{
			array: [byteLength]*node{
				nil, nil, nil, nil,
				/* 4 */ {zip: []byte{5}, data: []byte{1}},
			},
		},
	}
	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.split(2)
	want2 := &node{
		zip:  []byte{1, 2},
		data: nil,
		mapping: mapping{
			array: [byteLength]*node{
				nil, nil, nil,
				/* 3 */ {zip: nil, data: nil, mapping: mapping{
					array: [byteLength]*node{
						nil, nil, nil, nil,
						/* 4 */ {zip: []byte{5}, data: []byte{1}},
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
	got := &node{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
	}
	want := &node{
		zip:  []byte{1, 2, 3, 4, 5},
		data: []byte{1},
	}

	for i := 0; i != 5; i++ {
		got.split(i)
		got.join()
		if !reflect.DeepEqual(want, got) {
			t.Fail()
		}
	}
}

package trie

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_bytesDiff(t *testing.T) {

	tests := []struct {
		a    []byte
		b    []byte
		want int
	}{
		{[]byte{}, []byte{}, -1},
		{[]byte{1}, []byte{1}, -1},
		{[]byte{1}, []byte{0}, 0},
		{[]byte{1, 1, 1}, []byte{0, 0, 0}, 0},
		{[]byte{1, 1, 1}, []byte{1, 0, 0}, 1},
		{[]byte{1, 1, 1}, []byte{1, 1, 0}, 2},
		{[]byte{1, 1, 1}, []byte{1, 1, 1}, -1},
		{[]byte{1, 1, 1}, []byte{1, 1, 1, 1}, 3},
		{[]byte{1, 1, 1, 1}, []byte{1, 1, 1}, 3},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := bytesDiff(tt.a, tt.b); got != tt.want {
				t.Errorf("bytesDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapping_put_1(t *testing.T) {
	got := &mapping{}
	got.put([]byte{1, 2, 3}, []byte{0})
	want1 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: []byte{2, 3}, data: []byte{0}},
		},
	}

	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.put([]byte{1, 2, 3, 4, 5}, []byte{1})
	want2 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: []byte{2, 3}, data: []byte{0}, mapping: mapping{
				array: [byteLength]*node{
					nil, nil, nil, nil,
					/* 4 */ {zip: []byte{5}, data: []byte{1}},
				}},
			},
		},
	}
	if !reflect.DeepEqual(got, want2) {
		t.Errorf("put() = %s, want %s", got, want2)
	}
	got.put([]byte{1, 2, 4, 5}, []byte{2})
	want3 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: []byte{2}, data: nil, mapping: mapping{
				array: [byteLength]*node{
					nil, nil, nil,
					/* 3 */ {zip: nil, data: []byte{0}, mapping: mapping{
						array: [byteLength]*node{
							nil, nil, nil, nil,
							/* 4 */ {zip: []byte{5}, data: []byte{1}},
						}},
					},
					/* 4 */ {zip: []byte{5}, data: []byte{2}},
				}},
			},
		},
	}

	if !reflect.DeepEqual(got, want3) {
		t.Errorf("put() = %s, want %s", got, want3)
	}
}

func Test_mapping_put_2(t *testing.T) {
	got := &mapping{}
	got.put([]byte{1}, []byte{0})
	want1 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}},
		},
	}

	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.put([]byte{1, 2, 3}, []byte{1})
	want2 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}, mapping: mapping{
				array: [byteLength]*node{
					nil, nil,
					/* 2*/ {zip: []byte{3}, data: []byte{1}},
				},
			}},
		},
	}

	if !reflect.DeepEqual(got, want2) {
		t.Errorf("put() = %s, want %s", got, want2)
	}
	got.put([]byte{1, 2}, []byte{2})
	want3 := &mapping{
		array: [byteLength]*node{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}, mapping: mapping{
				array: [byteLength]*node{
					nil, nil,
					/* 2*/ {zip: nil, data: []byte{2}, mapping: mapping{
						array: [byteLength]*node{
							nil, nil, nil,
							/* 3*/ {zip: nil, data: []byte{1}},
						},
					}},
				},
			}},
		},
	}

	if !reflect.DeepEqual(got, want3) {
		t.Errorf("put() = %s, want %s", got, want3)
	}
}

func Test_mapping_get(t *testing.T) {
	got := &mapping{}
	got.put([]byte{1, 2, 3}, []byte{0})
	got.put([]byte{1, 2, 3, 4, 5}, []byte{1})
	got.put([]byte{1}, []byte{2})

	tests := []struct {
		key       []byte
		defaulted []byte
		val       []byte
		ok        bool
	}{
		{[]byte{}, nil, nil, false},
		{[]byte{1}, nil, []byte{2}, true},
		{[]byte{1, 2}, nil, []byte{2}, false},
		{[]byte{1, 2, 3}, nil, []byte{0}, true},
		{[]byte{1, 2, 3, 4}, nil, []byte{0}, false},
		{[]byte{1, 2, 3, 4, 5}, nil, []byte{1}, true},
		{[]byte{1, 2, 3, 4, 5, 6}, nil, []byte{1}, false},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, nil, []byte{1}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			val, ok := got.get(tt.key, nil)
			if ok != tt.ok {
				t.Errorf("get() = %v, want %v", ok, tt.ok)
			}
			if !bytes.Equal(val, tt.val) {
				t.Errorf("get() = %v, want %v", val, tt.val)
			}
		})
	}

}

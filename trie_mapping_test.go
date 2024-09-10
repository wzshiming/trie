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
	got := &Mapping[[]byte]{}
	got.put([]byte{1, 2, 3}, []byte{0})
	want1 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: []byte{2, 3}, data: []byte{0}, has: true},
		},
	}

	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.put([]byte{1, 2, 3, 4, 5}, []byte{1})
	want2 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: []byte{2, 3}, data: []byte{0}, has: true, mapping: &Mapping[[]byte]{
				array: [byteLength]*node[[]byte]{
					nil, nil, nil, nil,
					/* 4 */ {zip: []byte{5}, data: []byte{1}, has: true},
				}},
			},
		},
	}
	if !reflect.DeepEqual(got, want2) {
		t.Errorf("put() = %s, want %s", got, want2)
	}
	got.put([]byte{1, 2, 4, 5}, []byte{2})
	want3 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: []byte{2}, data: nil, mapping: &Mapping[[]byte]{
				array: [byteLength]*node[[]byte]{
					nil, nil, nil,
					/* 3 */ {zip: nil, data: []byte{0}, has: true, mapping: &Mapping[[]byte]{
						array: [byteLength]*node[[]byte]{
							nil, nil, nil, nil,
							/* 4 */ {zip: []byte{5}, data: []byte{1}, has: true},
						}},
					},
					/* 4 */ {zip: []byte{5}, data: []byte{2}, has: true},
				}},
			},
		},
	}

	if !reflect.DeepEqual(got, want3) {
		t.Errorf("put() = %s, want %s", got, want3)
	}
}

func Test_mapping_put_2(t *testing.T) {
	got := &Mapping[[]byte]{}
	got.put([]byte{1}, []byte{0})
	want1 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}, has: true},
		},
	}

	if !reflect.DeepEqual(got, want1) {
		t.Errorf("put() = %s, want %s", got, want1)
	}

	got.put([]byte{1, 2, 3}, []byte{1})
	want2 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}, has: true, mapping: &Mapping[[]byte]{
				array: [byteLength]*node[[]byte]{
					nil, nil,
					/* 2 */ {zip: []byte{3}, data: []byte{1}, has: true},
				},
			}},
		},
	}

	if !reflect.DeepEqual(got, want2) {
		t.Errorf("put() = %s, want %s", got, want2)
	}
	got.put([]byte{1, 2}, []byte{2})
	want3 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: nil, data: []byte{0}, has: true, mapping: &Mapping[[]byte]{
				array: [byteLength]*node[[]byte]{
					nil, nil,
					/* 2 */ {zip: nil, data: []byte{2}, has: true, mapping: &Mapping[[]byte]{
						array: [byteLength]*node[[]byte]{
							nil, nil, nil,
							/* 3 */ {zip: nil, data: []byte{1}, has: true},
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
	got := &Mapping[[]byte]{}
	got.put([]byte{1, 2, 3}, []byte{3})
	got.put([]byte{1, 2, 3, 4, 5}, []byte{5})
	got.put([]byte{1, 2, 3, 4, 5, 6}, []byte{6})
	got.put([]byte{1}, []byte{1})
	got.put([]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{8})
	got.put([]byte{1, 2, 3, 4, 5, 6, 7, 9, 9}, []byte{9})

	current7 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil, nil, nil, nil, nil, nil, nil,
			/* 7 */ {zip: nil, mapping: &Mapping[[]byte]{
				array: [byteLength]*node[[]byte]{
					nil, nil, nil, nil, nil, nil, nil, nil,
					{zip: nil, data: []byte{8}, has: true},
					{zip: []byte{9}, data: []byte{9}, has: true},
				},
			}},
		},
	}

	current6 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil, nil, nil, nil, nil, nil,
			/* 6 */ {zip: nil, data: []byte{6}, has: true, mapping: current7},
		},
	}
	current5 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil, nil, nil, nil,
			/* 4 */ {zip: []byte{5}, data: []byte{5}, has: true, mapping: current6},
		},
	}
	current3 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil, nil,
			/* 2 */ {zip: []byte{3}, data: []byte{3}, has: true, mapping: current5},
		},
	}
	current1 := &Mapping[[]byte]{
		array: [byteLength]*node[[]byte]{
			nil,
			/* 1 */ {zip: nil, data: []byte{1}, has: true, mapping: current3},
		},
	}
	tests := []struct {
		key       []byte
		defaulted []byte
		val       []byte
		current   *Mapping[[]byte]
		ok        bool
	}{
		{[]byte{}, nil, nil, nil, false},
		{[]byte{2}, nil, nil, nil, false},
		{[]byte{1}, nil, []byte{1}, current1, true},
		{[]byte{1, 3}, nil, []byte{1}, current1, true},
		{[]byte{1, 2}, nil, []byte{1}, current1, true},
		{[]byte{1, 2, 3}, nil, []byte{3}, current3, true},
		{[]byte{1, 2, 4}, nil, []byte{1}, current1, true},
		{[]byte{1, 2, 3, 4}, nil, []byte{3}, current3, true},
		{[]byte{1, 2, 3, 5}, nil, []byte{3}, current3, true},
		{[]byte{1, 2, 3, 4, 5}, nil, []byte{5}, current5, true},
		{[]byte{1, 2, 3, 4, 6}, nil, []byte{3}, current3, true},
		{[]byte{1, 2, 3, 4, 5, 6}, nil, []byte{6}, current6, true},
		{[]byte{1, 2, 3, 4, 5, 7}, nil, []byte{5}, current5, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, nil, []byte{6}, current6, true},
		{[]byte{1, 2, 3, 4, 5, 6, 8}, nil, []byte{6}, current6, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, nil, []byte{8}, current7.array[7].mapping, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil, []byte{8}, nil, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 9, 9}, nil, []byte{9}, current7.array[7].mapping, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.key), func(t *testing.T) {
			val, current, ok := got.get(nil, tt.key, nil, false)
			if ok != tt.ok {
				t.Errorf("get() ok = %v, want %v", ok, tt.ok)
			}
			if !reflect.DeepEqual(current, tt.current) {
				t.Errorf("get() current = %v, want %v", current, tt.current)
			}
			if !bytes.Equal(val, tt.val) {
				t.Errorf("get() val = %v, want %v", val, tt.val)
			}
		})
	}
}

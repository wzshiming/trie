package trie

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"
)

func BenchmarkTrie_Get1(b *testing.B) {
	mt := NewTrie[[]byte]()

	regs := []struct {
		key []byte
		val []byte
	}{
		{[]byte("GET /"), []byte("http")},
		{[]byte("DELETE /"), []byte("http")},
	}

	for _, reg := range regs {
		mt.Put(reg.key, reg.val)
	}

	d := []byte("DELETE /index")
	for i := 0; i != b.N; i++ {
		mt.Get(d)
	}
}

func BenchmarkTrie_Get2(b *testing.B) {
	mt := NewTrie[[]byte]()

	regs := []struct {
		key []byte
		val []byte
	}{
		{[]byte("PUT /"), []byte("http")},
		{[]byte("POST /"), []byte("http")},
	}

	for _, reg := range regs {
		mt.Put(reg.key, reg.val)
	}

	d := []byte("POST /index")
	for i := 0; i != b.N; i++ {
		mt.Get(d)
	}
}

func BenchmarkTrie_Put1(b *testing.B) {
	mt := NewTrie[[]byte]()
	data := []byte("http")
	var key [8]byte
	for i := 0; i != b.N; i++ {
		rand.Read(key[:])
		mt.Put(key[:], data)
	}
}

func TestTrie_GetAndPutAndKeys(t *testing.T) {
	mt := NewTrie[[]byte]()

	regs := []string{
		"AAAAAAAAAA",
		"AB",
		"ABC",
		"ABCD",
		"BCDE",
		"BCD",
		"BC",

		"AA",
		"AABB",
		"AABBCC",
		"AABBCCDD",
		"BBCCDDEE",
		"BBCCDD",
		"BBCC",
		"BB",

		"AABBCCDD",
	}

	for _, reg := range regs {
		tmp := []byte(reg)
		ok := mt.Put([]byte(reg), tmp)
		if !ok {
			t.Error(reg)
		}
	}

	got := mt.Keys()

	want := []string{
		"AAAAAAAAAA",
		"AB",
		"ABC",
		"ABCD",
		"BCDE",
		"BCD",
		"BC",
		"AA",
		"AABB",
		"AABBCC",
		"AABBCCDD",
		"BBCCDDEE",
		"BBCCDD",
		"BBCC",
		"BB",
	}

	if len(got) != len(want) {
		t.Errorf("Keys() len = %v want %v", len(got), len(want))
	}
	for _, vv := range got {
		t.Run("", func(t1 *testing.T) {
			got, _, _ := mt.Get([]byte(vv))
			if !bytes.Equal(got, []byte(vv)) {
				t.Errorf("Get() = %q, want %q", got, vv)
			}
		})
	}
}

func TestTrie_PutEmpty(t *testing.T) {
	mt := NewTrie[[]byte]()
	got := mt.Put(nil, nil)
	want := false
	if got != want {
		t.Errorf("Put(nil, nil) = %v want %v", got, want)
	}
}

func keyStrings(mt *Trie[int]) []string {
	var keys []string
	for _, k := range mt.Keys() {
		keys = append(keys, string(k))
	}
	return keys
}

func TestTrie_PutBranchNode(t *testing.T) {
	mt := NewTrie[int]()
	puts := []string{"ab", "ac", "a"}
	for i, key := range puts {
		if !mt.Put([]byte(key), i+1) {
			t.Errorf("Put(%q) = false want true", key)
		}
	}
	for i, key := range puts {
		if got, _, ok := mt.Get([]byte(key)); !ok || got != i+1 {
			t.Errorf("Get(%q) = %v, %v want %v, true", key, got, ok, i+1)
		}
	}
	if got, want := keyStrings(mt), []string{"a", "ab", "ac"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %q want %q", got, want)
	}
	if mt.Size() != 3 || mt.Depth() != 2 {
		t.Errorf("Size(), Depth() = %v, %v want 3, 2", mt.Size(), mt.Depth())
	}
}

func TestTrie_PutOverwrite(t *testing.T) {
	tests := []struct {
		puts  []string
		keys  []string
		depth int
	}{
		{[]string{"a"}, []string{"a"}, 1},
		{[]string{"abc"}, []string{"abc"}, 3},
		{[]string{"a", "ab"}, []string{"a", "ab"}, 2},
		{[]string{"ab", "ac", "a"}, []string{"a", "ab", "ac"}, 2},
	}
	for _, tt := range tests {
		key := tt.puts[len(tt.puts)-1]
		t.Run(key, func(t *testing.T) {
			mt := NewTrie[int]()
			for i, k := range tt.puts {
				mt.Put([]byte(k), i+1)
			}
			if !mt.Put([]byte(key), -1) {
				t.Errorf("Put(%q) = false want true", key)
			}
			if got, _, ok := mt.Get([]byte(key)); !ok || got != -1 {
				t.Errorf("Get(%q) = %v, %v want -1, true", key, got, ok)
			}
			if got := keyStrings(mt); !reflect.DeepEqual(got, tt.keys) {
				t.Errorf("Keys() = %q want %q", got, tt.keys)
			}
			if mt.Size() != len(tt.puts) || mt.Depth() != tt.depth {
				t.Errorf("Size(), Depth() = %v, %v want %v, %v", mt.Size(), mt.Depth(), len(tt.puts), tt.depth)
			}
		})
	}
}

func TestTrie_PutCopiesKey(t *testing.T) {
	mt := NewTrie[int]()
	buf := []byte("abc")
	mt.Put(buf, 1)
	buf[1] = 'X'
	if _, _, ok := mt.Get([]byte("abc")); !ok {
		t.Error(`Get("abc") = false want true`)
	}
	if _, _, ok := mt.Get([]byte("aXc")); ok {
		t.Error(`Get("aXc") = true want false`)
	}
	if got, want := keyStrings(mt), []string{"abc"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %q want %q", got, want)
	}
}

func TestTrie_PutCopiesKeyAfterSplit(t *testing.T) {
	mt := NewTrie[int]()
	buf := []byte("abcd")
	mt.Put(buf, 1)
	copy(buf, "abef")
	mt.Put(buf, 2)
	buf[3] = 'X'
	if _, _, ok := mt.Get([]byte("abef")); !ok {
		t.Error(`Get("abef") = false want true`)
	}
	if _, _, ok := mt.Get([]byte("abeX")); ok {
		t.Error(`Get("abeX") = true want false`)
	}
	if got, want := keyStrings(mt), []string{"abcd", "abef"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %q want %q", got, want)
	}
}

func TestTrie_Get(t *testing.T) {
	mt := NewTrie[[]byte]()

	regs := []struct {
		name  string
		magic string
	}{
		{"http", "GET /"},
		{"other", "GGG /"},
		{"http", "POST /"},
		{"http", "PUT /"},
		{"http", "DELETE /"},
		{"other", "GET //"},
		{"other", "GET "},
	}

	for _, reg := range regs {
		tmp := []byte(reg.name)
		mt.Put([]byte(reg.magic), tmp)
	}

	tests := []struct {
		want string
		data string
	}{
		{"http", "GET /"},
		{"other", "GGG /"},
		{"http", "POST /"},
		{"http", "PUT /"},
		{"http", "DELETE /"},
		{"other", "GET //"},
		{"other", "GET /////////////"},
		{"other", "GET //index"},
		{"other", "GET "},
		{"http", "GET /index"},
		{"http", "POST /index"},
		{"http", "PUT /index"},
		{"http", "DELETE /index"},
		{"", ""},
		{"", "GET"},
		{"", "POST"},
		{"", "PUT"},
		{"", "DELETE"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, _, _ := mt.Get([]byte(tt.data))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Get(%q) = %q want %q", tt.data, got, tt.want)
			}
		})
	}
}

func TestTrie_String(t *testing.T) {
	mt := NewTrie[[]byte]()

	regs := []struct {
		name  string
		magic string
	}{
		{"http", "GET /"},
		{"other", "GGG /"},
		{"http", "POST /"},
		{"http", "PUT /"},
		{"http", "DELETE /"},
		{"other", "GET //"},
		{"other", "GET "},
	}

	for _, reg := range regs {
		tmp := []byte(reg.name)
		mt.Put([]byte(reg.magic), tmp)
	}
	t.Log(mt.String())
}

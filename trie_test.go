package trie

import (
	"bytes"
	"testing"
)

func BenchmarkTrie_Get(b *testing.B) {
	mt := NewTrie()

	regs := []struct {
		key string
		val string
	}{
		{"GET /", "http"},
		{"DELETE /", "http"},
	}

	for _, reg := range regs {
		mt.Put([]byte(reg.key), []byte(reg.val))
	}

	d := []byte("DELETE /index")
	for i := 0; i != b.N; i++ {
		mt.Get(d)
	}
}

func TestTrie_GetAndPutAndKeys(t *testing.T) {
	mt := NewTrie()

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
		ok := mt.Put([]byte(reg), []byte(reg))
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

			got, _ := mt.Get([]byte(vv))
			if !bytes.Equal(got, []byte(vv)) {
				t.Errorf("Get() gotVal = %q, want %q", got, vv)
			}
		})
	}
}

func TestTrie_Get(t *testing.T) {
	mt := NewTrie()

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
		mt.Put([]byte(reg.magic), []byte(reg.name))
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
			got, _ := mt.Get([]byte(tt.data))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Error(mt.mapping.String())
				t.Errorf("Get(%q) = %q want %q", tt.data, got, tt.want)
			}
		})
	}
}

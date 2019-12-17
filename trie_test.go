package trie

import (
	"bytes"
	"crypto/rand"
	"sync"
	"testing"
)

func BenchmarkTrie_Get1(b *testing.B) {
	mt := NewTrie()

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
	mt := NewTrie()

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
	mt := NewTrie()
	data := []byte("http")
	var key [8]byte
	for i := 0; i != b.N; i++ {
		rand.Read(key[:])
		mt.Put(key[:], data)
	}
}

func BenchmarkTrie_Put2(b *testing.B) {
	mt := NewTrie()
	data := []byte("http")
	newFunc := func() interface{} {
		return make([]byte, 8)
	}
	pool := sync.Pool{
		New: newFunc,
	}
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	limit := make(chan struct{}, 10)
	for i := 0; i != b.N; i++ {
		limit <- struct{}{}
		go func() {
			buf := pool.Get().([]byte)
			defer pool.Put(buf)

			rand.Read(buf)
			mt.Put(buf, data)
			wg.Done()
			<-limit
		}()
	}
	wg.Wait()
}

func BenchmarkTrie_Put3(b *testing.B) {
	mt := NewTrie()
	key1 := []byte("key1")
	key2 := []byte("key2")
	data := []byte("http")
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i != b.N; i++ {
		go func() {
			mt.Get(key1)
			mt.Put(key1, data)
			mt.Get(key2)
			mt.Put(key2, data)
			wg.Done()
		}()
	}
	wg.Wait()
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

			got, _, _ := mt.Get([]byte(vv))
			if !bytes.Equal(got, []byte(vv)) {
				t.Errorf("Get() gotVal = %q, want %q", got, vv)
			}
		})
	}
}

func TestTrie_PutEmpty(t *testing.T) {
	mt := NewTrie()
	got := mt.Put(nil, nil)
	want := false
	if got != want {
		t.Errorf("Put(nil, nil) = %v want %v", got, want)
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
			got, _, _ := mt.Get([]byte(tt.data))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Get(%q) = %q want %q", tt.data, got, tt.want)
			}
		})
	}
}

func TestTrie_String(t *testing.T) {
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
	t.Log(mt.String())
}

package trie

const (
	byteLength = 256
)

// Trie is a trie tree implementation.
type Trie[T any] struct {
	mapping Mapping[T]
	depth   int
	size    int
}

// NewTrie returns a new Trie.
func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{}
}

// String returns format trie in a friendly.
func (t *Trie[T]) String() string {
	return t.mapping.String()
}

// Keys returns all key of Trie.
func (t *Trie[T]) Keys() [][]byte {
	out := make([][]byte, 0, t.size)
	t.Walk(func(k []byte, v T) {
		n := make([]byte, len(k))
		copy(n, k)
		out = append(out, n)
	})
	return out
}

// Walk calls f sequentially for each key and value present in the trie.
func (t *Trie[T]) Walk(f func(k []byte, v T)) {
	buf := make([]byte, 0, t.depth)
	t.mapping.walk(buf, f)
}

// Put sets the val in the trie for a key.
func (t *Trie[T]) Put(key []byte, val T) (finish bool) {
	finish = t.mapping.put(key, val)
	if !finish {
		return false
	}
	if t.depth < len(key) {
		t.depth = len(key)
	}
	t.size++
	return true
}

// Mapping gets the Mapping for get only.
func (t *Trie[T]) Mapping() (m *Mapping[T]) {
	return &t.mapping
}

// Get returns the val in the trie for a key.
func (t *Trie[T]) Get(key []byte) (val T, current *Mapping[T], finish bool) {
	return t.mapping.Get(key)
}

// Size returns the size of the trie.
func (t *Trie[T]) Size() int {
	return t.size
}

// Depth returns the depth of the trie.
func (t *Trie[T]) Depth() int {
	return t.depth
}

package trie

const (
	byteLength = 256
)

// Trie is a trie tree implementation.
type Trie struct {
	mapping
	deep int
	size int
}

// NewTrie returns a new Trie.
func NewTrie() *Trie {
	return &Trie{}
}

// String returns format trie in a friendly.
func (t *Trie) String() string {
	return t.mapping.String()
}

// Keys returns all key of Trie.
func (t *Trie) Keys() [][]byte {
	out := make([][]byte, 0, t.size)
	t.Walk(func(k, v []byte) {
		n := make([]byte, len(k))
		copy(n, k)
		out = append(out, n)
	})
	return out
}

// Walk calls f sequentially for each key and value present in the trie.
func (t *Trie) Walk(f func(k, v []byte)) {
	buf := make([]byte, 0, t.deep)
	t.mapping.walk(buf, f)
}

// Put sets the val in the trie for a key.
func (t *Trie) Put(key, val []byte) (finish bool) {
	finish = t.mapping.put(key, val)
	if !finish {
		return false
	}
	if t.deep < len(key) {
		t.deep = len(key)
	}
	t.size++
	return true
}

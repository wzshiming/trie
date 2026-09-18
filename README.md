# trie

[![GoDoc](https://godoc.org/github.com/wzshiming/trie?status.svg)](https://godoc.org/github.com/wzshiming/trie)
[![GitHub license](https://img.shields.io/github/license/wzshiming/trie.svg)](https://github.com/wzshiming/trie/blob/master/LICENSE)

Trie is a Compressed Prefix Tree Implementation in Golang Generic.

- `Put` is not safe to run concurrently with any other operation on the same trie.
- Once construction is finished, read-only operations on the trie (`Get`, `Walk`, `Keys`, `Size`, `Depth`, `String`, `MatchWithReader`) may run concurrently.
- `Put` copies the key; the caller keeps ownership of the passed slice.
- `Get` returns the value of the longest stored key that is a prefix of the lookup key.
- There is no deletion; removing keys means building a new trie.

## License

Licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/trie/blob/master/LICENSE) for the full license text.

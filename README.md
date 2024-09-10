# trie

[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/trie)](https://goreportcard.com/report/github.com/wzshiming/trie)
[![GoDoc](https://godoc.org/github.com/wzshiming/trie?status.svg)](https://godoc.org/github.com/wzshiming/trie)
[![GitHub license](https://img.shields.io/github/license/wzshiming/trie.svg)](https://github.com/wzshiming/trie/blob/master/LICENSE)

Trie is a Compressed Prefix Tree Implementation in Golang Generic.

``` console
> go test -benchmem -run=^$ -bench . github.com/wzshiming/trie -v
goos: linux
goarch: amd64
pkg: github.com/wzshiming/trie
cpu: Intel(R) Xeon(R) Platinum 8272CL CPU @ 2.60GHz
BenchmarkTrie_Get1
BenchmarkTrie_Get1-4    60400056                19.19 ns/op            0 B/op          0 allocs/op
BenchmarkTrie_Get2
BenchmarkTrie_Get2-4    64019403                19.07 ns/op            0 B/op          0 allocs/op
BenchmarkTrie_Put1
BenchmarkTrie_Put1-4     1234398               939.6 ns/op             0 B/op          0 allocs/op
BenchmarkTrie_Put2
BenchmarkTrie_Put2-4     1000000              1228 ns/op             339 B/op          3 allocs/op
BenchmarkTrie_Put3
BenchmarkTrie_Put3-4     3234558               398.5 ns/op            96 B/op          1 allocs/op
PASS
ok      github.com/wzshiming/trie       7.499s
```

## License

Licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/trie/blob/master/LICENSE) for the full license text.

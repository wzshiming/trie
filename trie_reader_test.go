package trie

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

type entry struct{ key, handler string }

type spyReader struct {
	io.Reader
	calls, empty int
}

func (spy *spyReader) Read(buf []byte) (int, error) {
	spy.calls++
	if len(buf) == 0 {
		spy.empty++
	}
	return spy.Reader.Read(buf)
}

func newTrie(entries []entry) *Trie[string] {
	mt := NewTrie[string]()
	for _, ent := range entries {
		mt.Put([]byte(ent.key), ent.handler)
	}
	return mt
}

func TestTrie_MatchWithReader(t *testing.T) {
	cmux := []entry{
		{"GET ", "http"}, {"HEAD ", "http"}, {"POST ", "http"}, {"PUT ", "http"},
		{"PRI * HTTP/2.0", "http2"}, {"SSH-", "ssh"},
		{"\x16\x03\x01", "tls"}, {"\x05\x01", "socks5"}, {"\x05\x02", "socks5"},
	}
	overlap := []entry{{"GET ", "get"}, {"GET /x", "getx"}}
	tests := []struct {
		name    string
		entries []entry
		reader  io.Reader
		handler string
		prefix  string
		err     error
		calls   int
	}{
		{"http truncated at depth", cmux, strings.NewReader("GET / HTTP/1.1\r\n"), "http", "GET / HTTP/1.1", nil, 1},
		{"exact GET no eof read", cmux, strings.NewReader("GET "), "http", "GET ", nil, 1},
		{"exact http2 no eof read", cmux, strings.NewReader("PRI * HTTP/2.0"), "http2", "PRI * HTTP/2.0", nil, 1},
		{"http2 byte at a time", cmux, iotest.OneByteReader(strings.NewReader("PRI * HTTP/2.0\r\nSM")), "http2", "PRI * HTTP/2.0", nil, 14},
		{"tls byte at a time", cmux, iotest.OneByteReader(bytes.NewReader([]byte{0x16, 0x03, 0x01, 0x00, 0xf5})), "tls", "\x16\x03\x01", nil, 3},
		{"socks5 byte at a time", cmux, iotest.OneByteReader(bytes.NewReader([]byte{0x05, 0x02, 0x00, 0x01})), "socks5", "\x05\x02", nil, 2},
		{"socks5 whole", cmux, bytes.NewReader([]byte{0x05, 0x01, 0x00}), "socks5", "\x05\x01\x00", nil, 1},
		{"ssh bytes with eof", cmux, iotest.DataErrReader(strings.NewReader("SSH-2.0-x")), "ssh", "SSH-2.0-x", nil, 1},
		{"ssh bytes with timeout", cmux, iotest.DataErrReader(io.MultiReader(strings.NewReader("SSH-2.0-x"), iotest.ErrReader(iotest.ErrTimeout))), "ssh", "SSH-2.0-x", nil, 1},
		{"exact http2 timeout reader", cmux, iotest.TimeoutReader(strings.NewReader("PRI * HTTP/2.0")), "http2", "PRI * HTTP/2.0", nil, 1},
		{"no match", cmux, strings.NewReader("XYZ"), "", "XYZ", ErrNotFound, 1},
		{"empty reader", cmux, strings.NewReader(""), "", "", io.EOF, 1},
		{"incomplete then eof", cmux, strings.NewReader("GET"), "", "GET", io.EOF, 2},
		{"incomplete with eof", cmux, iotest.DataErrReader(strings.NewReader("GET")), "", "GET", io.EOF, 1},
		{"incomplete then timeout", cmux, io.MultiReader(strings.NewReader("GET"), iotest.ErrReader(iotest.ErrTimeout)), "", "GET", iotest.ErrTimeout, 2},
		{"zero read terminates", cmux, io.MultiReader(strings.NewReader("GET"), iotest.ErrReader(nil)), "", "GET", ErrNotFound, 2},
		{"longer overlap", overlap, iotest.OneByteReader(strings.NewReader("GET /xyz")), "getx", "GET /x", nil, 6},
		{"shorter overlap diverges", overlap, iotest.OneByteReader(strings.NewReader("GET /yz")), "get", "GET /y", nil, 6},
		{"shorter overlap then eof", overlap, strings.NewReader("GET /"), "get", "GET /", io.EOF, 2},
		{"shorter overlap then timeout", overlap, io.MultiReader(strings.NewReader("GET /"), iotest.ErrReader(iotest.ErrTimeout)), "get", "GET /", iotest.ErrTimeout, 2},
		{"shorter overlap with timeout", overlap, iotest.DataErrReader(io.MultiReader(strings.NewReader("GET /"), iotest.ErrReader(iotest.ErrTimeout))), "get", "GET /", iotest.ErrTimeout, 1},
		{"shorter overlap then zero read", overlap, io.MultiReader(strings.NewReader("GET /"), iotest.ErrReader(nil)), "get", "GET /", nil, 2},
		{"zero handler", []entry{{"GET ", ""}}, strings.NewReader("GET /"), "", "GET ", nil, 1},
		{"empty trie", nil, strings.NewReader("GET "), "", "", ErrNotFound, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newTrie(tt.entries)
			spy := &spyReader{Reader: tt.reader}
			handler, prefix, err := mt.MatchWithReader(spy)
			if handler != tt.handler || string(prefix) != tt.prefix || err != tt.err {
				t.Errorf("MatchWithReader() = %q, %q, %v want %q, %q, %v", handler, prefix, err, tt.handler, tt.prefix, tt.err)
			}
			if spy.calls != tt.calls || spy.empty != 0 || len(prefix) > mt.Depth() {
				t.Errorf("reads = %d, empty reads = %d, len(prefix) = %d want %d, 0, <= %d", spy.calls, spy.empty, len(prefix), tt.calls, mt.Depth())
			}
		})
	}
}

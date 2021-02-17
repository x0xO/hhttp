package hhttp

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

type body struct {
	headers headers
	bytes   []byte
}

func (b body) String() string {
	return string(b.bytes)
}

func (b body) Bytes() []byte {
	return b.bytes
}

func (b body) Contains(pattern interface{}) bool {
	switch pattern.(type) {
	case []byte:
		return bytes.Contains(bytes.ToLower(b.bytes), bytes.ToLower(pattern.([]byte)))
	case string:
		return strings.Contains(strings.ToLower(b.String()), strings.ToLower(pattern.(string)))
	case *regexp.Regexp:
		return pattern.(*regexp.Regexp).Match(b.bytes)
	default:
		return false
	}
}

func (b body) UTF8() body {
	contentType := b.headers.Get("Content-Type")
	utf8Reader, err := charset.NewReader(bytes.NewReader(b.bytes), contentType)
	if err != nil {
		return b
	}

	utf8Body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return b
	}

	return body{b.headers, utf8Body}
}

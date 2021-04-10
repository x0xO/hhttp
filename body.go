package hhttp

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

type body struct {
	body    io.ReadCloser
	headers headers
	stream  *bufio.Reader
	content []byte
	limit   int64
	deflate bool
	utf8    bool
}

func (b *body) Stream() *bufio.Reader { return b.stream }

func (b *body) String() string { return string(b.Bytes()) }

func (b *body) Close() error { return b.body.Close() }

func (b *body) Limit(limit int64) *body { b.limit = limit; return b }

func (b *body) UTF8() string {
	var (
		content string
		err     error
	)

	if b.content != nil {
		content = string(b.content)
	} else {
		content = b.String()
	}

	utf8Reader, err := charset.NewReader(strings.NewReader(content), b.headers.Get("Content-Type"))
	if err != nil {
		return content
	}

	utf8Content, err := io.ReadAll(utf8Reader)
	if err != nil {
		return content
	}

	return string(utf8Content)
}

func (b *body) Bytes() []byte {
	if b.content != nil {
		return b.content
	}

	if b.stream != nil {
		b.body = io.NopCloser(b.stream)
	}

	defer b.Close()

	var err error
	if b.deflate {
		if b.body, err = zlib.NewReader(b.body); err != nil {
			return nil
		}
	}

	if b.limit != -1 {
		b.content, err = io.ReadAll(io.LimitReader(b.body, b.limit))
	} else {
		b.content, err = io.ReadAll(b.body)
	}

	if err != nil {
		return nil
	}

	return b.content
}

func (b *body) Contains(pattern interface{}) bool {
	switch pattern.(type) {
	case []byte:
		return bytes.Contains(bytes.ToLower(b.Bytes()), bytes.ToLower(pattern.([]byte)))
	case string:
		return strings.Contains(strings.ToLower(b.String()), strings.ToLower(pattern.(string)))
	case *regexp.Regexp:
		return pattern.(*regexp.Regexp).Match(b.Bytes())
	default:
		return false
	}
}

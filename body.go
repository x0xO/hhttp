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
	headers headers
	body    io.ReadCloser
	stream  *bufio.Reader
	deflate bool
}

func (b body) String() string { return string(b.Bytes()) }

func (b body) Stream() *bufio.Reader { return b.stream }

func (b body) Limit(limiter int64) body {
	return body{b.headers, io.NopCloser(io.LimitReader(b.body, limiter)), b.stream, b.deflate}
}

func (b *body) Bytes() []byte {
	if b.stream != nil {
		b.body = io.NopCloser(b.stream)
	}

	defer b.body.Close()

	if b.deflate {
		var err error
		if b.body, err = zlib.NewReader(b.body); err != nil {
			return nil
		}
	}

	bodyBytes, err := io.ReadAll(b.body)
	if err != nil {
		return nil
	}

	return bodyBytes
}

func (b body) Contains(pattern interface{}) bool {
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

func (b body) UTF8() body {
	contentType := b.headers.Get("Content-Type")
	utf8Reader, err := charset.NewReader(bytes.NewReader(b.Bytes()), contentType)
	if err != nil {
		return b
	}

	return body{b.headers, io.NopCloser(utf8Reader), b.stream, b.deflate}
}

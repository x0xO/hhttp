package hhttp

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

type body []byte

func (b body) String() string {
	return string(b)
}

func (b body) Contains(pattern interface{}) bool {
	switch pattern.(type) {
	case []byte:
		return bytes.Contains(bytes.ToLower(b), bytes.ToLower(pattern.([]byte)))
	case string:
		return strings.Contains(strings.ToLower(b.String()), strings.ToLower(pattern.(string)))
	case *regexp.Regexp:
		return pattern.(*regexp.Regexp).Match(b)
	default:
		return false
	}
}

func (b body) UTF8(data ...interface{}) body {
	var contentType string

	if len(data) != 0 {
		switch data[0].(type) {
		case headers:
			contentType = data[0].(headers).Get("Content-Type")
		case string:
			contentType = strings.ToLower(data[0].(string))
		}
	}

	utf8Reader, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return b
	}

	utf8Body, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return b
	}

	return utf8Body
}

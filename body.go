package hhttp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
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

func (b body) UTF8() body {
	peek, err := bufio.NewReader(bytes.NewReader(b)).Peek(1024)
	if err != nil {
		return b
	}

	e, _, _ := charset.DetermineEncoding(peek, "")
	bodyUTF8, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), e.NewDecoder()))
	if err != nil {
		return b
	}

	return bodyUTF8
}

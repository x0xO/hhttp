package hhttp

import (
	"bytes"
	"regexp"
	"strings"
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

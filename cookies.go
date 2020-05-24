package hhttp

import (
	"net/http"
	"regexp"
	"strings"
)

type cookies []*http.Cookie

func (c *cookies) Contains(pattern interface{}) bool {
	for _, cookie := range *c {
		switch pattern.(type) {
		case string:
			if strings.Contains(strings.ToLower(cookie.String()), strings.ToLower(pattern.(string))) {
				return true
			}
		case *regexp.Regexp:
			if pattern.(*regexp.Regexp).Match([]byte(cookie.String())) {
				return true
			}
		}
	}

	return false
}

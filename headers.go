package hhttp

import (
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
)

type headers http.Header

func (h headers) Contains(header string, patterns interface{}) bool {
	if h.Values(header) != nil {
		for _, value := range h.Values(header) {
			switch patterns.(type) {
			case string:
				if strings.Contains(strings.ToLower(value), strings.ToLower(patterns.(string))) {
					return true
				}
			case []string:
				for _, pattern := range patterns.([]string) {
					if strings.Contains(strings.ToLower(value), strings.ToLower(pattern)) {
						return true
					}
				}
			case []*regexp.Regexp:
				for _, pattern := range patterns.([]*regexp.Regexp) {
					if pattern.Match([]byte(value)) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (h headers) Values(key string) []string {
	return textproto.MIMEHeader(h).Values(key)
}

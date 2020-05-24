package hhttp

import (
	"net/http"
	"net/url"
)

type history []*http.Response

func (his history) URLS() []*url.URL {
	var URLs []*url.URL
	for _, h := range his {
		if h.Request.URL != nil {
			URLs = append(URLs, h.Request.URL)
		}
	}
	return URLs
}

func (his history) Referers() []string {
	var referers []string
	for _, h := range his {
		if h.Request.Referer() != "" {
			referers = append(referers, h.Request.Referer())
		}
	}
	return referers
}

func (his history) StatusCodes() []int {
	var statusCodes []int
	for _, h := range his {
		statusCodes = append(statusCodes, h.StatusCode)
	}
	return statusCodes
}

func (his history) Cookies() [][]*http.Cookie {
	var cookies [][]*http.Cookie
	for _, h := range his {
		if len(h.Cookies()) != 0 {
			cookies = append(cookies, h.Cookies())
		}
	}
	return cookies
}

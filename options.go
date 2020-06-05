package hhttp

import "time"

type Options struct {
	BasicAuth     interface{}
	DNS           string
	DNSoverTLS    *dnsOverTLS
	History       bool
	InterfaceAddr string
	MaxRedirect   int
	Proxy         interface{}
	Session       bool
	Timeout       time.Duration
	UserAgent     interface{}
}

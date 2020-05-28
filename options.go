package hhttp

import (
	"time"
)

type Options struct {
	BasicAuth   interface{}
	History     bool
	MaxRedirect int
	Proxy       interface{}
	Timeout     time.Duration
	UserAgent   interface{}
	IP          string
	DNS         string
	DNSoverTLS  *dnsOverTLS
}

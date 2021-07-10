package hhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type options struct {
	basicAuth      interface{}
	proxy          interface{}
	userAgent      interface{}
	dotResolver    *net.Resolver
	redirectPolicy func(*http.Request, []*http.Request) error
	interfaceAddr  string
	dns            string
	timeout        time.Duration
	maxRedirects   int
	session        bool
	history        bool
	http2          bool
	stream         bool
	keepAlive      bool
	retryMax       int
	retryWait      time.Duration
}

func NewOptions() *options { return &options{keepAlive: true} }

func (opt *options) BasicAuth(
	basicAuth interface{},
) *options {
	opt.basicAuth = basicAuth
	return opt
}

func (opt *options) DNS(dns string) *options { opt.dns = dns; return opt }

func (opt *options) DNSOverTLS() *dnsOverTLS { return &dnsOverTLS{opt: opt} }

func (opt *options) Timeout(timeout time.Duration) *options { opt.timeout = timeout; return opt }

func (opt *options) UserAgent(
	userAgent interface{},
) *options {
	opt.userAgent = userAgent
	return opt
}

func (opt *options) InterfaceAddr(addr string) *options { opt.interfaceAddr = addr; return opt }

func (opt *options) Proxy(proxy interface{}) *options { opt.proxy = proxy; return opt }

func (opt options) String() string { return fmt.Sprintf("%#v", opt) }

func (opt *options) KeepAlive(enable ...bool) *options {
	if len(enable) != 0 {
		opt.keepAlive = enable[0]
	}
	return opt
}

func (opt *options) Retry(retryMax int, retryWait ...time.Duration) *options {
	opt.retryWait = time.Second * 1
	if len(retryWait) != 0 {
		opt.retryWait = retryWait[0]
	}
	opt.retryMax = retryMax
	return opt
}

func (opt *options) Stream(enable ...bool) *options {
	if len(enable) != 0 {
		opt.stream = enable[0]
	} else {
		opt.stream = true
	}
	return opt
}

func (opt *options) History(enable ...bool) *options {
	if len(enable) != 0 {
		opt.history = enable[0]
	} else {
		opt.history = true
	}
	return opt
}

func (opt *options) HTTP2(enable ...bool) *options {
	if len(enable) != 0 {
		opt.http2 = enable[0]
	} else {
		opt.http2 = true
	}
	return opt
}

func (opt *options) RedirectPolicy(f func(*http.Request, []*http.Request) error) *options {
	opt.redirectPolicy = f
	return opt
}

func (opt *options) MaxRedirects(maxRedirects int) *options {
	opt.maxRedirects = maxRedirects
	return opt
}

func (opt *options) Session(enable ...bool) *options {
	if len(enable) != 0 {
		opt.session = enable[0]
	} else {
		opt.session = true
	}
	return opt
}

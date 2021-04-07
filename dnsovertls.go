package hhttp

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"time"
)

type dnsOverTLS struct{ opt *options }

func (dot *dnsOverTLS) Google() *options {
	return dot.AddProvider("dns.google", "8.8.8.8:853", "8.8.4.4:853")
}

func (dot *dnsOverTLS) Cloudflare() *options {
	return dot.AddProvider("cloudflare-dns.com", "1.1.1.1:853", "1.0.0.1:853")
}

func (dot *dnsOverTLS) Libredns() *options {
	return dot.AddProvider("dot.libredns.gr", "116.202.176.26:853")
}

func (dot *dnsOverTLS) Quad9() *options {
	return dot.AddProvider("dns.quad9.net", "9.9.9.9:853", "149.112.112.112:853")
}

func (dot *dnsOverTLS) Switch() *options {
	return dot.AddProvider("dns.switch.ch", "130.59.31.248:853", "130.59.31.251:853")
}

func (dot dnsOverTLS) resolver(serverName string, addresses ...string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial:     dial(serverName, addresses...),
	}
}

func (dot *dnsOverTLS) AddProvider(serverName string, addresses ...string) *options {
	dot.opt.dotResolver = dot.resolver(serverName, addresses...)
	return dot.opt
}

func dial(serverName string, addresses ...string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		conn, err := dialer.DialContext(ctx, "tcp", addresses[rand.Intn(len(addresses))])
		if err != nil {
			return nil, err
		}
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute)
		return tls.Client(conn, &tls.Config{
			ServerName:         serverName,
			ClientSessionCache: tls.NewLRUClientSessionCache(0),
		}), nil
	}
}

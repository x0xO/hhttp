package hhttp

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"time"
)

type dnsOverTLS struct {
	dnsResolver *net.Resolver
}

func DNSoverTLS() *dnsOverTLS { return &dnsOverTLS{} }

func (dot *dnsOverTLS) Google() *dnsOverTLS {
	dot.dnsResolver = dot.resolver("dns.google", "8.8.8.8:853", "8.8.4.4:853")
	return dot
}

func (dot *dnsOverTLS) Cloudflare() *dnsOverTLS {
	dot.dnsResolver = dot.resolver("cloudflare-dns.com", "1.1.1.1:853", "1.0.0.1:853")
	return dot
}

func (dot *dnsOverTLS) Libredns() *dnsOverTLS {
	dot.dnsResolver = dot.resolver("dot.libredns.gr", "116.202.176.26:853")
	return dot
}

func (dot *dnsOverTLS) Quad9() *dnsOverTLS {
	dot.dnsResolver = dot.resolver("dns.quad9.net", "9.9.9.9:853", "149.112.112.112:853")
	return dot
}

func (dot dnsOverTLS) resolver(serverName string, addresses ...string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial:     dial(serverName, addresses...),
	}
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

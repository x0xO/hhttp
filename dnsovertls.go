package hhttp

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"time"
)

type dnsOverTLS struct {
	dns         string
	addrs       []string
	dnsResolver *net.Resolver
}

func DNSoverTLS() *dnsOverTLS {
	return &dnsOverTLS{}
}

func (dot *dnsOverTLS) Google() *dnsOverTLS {
	dot.dns = "dns.google"
	dot.addrs = []string{"8.8.8.8:853", "8.8.4.4:853"}
	dot.dnsResolver = dot.resolver()
	return dot
}

func (dot *dnsOverTLS) Cloudflare() *dnsOverTLS {
	dot.dns = "cloudflare-dns.com"
	dot.addrs = []string{"1.1.1.1:853", "1.0.0.1:853"}
	dot.dnsResolver = dot.resolver()
	return dot
}

func (dot *dnsOverTLS) Libredns() *dnsOverTLS {
	dot.dns = "dot.libredns.gr"
	dot.addrs = []string{"116.202.176.26:853"}
	dot.dnsResolver = dot.resolver()
	return dot
}

func (dot *dnsOverTLS) Quad9() *dnsOverTLS {
	dot.dns = "dns.quad9.net"
	dot.addrs = []string{"9.9.9.9:853", "149.112.112.112:853"}
	dot.dnsResolver = dot.resolver()
	return dot
}

func (dot dnsOverTLS) resolver() *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial:     dial(dot.dns, dot.addrs...),
	}
}

func dial(serverName string, addrs ...string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		conn, err := dialer.DialContext(ctx, "tcp", addrs[rand.Intn(len(addrs))])
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

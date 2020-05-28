package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.Options{DNSoverTLS: hhttp.DNSoverTLS().Google()}
	// opt := hhttp.Options{DNSoverTLS: hhttp.DNSoverTLS().Cloudflare()}
	// opt := hhttp.Options{DNSoverTLS: hhttp.DNSoverTLS().Libredns()}
	// opt := hhttp.Options{DNSoverTLS: hhttp.DNSoverTLS().Quad9()}

	r, err := hhttp.NewClient().SetOptions(&opt).Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

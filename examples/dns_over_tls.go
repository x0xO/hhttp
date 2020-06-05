package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.NewOptions()

	opt.DNSOverTLS().Google()
	// opt.DNSOverTLS().Cloudflare()
	// opt.DNSOverTLS().Libredns()
	// opt.DNSOverTLS().Quad9()

	r, err := hhttp.NewClient().SetOptions(opt).Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

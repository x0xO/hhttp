package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.NewOptions()

	opt.DNS("1.1.1.1:53") // cloudflare dns
	// opt.DNS("127.0.0.1:9053") // tor dns
	// opt.DNS("8.8.8.8:53")     // google dns
	// opt.DNS("9.9.9.9:53")     // quad9 dns

	r, err := hhttp.NewClient().SetOptions(opt).Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

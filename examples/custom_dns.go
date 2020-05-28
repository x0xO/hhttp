package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	// opt := hhttp.Options{DNS: "127.0.0.1:9053"} // tor dns
	// opt := hhttp.Options{DNS: "8.8.8.8:53"} // google dns
	// opt := hhttp.Options{DNS: "1.1.1.1:53"} // cloudflare dns
	opt := hhttp.Options{DNS: "9.9.9.9:53"} // quad9 dns

	r, err := hhttp.NewClient().SetOptions(&opt).Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

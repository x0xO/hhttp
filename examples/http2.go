package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().HTTP2()).Get("https://http2cdn.cdnsun.com").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Proto)
	fmt.Println(r.Body)

	r.Debug()
}

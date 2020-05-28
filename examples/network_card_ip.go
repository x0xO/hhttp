package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.Options{IP: "127.0.0.1"} // network adapter ip address

	r, err := hhttp.NewClient().SetOptions(&opt).Get("http://myip.dnsomatic.com").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

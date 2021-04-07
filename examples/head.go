package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().Head("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Status)
}

package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().SetOptions(&hhttp.Options{Timeout: 2}).Get("httpbingo.org/delay/3").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

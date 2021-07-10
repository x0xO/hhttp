package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.NewOptions().Retry(5)
	// opt := hhttp.NewOptions().Retry(3, time.Millisecond*50)
	r, err := hhttp.NewClient().SetOptions(opt).Get("http://httpbingo.org/unstable").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
	fmt.Println(r.Attempts)
}

package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().Timeout(1)).Get("httpbingo.org/delay/5").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

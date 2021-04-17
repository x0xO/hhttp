package main

import (
	"fmt"
	"log"
	"time"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().Timeout(time.Second)).Get("httpbingo.org/delay/2").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

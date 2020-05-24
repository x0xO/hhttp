package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	type Get struct {
		Headers struct {
			UserAgent []string `json:"User-Agent"`
		} `json:"headers"`
	}

	r, err := hhttp.NewClient().Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	var get Get
	r.JSON(&get)

	fmt.Println(get.Headers.UserAgent)
	fmt.Println(r.UserAgent)
}

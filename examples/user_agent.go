package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	type Get struct {
		Headers struct {
			UserAgent []string `json:"User-Agent"`
		} `json:"headers"`
	}

	URL := "https://httpbingo.org/get"

	r, _ := hhttp.NewClient().Get(URL).Do()

	var get Get
	r.JSON(&get)

	fmt.Printf("default user agent: %s\n", get.Headers.UserAgent)

	// change user-agent header
	options := hhttp.Options{UserAgent: "From Root with love!!!"}

	r, _ = hhttp.NewClient().SetOptions(&options).Get(URL).Do()

	r.JSON(&get)

	fmt.Printf("changed user agent: %s\n", get.Headers.UserAgent)
	fmt.Println(r.UserAgent)
}

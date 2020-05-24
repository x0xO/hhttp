package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	type Headers struct {
		Headers struct {
			Referer   []string `json:"Referer"`
			UserAgent []string `json:"User-Agent"`
		} `json:"headers"`
	}

	URL := "https://httpbingo.org/headers"

	h1 := map[string]string{"Referer": "Hell"}
	h2 := map[string]string{"Referer": "Paradise"}

	r, err := hhttp.NewClient().Get(URL).SetHeaders(h1).AddHeaders(h2).Do()
	if err != nil {
		log.Fatal(err)
	}

	var headers Headers
	r.JSON(&headers)

	fmt.Println(headers.Headers.Referer)
	fmt.Println(r.Referer()) // return first only

	fmt.Println(r.Headers)
	fmt.Println(r.Headers.Values("date"))
}

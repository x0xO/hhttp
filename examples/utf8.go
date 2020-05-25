package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	r, _ := hhttp.NewClient().Get("https://httpbingo.org/encoding/utf8").Do()
	fmt.Println(r.Body)

	r, _ = hhttp.NewClient().Get("http://vk.com").Do()

	fmt.Println(r.Body.UTF8())          // by body
	fmt.Println(r.Body.UTF8(r.Headers)) // by headers and body if charset not not found in headers
}

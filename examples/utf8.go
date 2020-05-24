package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	r, _ := hhttp.NewClient().Get("https://httpbingo.org/encoding/utf8").Do()
	fmt.Println(r.Body)
}

package main

import (
	"fmt"
	"time"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "https://httpbingo.org/cache"
	r, _ := hhttp.NewClient().
		Get(URL).
		AddHeaders(map[string]string{"If-Modified-Since": time.Now().Format("02.01.2006-15:04:05")}).
		Do()

	fmt.Println(r.StatusCode)
	r.Debug()
}

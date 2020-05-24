package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "https://httpbingo.org/get"

	cli := hhttp.NewClient()
	req := cli.Get(URL)

	resp, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)
	fmt.Println(resp.Cookies)
	fmt.Println(resp.Headers)
	fmt.Println(resp.URL)
}

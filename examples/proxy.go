package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	type Proxy struct {
		ISTor bool   `json:"IsTor"`
		IP    string `json:"IP"`
	}

	URL := "https://check.torproject.org/api/ip"
	r, err := hhttp.NewClient().SetOptions(&hhttp.Options{Proxy: "socks5://127.0.0.1:9050"}).Get(URL).Do()
	// // for random select proxy from slice
	// r, err := hhttp.NewClient().SetOptions(&hhttp.Options{Proxy: []string{"socks5://127.0.0.1:9050", "socks5://127.0.0.1:9050"}}).Get(URL).Do()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(r.Body)

	var proxy Proxy
	r.JSON(&proxy)

	fmt.Printf("is tor: %v, ip: %s", proxy.ISTor, proxy.IP)
}

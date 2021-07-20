package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {

	type Ja3 struct {
		Ja3Hash   string `json:"ja3_hash"`
		Ja3       string `json:"ja3"`
		UserAgent string `json:"User-Agent"`
	}

	opt := hhttp.NewOptions()

	// opt := hhttp.NewOptions().DNS("8.8.8.8:53")
	// opt := hhttp.NewOptions().DNSOverTLS().Google()

	// opt := hhttp.NewOptions().Proxy("socks5://127.0.0.1:9050")
	// opt := hhttp.NewOptions().Proxy("http://127.0.0.1:18080")

	opt.TLSFingerprint().
		JA3("771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-22-24,0").
		UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.154 Chrome/83.0.4103.122 Electron/9.3.5 Safari/537.36")

	r, err := hhttp.NewClient().SetOptions(opt).Get("https://ja3er.com/json").Do()
	if err != nil {
		log.Fatal(err)
	}

	var ja3 Ja3
	r.JSON(&ja3)

	fmt.Println(ja3.Ja3Hash == "34a0de67a0bdb5aab4df9962613cf620")

	fmt.Println(ja3.Ja3Hash)
	fmt.Println(ja3.Ja3)
	fmt.Println(ja3.UserAgent)

}

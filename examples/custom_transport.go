package main

import (
	"log"
	"time"

	"github.com/x0xO/hhttp"
)

func main() {
	// transport custom settings
	cli := hhttp.NewClient()

	// cli.GetTransport().TLSHandshakeTimeout = time.Nanosecond

	tr := cli.GetTransport()
	tr.TLSHandshakeTimeout = time.Nanosecond

	_, err := cli.Get("https://google.com").Do()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/x0xO/hhttp"
)

func main() {

	r, err := hhttp.NewClient().
		SetOptions(hhttp.NewOptions().KeepAlive(false)).
		Get("http://www.keycdn.com").
		Do()

	if err != nil {
		log.Fatal(err)
	}

	r.Debug().Print() // Connection: close
}

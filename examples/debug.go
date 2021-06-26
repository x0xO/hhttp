package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	r, err := hhttp.NewClient().Get("https://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	r.Debug().Print()     // without body
	r.Debug(true).Print() // with body

	fmt.Println(r.Time)
}

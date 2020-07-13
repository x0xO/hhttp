package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/x0xO/hhttp"
)

func main() {
	opt := hhttp.NewOptions().RedirectPolicy(
		func(req *http.Request, via []*http.Request) error {
			if len(via) >= 4 {
				return fmt.Errorf("stopped after %d redirects", 4)
			}
			return nil
		},
	)

	r, err := hhttp.NewClient().SetOptions(opt).Get("https://httpbingo.org/redirect/6").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

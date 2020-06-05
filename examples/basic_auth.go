package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	type basicAuth struct {
		Authorized bool   `json:"authorized"`
		User       string `json:"user"`
	}

	URL := "https://httpbingo.org/basic-auth/root/passwd"

	r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().BasicAuth("root:passwd")).Get(URL).Do()
	// r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().BasicAuth([]string{"root", "passwd"})).Get(URL).Do()
	// r, err := hhttp.NewClient().SetOptions(hhttp.NewOptions().BasicAuth(map[string]string{"root": "passwd"})).Get(URL).Do()
	if err != nil {
		log.Fatal(err)
	}

	var ba basicAuth
	r.JSON(&ba)

	fmt.Printf("authorized: %v, user: %s", ba.Authorized, ba.User)
}

package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "https://httpbingo.org/cookies"

	r, err := hhttp.NewClient().Get(URL + "/set?name1=value1&name2=value2").Do()
	if err != nil {
		log.Fatal(err)
	}

	// check cookies
	log.Println(r.GetCookies(URL))

	// second request, returns cookie data
	r, _ = r.Session.Get(URL).Do()

	// check if cookies in response {"name1":"value1","name2":"value2"}
	fmt.Println(r.Body)
}

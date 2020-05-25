package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "https://httpbingo.org/cookies"

	// example 1
	// chains session
	r, _ := hhttp.NewClient().Get(URL + "/set?name1=value1&name2=value2").Do()
	r, _ = r.Session.Get(URL).Do()
	fmt.Println(r.Body) // check if cookies in response {"name1":"value1","name2":"value2"}

	// example 2
	// splited session
	cli := hhttp.NewClient()
	cli.Get(URL + "/set?name1=value1&name2=value2").Do()
	s, _ := cli.Get(URL).Do()
	fmt.Println(s.Body) // check if cookies in response {"name1":"value1","name2":"value2"}
}

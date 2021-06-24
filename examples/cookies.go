package main

import (
	"fmt"
	"net/http"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "http://google.com"

	// cookie before request
	c1 := http.Cookie{Name: "root1", Value: "cookie1"}
	c2 := http.Cookie{Name: "root2", Value: "cookie2"}

	r, _ := hhttp.NewClient().
		SetOptions(hhttp.NewOptions().Session()).
		Get(URL).
		AddCookies(c1, c2).
		Do()

	r.Debug()

	// set cookie after first request
	r.SetCookies(URL, []*http.Cookie{{Name: "root", Value: "cookie"}})

	r, _ = r.Get(URL).Do()
	r.Debug()

	fmt.Println(r.GetCookies(URL)) // request url cookies
	fmt.Println(r.Cookies)
}

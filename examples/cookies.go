package main

import (
	"fmt"
	"net/http"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "http://google.com"

	r, _ := hhttp.NewClient().Get(URL).Do()

	r.SetCookie(URL, []*http.Cookie{{Name: "root", Value: "cookies"}})

	r, _ = r.Session.Get(URL).Do()
	r.Debug()

	fmt.Println(r.GetCookies(URL)) // request url cookies
	fmt.Println(r.Cookies)
}

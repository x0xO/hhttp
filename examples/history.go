package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	r, _ := hhttp.NewClient().SetOptions(hhttp.NewOptions().History()).Get("http://google.com").Do()

	fmt.Println(r.History.Referers())
	fmt.Println(r.History.StatusCodes())
	fmt.Println(r.History.Cookies())
	fmt.Println(r.History.URLS())
}

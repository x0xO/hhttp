package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	URL := "https://httpbin.org/gzip"
	r, _ := hhttp.NewClient().Get(URL).Do()
	fmt.Println(r.Body)

	URL = "https://httpbin.org/deflate"
	r, _ = hhttp.NewClient().Get(URL).Do()
	fmt.Println(r.Body)

	// resp, err := http.Get(URL)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// body, _ := ioutil.ReadAll(resp.Body)
	//
	// fmt.Println(string(body))
}

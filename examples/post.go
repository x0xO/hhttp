package main

import (
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {
	type Post struct {
		Form struct {
			Custemail []string `json:"custemail"`
			Custname  []string `json:"custname"`
			Custtel   []string `json:"custtel"`
		} `json:"form"`
	}

	URL := "https://httpbingo.org/post"

	// string post data
	data := "custname=root&custtel=999999999&custemail=some@email.com"

	r, _ := hhttp.NewClient().Post(URL, data).Do()

	var post Post
	r.JSON(&post)

	fmt.Println(post.Form.Custname)
	fmt.Println(post.Form.Custtel)
	fmt.Println(post.Form.Custemail)

	// map post data
	mapData := map[string]string{
		"custname":  "toor",
		"custtel":   "88888888",
		"custemail": "rest@gmail.com",
	}

	r, _ = hhttp.NewClient().Post(URL, mapData).Do()

	r.JSON(&post)

	fmt.Println(post.Form.Custname)
	fmt.Println(post.Form.Custtel)
	fmt.Println(post.Form.Custemail)
}

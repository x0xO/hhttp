package main

import (
	"github.com/x0xO/hhttp"
)

func main() {
	URL := "http://ptsv2.com/t/ys04l-1590171554/post"

	hhttp.NewClient().FileUpload(URL, "filefield", "/path/to/file.txt").Do()             // with file path
	hhttp.NewClient().FileUpload(URL, "filefield", "info.txt", "Hello from hhttp!").Do() // without phisical file

	// with multipart data
	multipartValues := map[string]string{"some": "values"}
	hhttp.NewClient().FileUpload(URL, "filefield", "/path/to/file.txt", multipartValues).Do()             // with file path
	hhttp.NewClient().FileUpload(URL, "filefield", "info.txt", "Hello from hhttp!", multipartValues).Do() // without phisical file
}

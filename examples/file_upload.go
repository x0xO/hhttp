package main

import (
	"github.com/x0xO/hhttp"
)

func main() {
	URL := "http://ptsv2.com/t/ys04l-1590171554/post"

	// with file path
	hhttp.NewClient().
		FileUpload(URL, "filefield", "/path/to/file.txt").
		Do()

	// without phisical file
	hhttp.NewClient().
		FileUpload(URL, "filefield", "info.txt", "Hello from hhttp!").
		Do()

	// with multipart data
	multipartValues := map[string]string{"some": "values"}

	// with file path
	hhttp.NewClient().
		FileUpload(URL, "filefield", "/path/to/file.txt", multipartValues).
		Do()

	// without phisical file
	hhttp.NewClient().
		FileUpload(URL, "filefield", "info.txt", "Hello from hhttp!", multipartValues).
		Do()
}

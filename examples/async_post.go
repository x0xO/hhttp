package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/x0xO/hhttp"
)

func main() {
	type MultiPost struct {
		Form struct {
			Comments  []string `json:"comments"`
			Custemail []string `json:"custemail"`
			Custname  []string `json:"custname"`
		} `json:"form"`
	}

	var URLs []string
	for i := 0; i < 50; i++ {
		URLs = append(URLs, "https://httpbingo.org/post")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data := "custname=root&custtel=999999999&custemail=some@email.com"

	var post MultiPost

	// with defaultMaxWorkers limited to 10 requests, no context
	// jobs, errors := hhttp.NewClient().Async.Post(URLs, data).Do()

	// one
	// with context and pool worker limited to 5 requests
	jobs, errors := hhttp.NewClient().Async.WithContext(ctx).Post(URLs, data).Pool(5).Do()

	for jobs != nil && errors != nil {
		select {
		case job, ok := <-jobs:
			if !ok {
				jobs = nil
				continue
			}
			job.JSON(&post)
			if post.Form.Custname[0] == "root" {
				fmt.Println("FOUND")
				cancel()
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
				continue
			}
			fmt.Println(err)
		}
	}

	fmt.Println(strings.Repeat("=", 80))

	// two
	// with custom pool worker limited to 100 requests
	jobs, errors = hhttp.NewClient().Async.Post(URLs, data).Pool(100).Do()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for job := range jobs {
			job.JSON(&post)
			fmt.Println(post.Form.Custname)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range errors {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("FINISH")
}

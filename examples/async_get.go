package main

import (
	"context"
	"fmt"

	"github.com/x0xO/hhttp"
)

func main() {

	var URLs []string
	for i := 0; i < 100; i++ {
		URLs = append(URLs, "https://httpbingo.org/get")
	}

	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	/* == URLS CHAN START == */
	urlsChan := make(chan string)
	go func() {
		defer close(urlsChan)
		for _, URL := range URLs {
			urlsChan <- URL
		}
	}()

	jobs, errors := hhttp.NewClient().
		Async.WithContext(ctx).
		Get(urlsChan). // urls chan string
		Pool(20).
		Do()

	/* == URLS CHAN END == */

	// with context and pool worker, limit to 20 requests

	// jobs, errors := hhttp.NewClient().
	// 	Async.WithContext(ctx).
	// 	Get(URLs). //urls []string
	// 	Pool(20).
	// 	Do()

	for jobs != nil && errors != nil {
		select {
		case job, ok := <-jobs:
			if !ok {
				jobs = nil
				continue
			}
			if job.Body.Contains("httpbingo") {
				cancel() // stop gorutines
				fmt.Println("FOUND")
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
				continue
			}
			fmt.Println(err)
		}
	}

	// var wg sync.WaitGroup
	// wg.Add(2)
	//
	// go func() {
	// 	defer wg.Done()
	// 	for job := range jobs {
	// 		if job.Body.Contains("google") {
	// 			cancel() // stop gorutines
	// 			fmt.Println("FOUND")
	// 		}
	// 	}
	// }()
	//
	// go func() {
	// 	defer wg.Done()
	// 	for err := range errors {
	// 		fmt.Println(err)
	// 	}
	// }()
	//
	// wg.Wait()

	fmt.Println("FINISH")
}

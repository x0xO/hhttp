package main

import (
	"fmt"
	"sync"

	"github.com/x0xO/hhttp"
)

func main() {
	var URLs []string
	for i := 0; i < 10; i++ {
		URLs = append(URLs, "https://httpbingo.org/get")
	}

	type Get struct {
		Headers struct {
			UserAgent []string `json:"User-Agent"`
		} `json:"headers"`
	}

	opt := hhttp.NewOptions().UserAgent([]string{"one", "two", "three", "four", "five"})

	jobs, errors := hhttp.NewClient().SetOptions(opt).Async.Get(URLs).Do()

	var wg sync.WaitGroup
	wg.Add(2)

	var get Get

	go func() {
		defer wg.Done()
		for job := range jobs {
			job.JSON(&get)
			fmt.Println(get.Headers.UserAgent)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range errors {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	fmt.Println("FINISH")
}

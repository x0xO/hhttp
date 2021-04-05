package hhttp

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/x0xO/hhttp/hsyscall"
)

type Requests struct {
	jobs       chan *Request
	setHeaders map[string]string
	addHeaders map[string]string
	maxWorkers int
}

func (reqs *Requests) Do() (chan *Response, chan error) {
	maxWorkers := defaultMaxWorkers

	if reqs.maxWorkers != 0 {
		if runtime.GOOS != "windows" {
			reqs.maxWorkers = hsyscall.RlimitStack(reqs.maxWorkers)
		}
		maxWorkers = reqs.maxWorkers
		fmt.Println(maxWorkers)
	}

	results := make(chan *Response)
	errors := make(chan error)

	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range reqs.jobs {
				if resp, err := job.SetHeaders(reqs.setHeaders).AddHeaders(reqs.addHeaders).Do(); err != nil {
					errors <- err
				} else {
					results <- resp
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	return results, errors
}

func (reqs *Requests) Pool(workers int) *Requests {
	reqs.maxWorkers = workers
	return reqs
}

func (reqs *Requests) SetHeaders(headers map[string]string) *Requests {
	reqs.setHeaders = headers
	return reqs
}

func (reqs *Requests) AddHeaders(headers map[string]string) *Requests {
	reqs.addHeaders = headers
	return reqs
}

package hhttp

import (
	"sync"
	"syscall"
)

type Requests struct {
	jobs       chan *Request
	maxWorkers int
	setHeaders map[string]string
	addHeaders map[string]string
}

func (reqs *Requests) Do() (chan *Response, chan error) {
	maxWorkers := defaultMaxWorkers

	if reqs.maxWorkers != 0 {
		var rLimit syscall.Rlimit
		syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if uint64(reqs.maxWorkers) > rLimit.Cur {
			reqs.maxWorkers = int(float64(rLimit.Cur) * 0.7)
		}
		maxWorkers = reqs.maxWorkers
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

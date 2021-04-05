package hhttp

import "context"

type async struct {
	client *Client
	ctx    context.Context
}

func (a *async) WithContext(ctx context.Context) *async {
	a.ctx = ctx
	return a
}

func (a *async) Get(URLS interface{}, data ...interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						URLS = nil
						return
					default:
						jobs <- a.client.Get(URL, data...)
					}
					continue
				}
				jobs <- a.client.Get(URL, data...)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Get(URL, data...)
					}
					continue
				}
				jobs <- a.client.Get(URL, data...)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Delete(URLS interface{}, data ...interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Delete(URL, data...)
					}
					continue
				}
				jobs <- a.client.Delete(URL, data...)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Delete(URL, data...)
					}
					continue
				}
				jobs <- a.client.Delete(URL, data...)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Head(URLS interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Head(URL)
					}
					continue
				}
				jobs <- a.client.Head(URL)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Head(URL)
					}
					continue
				}
				jobs <- a.client.Head(URL)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Post(URLS interface{}, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Post(URL, data)
					}
					continue
				}
				jobs <- a.client.Post(URL, data)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Post(URL, data)
					}
					continue
				}
				jobs <- a.client.Post(URL, data)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Put(URLS interface{}, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Put(URL, data)
					}
					continue
				}
				jobs <- a.client.Put(URL, data)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.Put(URL, data)
					}
					continue
				}
				jobs <- a.client.Put(URL, data)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) PostFile(URLS interface{}, fieldName, filePath string, uploadForm map[string]string) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		switch URLS.(type) {
		case chan string:
			for URL := range URLS.(chan string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.PostFile(URL, fieldName, filePath, uploadForm)
					}
					continue
				}
				jobs <- a.client.PostFile(URL, fieldName, filePath, uploadForm)
			}
		case []string:
			for _, URL := range URLS.([]string) {
				if a.ctx != nil {
					select {
					case <-a.ctx.Done():
						return
					default:
						jobs <- a.client.PostFile(URL, fieldName, filePath, uploadForm)
					}
					continue
				}
				jobs <- a.client.PostFile(URL, fieldName, filePath, uploadForm)
			}
		}
	}()

	return &Requests{jobs: jobs}
}

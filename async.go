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

func (a *async) Get(URLS []string, data ...interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Delete(URLS []string, data ...interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Head(URLS []string) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)

		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Post(URLS []string, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)
		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

func (a *async) PostJSON(URLS []string, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)
		for _, URL := range URLS {
			if a.ctx != nil {
				select {
				case <-a.ctx.Done():
					return
				default:
					jobs <- a.client.PostJSON(URL, data)
				}
				continue
			}
			jobs <- a.client.PostJSON(URL, data)
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) Put(URLS []string, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)
		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

func (a *async) PutJSON(URLS []string, data interface{}) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)
		for _, URL := range URLS {
			if a.ctx != nil {
				select {
				case <-a.ctx.Done():
					return
				default:
					jobs <- a.client.PutJSON(URL, data)
				}
				continue
			}
			jobs <- a.client.PutJSON(URL, data)
		}
	}()

	return &Requests{jobs: jobs}
}

func (a *async) PostFile(URLS []string, fieldName, filePath string, uploadForm map[string]string) *Requests {
	jobs := make(chan *Request)

	go func() {
		defer close(jobs)
		for _, URL := range URLS {
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
	}()

	return &Requests{jobs: jobs}
}

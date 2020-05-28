package hhttp

import (
	"compress/zlib"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	request *http.Request
	client  *client
	error   error
}

func (req *Request) Do() (*Response, error) {
	if req.error != nil {
		return nil, req.error
	}

	if err := req.acceptOptions(); err != nil {
		return nil, err
	}

	start := time.Now()

	resp, err := req.client.cli.Do(req.request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "deflate":
		reader, err = zlib.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &Response{
		Body:          body,
		Session:       req.client,
		ContentLength: resp.ContentLength,
		Cookies:       resp.Cookies(),
		Headers:       headers(resp.Header),
		History:       req.client.history,
		StatusCode:    resp.StatusCode,
		Time:          elapsed,
		URL:           resp.Request.URL,
		UserAgent:     req.request.UserAgent(),
		request:       req.request,
		response:      resp,
	}, nil
}

func (req *Request) SetHeaders(headers map[string]string) *Request {
	if headers == nil {
		return req
	}
	for header, data := range headers {
		req.request.Header.Set(header, data)
	}
	return req
}

func (req *Request) AddHeaders(headers map[string]string) *Request {
	if headers == nil {
		return req
	}
	for header, data := range headers {
		req.request.Header.Add(header, data)
	}
	return req
}

func (req *Request) acceptOptions() error {
	userAgent := defaultUserAgent
	req.client.transport.Proxy = http.ProxyFromEnvironment

	if req.client.opt != nil {

		if req.client.opt.BasicAuth != nil && req.request.Header.Get("Authorization") == "" {
			err := req.basicAuth()
			if err != nil {
				return err
			}
		}

		if req.client.opt.UserAgent != nil {
			switch req.client.opt.UserAgent.(type) {
			case string:
				userAgent = req.client.opt.UserAgent.(string)
			case []string:
				userAgent = req.client.opt.UserAgent.([]string)[rand.Intn(len(req.client.opt.UserAgent.([]string)))]
			}
		}

		if req.client.opt.Proxy != nil {
			var proxy string

			switch req.client.opt.Proxy.(type) {
			case string:
				proxy = req.client.opt.Proxy.(string)
			case []string:
				proxy = req.client.opt.Proxy.([]string)[rand.Intn(len(req.client.opt.Proxy.([]string)))]
			}

			if proxyURL, err := url.Parse(proxy); err == nil && proxyURL.Scheme != "" {
				req.client.transport.Proxy = http.ProxyURL(proxyURL)
			}
		}
	}

	req.request.Header.Set("User-Agent", userAgent)
	req.request.Header.Add("Connection", "keep-alive")

	return nil
}

func (req *Request) basicAuth() error {
	baError := errors.New("basic authorization option parameter error")
	user, password := "", ""

	switch req.client.opt.BasicAuth.(type) {
	case string:
		ba := strings.Split(req.client.opt.BasicAuth.(string), ":")
		if len(ba) != 2 {
			return baError
		}
		user = ba[0]
		password = ba[1]
	case []string:
		ba := req.client.opt.BasicAuth.([]string)
		if len(ba) != 2 {
			return baError
		}
		user = ba[0]
		password = ba[1]
	case map[string]string:
		ba := req.client.opt.BasicAuth.(map[string]string)
		if len(ba) != 1 {
			return baError
		}
		for u, p := range ba {
			user = u
			password = p
			break
		}
	default:
		return baError
	}

	req.request.SetBasicAuth(user, password)
	return nil
}

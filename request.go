package hhttp

import (
	"bufio"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	request *http.Request
	client  *Client
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

	elapsed := time.Since(start)

	var streamReader *bufio.Reader

	if req.client.opt != nil && req.client.opt.stream {
		streamReader = bufio.NewReader(resp.Body)
		resp.Body = nil
	}

	deflate := resp.Header.Get("Content-Encoding") == "deflate"

	return &Response{
		Client:        req.client,
		Body:          &body{headers(resp.Header), resp.Body, streamReader, deflate, -1},
		ContentLength: resp.ContentLength,
		Cookies:       resp.Cookies(),
		Headers:       headers(resp.Header),
		History:       req.client.history,
		Proto:         resp.Proto,
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Time:          elapsed,
		URL:           resp.Request.URL,
		UserAgent:     req.request.UserAgent(),
		request:       req.request,
		response:      resp,
	}, nil
}

func (req *Request) AddCookie(cookie *http.Cookie) *Request {
	req.request.AddCookie(cookie)
	return req
}

func (req *Request) SetHeaders(headers map[string]string) *Request {
	if headers != nil && req.request != nil {
		for header, data := range headers {
			req.request.Header.Set(header, data)
		}
	}
	return req
}

func (req *Request) AddHeaders(headers map[string]string) *Request {
	if headers != nil && req.request != nil {
		for header, data := range headers {
			req.request.Header.Add(header, data)
		}
	}
	return req
}

func (req *Request) acceptOptions() error {
	userAgent := defaultUserAgent
	req.client.transport.Proxy = http.ProxyFromEnvironment

	if req.client.opt != nil {
		if req.client.opt.basicAuth != nil && req.request.Header.Get("Authorization") == "" {
			err := req.basicAuth()
			if err != nil {
				return err
			}
		}

		if req.client.opt.userAgent != nil {
			switch req.client.opt.userAgent.(type) {
			case string:
				userAgent = req.client.opt.userAgent.(string)
			case []string:
				userAgent = req.client.opt.userAgent.([]string)[rand.Intn(len(req.client.opt.userAgent.([]string)))]
			}
		}

		if req.client.opt.proxy != nil {
			var proxy string

			switch req.client.opt.proxy.(type) {
			case string:
				proxy = req.client.opt.proxy.(string)
			case []string:
				proxy = req.client.opt.proxy.([]string)[rand.Intn(len(req.client.opt.proxy.([]string)))]
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

	switch req.client.opt.basicAuth.(type) {
	case string:
		ba := strings.Split(req.client.opt.basicAuth.(string), ":")
		if len(ba) != 2 {
			return baError
		}
		user = ba[0]
		password = ba[1]
	case []string:
		ba := req.client.opt.basicAuth.([]string)
		if len(ba) != 2 {
			return baError
		}
		user = ba[0]
		password = ba[1]
	case map[string]string:
		ba := req.client.opt.basicAuth.(map[string]string)
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

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

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &Response{
		Client:        req.client,
		Body:          body{bodyBytes, headers(resp.Header)},
		ContentLength: resp.ContentLength,
		Cookies:       resp.Cookies(),
		Headers:       headers(resp.Header),
		History:       req.client.history,
		Proto:         resp.Proto,
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

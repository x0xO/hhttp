package hhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	Async     *async
	cli       *http.Client
	dialer    *net.Dialer
	history   history
	opt       *Options
	transport *http.Transport
}

func NewClient() *Client {
	c := Client{Async: &async{}}
	c.Async.client = &c

	c.dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	c.transport = &http.Transport{
		DialContext:           c.dialer.DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100, // http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	c.cli = &http.Client{
		Transport: c.transport,
		Timeout:   time.Second * 180,
	}

	return &c
}

func (c *Client) SetOptions(opt *Options) *Client {
	c.opt = opt

	if c.opt.Session {
		c.cli.Jar, _ = cookiejar.New(nil)
	}

	maxRedirects := defaultRedirects
	if c.opt.MaxRedirect != 0 {
		maxRedirects = c.opt.MaxRedirect
	}

	if c.opt.DNS != "" && c.opt.DNSoverTLS == nil {
		c.dialer.Resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				var dialer net.Dialer
				return dialer.DialContext(ctx, "udp", c.opt.DNS)
			},
		}
	}

	if c.opt.DNSoverTLS != nil && c.opt.DNSoverTLS.dnsResolver != nil {
		c.dialer.Resolver = c.opt.DNSoverTLS.dnsResolver
	}

	if c.opt.InterfaceAddr != "" {
		if ip, err := net.ResolveTCPAddr("tcp", c.opt.InterfaceAddr+":0"); err == nil {
			c.dialer.LocalAddr = ip
		}
	}

	if c.opt.Timeout != 0 {
		c.cli.Timeout = time.Second * c.opt.Timeout
	}

	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRedirects)
		}
		if c.opt.History {
			c.history = append(c.history, req.Response)
		}
		return nil
	}

	c.cli.CheckRedirect = redirectPolicy
	return c
}

func (c *Client) Get(URL string, data ...interface{}) *Request {
	if len(data) != 0 {
		return c.buildRequest(URL, http.MethodGet, data[0])
	}
	return c.buildRequest(URL, http.MethodGet, nil)
}

func (c *Client) Delete(URL string, data ...interface{}) *Request {
	if len(data) != 0 {
		return c.buildRequest(URL, http.MethodDelete, data[0])
	}
	return c.buildRequest(URL, http.MethodDelete, nil)
}

func (c *Client) Head(URL string) *Request {
	return c.buildRequest(URL, http.MethodHead, nil)
}

func (c *Client) Post(URL string, data interface{}) *Request {
	return c.buildRequest(URL, http.MethodPost, data)
}

func (c *Client) PostJSON(URL string, data interface{}) *Request {
	return c.buildRequest(URL, http.MethodPost, data)
}

func (c *Client) Put(URL string, data interface{}) *Request {
	return c.buildRequest(URL, http.MethodPut, data)
}

func (c *Client) PutJSON(URL string, data interface{}) *Request {
	return c.buildRequest(URL, http.MethodPut, data)
}

func (c *Client) PostFile(URL, fieldName, filePath string, data ...interface{}) *Request {
	URL = c.urlFormater(URL)

	var (
		reader          io.Reader
		multipartValues map[string]string
	)

	if len(data) > 2 {
		data = data[:2]
	}

	for _, v := range data {
		switch v.(type) {
		case map[string]string:
			multipartValues = v.(map[string]string)
		case string:
			reader = strings.NewReader(v.(string))
		}
	}

	if reader == nil {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return &Request{error: err}
		}
		reader = bytes.NewReader(file)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return &Request{error: err}
	}

	io.Copy(part, reader)

	if multipartValues != nil {
		for field, value := range multipartValues {
			writer.WriteField(field, value)
		}
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, URL, body)
	if err != nil {
		return &Request{error: err}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return &Request{request: req, client: c}
}

func (c Client) getCookies(URL string) []*http.Cookie {
	if c.cli.Jar == nil {
		return nil
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil
	}

	return c.cli.Jar.Cookies(parsedURL)
}

func (c *Client) setCookies(URL string, cookies []*http.Cookie) error {
	if c.cli.Jar == nil {
		return errors.New("cookie jar is not available")
	}

	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	c.cli.Jar.SetCookies(u, cookies)

	return nil
}

func (c *Client) buildRequest(URL, methodType string, data interface{}) *Request {
	URL = c.urlFormater(URL)

	body, contentType, err := c.buildBody(data)
	if err != nil {
		return &Request{error: err}
	}

	req, err := http.NewRequest(methodType, URL, body)
	if err != nil {
		return &Request{error: err}
	}

	if data != nil && contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	return &Request{request: req, client: c}
}

func (c *Client) buildBody(data interface{}) (io.Reader, string, error) {
	var reader io.Reader
	var contentType string

	if data == nil {
		return reader, contentType, nil
	}

	switch data.(type) {
	case []byte:
		contentType = http.DetectContentType(data.([]byte))
		reader = bytes.NewReader(data.([]byte))
	case string:
		var in interface{}
		if json.Unmarshal([]byte(data.(string)), &in) == nil { // if json
			contentType = "application/json; charset=utf-8"
		} else if xml.Unmarshal([]byte(data.(string)), &in) == nil { // if xml
			contentType = "application/xml; charset=utf-8"
		} else {
			contentType = http.DetectContentType([]byte(data.(string)))
		}
		if contentType == "text/plain; charset=utf-8" && strings.ContainsAny(data.(string), "=&") {
			contentType = "application/x-www-form-urlencoded"
		}
		reader = strings.NewReader(data.(string))
	case map[string]string:
		contentType = "application/x-www-form-urlencoded"
		reader = strings.NewReader("")
		form := url.Values{}
		for key, value := range data.(map[string]string) {
			form.Add(key, value)
		}
		reader = strings.NewReader(form.Encode())
	default:
		// TODO: check other types
		switch c.detectDataType(data) {
		case "json":
			contentType = "application/json; charset=utf-8"
			buf, err := json.Marshal(data)
			if err != nil {
				return reader, contentType, err
			}
			reader = bytes.NewBuffer(buf)
		case "xml":
			contentType = "application/xml; charset=utf-8"
			buf, err := xml.Marshal(data)
			if err != nil {
				return reader, contentType, err
			}
			reader = bytes.NewBuffer(buf)
		default:
			return reader, contentType, errors.New("data type not detected")
		}
	}

	return reader, contentType, nil
}

func (c *Client) urlFormater(URL string) string {
	URL = strings.Trim(URL, ".")
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = "http://" + URL
	}

	return URL
}

func (c Client) detectDataType(data interface{}) string {
	value := reflect.ValueOf(data)
	for i := 0; i < value.Type().NumField(); i++ {
		if _, ok := value.Type().Field(i).Tag.Lookup("json"); ok {
			return "json"
		}
		if _, ok := value.Type().Field(i).Tag.Lookup("xml"); ok {
			return "xml"
		}
	}

	return ""
}

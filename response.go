package hhttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	Headers headers
	request *http.Request
	URL     *url.URL
	*Client
	response      *http.Response
	UserAgent     string
	Proto         string
	Status        string
	Body          body
	History       history
	Cookies       cookies
	StatusCode    int
	Time          time.Duration
	ContentLength int64
}

func (resp Response) Referer() string {
	return resp.response.Request.Referer()
}

func (resp Response) GetCookies(URL string) []*http.Cookie {
	return resp.getCookies(URL)
}

func (resp *Response) SetCookie(URL string, cookies []*http.Cookie) error {
	return resp.setCookies(URL, cookies)
}

func (resp Response) Dump(filename string) error {
	if path, err := filepath.Abs(filepath.Dir(filename)); err == nil {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, 0o755)
		}
	}

	return os.WriteFile(filename, resp.Body.bytes, 0o644)
}

func (resp Response) XML(data interface{}) error {
	return xml.Unmarshal(resp.Body.bytes, data)
}

func (resp Response) JSON(data interface{}) error {
	return json.Unmarshal(resp.Body.bytes, data)
}

func (resp Response) Debug(verbos ...bool) {
	body, err := httputil.DumpRequestOut(resp.request, false)
	if err != nil {
		return
	}

	fmt.Println("========= Request ==========")
	fmt.Println(strings.TrimSpace(string(body)))

	cookies := resp.getCookies(resp.request.URL.String())
	if len(cookies) != 0 {
		fmt.Println("\nCookies:")
		for _, cookie := range cookies {
			fmt.Println(cookie)
		}
	}

	fmt.Println("========= Response =========")
	body, err = httputil.DumpResponse(resp.response, false)
	if err != nil {
		return
	}

	fmt.Println(strings.TrimSpace(string(body)))
	fmt.Println("============================")

	if len(verbos) != 0 && verbos[0] {
		fmt.Println(resp.Body)
	}
}

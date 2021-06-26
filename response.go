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
	Body          *body
	History       history
	Cookies       cookies
	StatusCode    int
	Time          time.Duration
	ContentLength int64
}

func (resp Response) Referer() string { return resp.response.Request.Referer() }

func (resp Response) GetCookies(URL string) []*http.Cookie { return resp.getCookies(URL) }

func (resp Response) XML(data interface{}) error { return xml.Unmarshal(resp.Body.Bytes(), data) }

func (resp Response) JSON(data interface{}) error { return json.Unmarshal(resp.Body.Bytes(), data) }

func (resp *Response) SetCookies(URL string, cookies []*http.Cookie) error {
	return resp.setCookies(URL, cookies)
}

func (resp Response) Dump(filename string) error {
	if path, err := filepath.Abs(filepath.Dir(filename)); err == nil {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, 0o755)
		}
	}

	return os.WriteFile(filename, resp.Body.Bytes(), 0o644)
}

type printer string

func (resp Response) Debug(verbos ...bool) printer {

	var builder strings.Builder

	body, err := httputil.DumpRequestOut(resp.request, false)
	if err != nil {
		return printer(builder.String() + "\n")
	}

	builder.WriteString("========= Request ==========\n")
	builder.WriteString(strings.TrimSpace(string(body)) + "\n")

	cookies := resp.getCookies(resp.request.URL.String())
	if len(cookies) != 0 {
		builder.WriteString("\nCookies:\n")
		for _, cookie := range cookies {
			builder.WriteString(fmt.Sprint(cookie) + "\n")
		}
	}

	builder.WriteString("========= Response =========\n")
	body, err = httputil.DumpResponse(resp.response, false)
	if err != nil {
		return printer(builder.String() + "\n")
	}

	builder.WriteString(strings.TrimSpace(string(body)))
	builder.WriteString("\n============================\n")

	if len(verbos) != 0 && verbos[0] {
		builder.WriteString("=========== Body ===========\n")
		builder.WriteString(resp.Body.String())
	}

	return printer(builder.String() + "\n")
}

func (p printer) Print() {
	fmt.Println(p)
}

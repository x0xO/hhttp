package hhttp

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestMultipart(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		var buff bytes.Buffer
		if r.FormValue("some") == "values" {
			buff.WriteString(r.FormValue("some"))
		}
		w.Write(buff.Bytes())
	}))
	defer ts.Close()

	multipartData := map[string]string{"some": "values"}

	r, err := NewClient().Multipart(ts.URL, multipartData).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if r.Body.String() != "values" {
		t.Error()
	}
}

func TestFileUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		var buff bytes.Buffer
		if r.FormValue("some") == "values" {
			buff.WriteString(r.FormValue("some"))
		}

		file, _, _ := r.FormFile("file")
		defer file.Close()

		io.Copy(&buff, file)
		w.Write(buff.Bytes())
	}))
	defer ts.Close()

	r, err := NewClient().FileUpload(ts.URL, "file", "info.txt", "justfile").Do()
	if err != nil {
		t.Error(err)
		return
	}

	multipartValues := map[string]string{"some": "values"}
	r2, err := NewClient().FileUpload(ts.URL, "file", "info.txt", "multipart", multipartValues).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if r.Body.String() != "justfile" || r2.Body.String() != "valuesmultipart" {
		t.Error()
	}
}

func TestDeflate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		w2 := zlib.NewWriter(buf)
		w2.Write([]byte("OK"))
		w2.Close()

		w.Header().Set("Content-Encoding", "deflate")
		w.Write(buf.Bytes())
	}))
	defer ts.Close()

	r, err := NewClient().Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") || !r.Body.Contains([]byte("OK")) {
		t.Error()
	}
}

func TestBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))
	defer ts.Close()

	r, err := NewClient().Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") || !r.Body.Contains([]byte("OK")) {
		t.Error()
	}
}

func TestTimeOut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Nanosecond)
		io.WriteString(w, "OK")
	}))
	defer ts.Close()

	_, err := NewClient().SetOptions(NewOptions().Timeout(time.Microsecond)).Get(ts.URL).Do()
	r, _ := NewClient().SetOptions(NewOptions().Timeout(time.Second)).Get(ts.URL).Do()

	if err == nil || !r.Body.Contains("OK") {
		t.Error()
	}
}

func TestSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cookie http.Cookie

		cookies, err := r.Cookie("username")
		if err == http.ErrNoCookie {
			cookie = http.Cookie{Name: "username", Value: "root"}
		} else if cookies.Value == "root" {
			cookie = http.Cookie{Name: "username", Value: "toor"}
		}

		http.SetCookie(w, &cookie)
	}))
	defer ts.Close()

	r, err := NewClient().SetOptions(NewOptions().Session()).Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	r.Body.Close()

	r, err = r.Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	cookies := r.GetCookies(ts.URL)

	if !reflect.DeepEqual(cookies, []*http.Cookie{{Name: "username", Value: "toor"}}) {
		t.Error()
	}
}

func TestBasicAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		username, password, ok := r.BasicAuth()

		if !ok {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		if username != "good" || password != "password" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
	}))

	defer ts.Close()

	r, err := NewClient().SetOptions(NewOptions().BasicAuth("good:password")).Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Body.Close()

	r2, err := NewClient().SetOptions(NewOptions().BasicAuth("bad:password")).Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}
	defer r2.Body.Close()

	if r.StatusCode != http.StatusOK || r2.StatusCode != http.StatusUnauthorized {
		t.Error()
	}
}

func TestCookies(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("root"); err == nil {
			if cookie.Value == "cookie" {
				io.WriteString(w, "OK")
			}
		}
	}))
	defer ts.Close()

	c1 := http.Cookie{Name: "root", Value: "cookie"}
	r, err := NewClient().Get(ts.URL).AddCookies(c1).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") {
		t.Error()
	}
}

func TestUserAgent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.UserAgent())
	}))
	defer ts.Close()

	agent := "Hi from hhttp"
	r, err := NewClient().SetOptions(NewOptions().UserAgent(agent)).Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains(agent) {
		t.Error()
	}
}

func TestHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		some := r.Header.Get("Some")
		if some == "header" {
			io.WriteString(w, "OK")
		}
	}))
	defer ts.Close()

	headers := map[string]string{"some": "header"}
	r, err := NewClient().Get(ts.URL).AddHeaders(headers).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") {
		t.Error()
	}
}

func TestHTTP2(t *testing.T) {
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.Proto)
		}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()

	r, err := NewClient().SetOptions(NewOptions().HTTP2()).Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("Hello, HTTP/2.0") {
		t.Error()
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))
	defer ts.Close()

	r, err := NewClient().Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") {
		t.Error()
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("test") == "data" {
			io.WriteString(w, "OK")
		}
	}))
	defer ts.Close()

	r, err := NewClient().Post(ts.URL, "test=data").Do()
	if err != nil {
		t.Error(err)
		return
	}

	if !r.Body.Contains("OK") {
		t.Error()
	}
}

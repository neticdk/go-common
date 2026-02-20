package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/neticdk/go-stdlib/file"
)

type Request struct {
	req    *http.Request
	client *http.Client
	resp   *http.Response
}

func NewRequest(options ...Option) *Request {
	r := &Request{
		client: &http.Client{
			Transport: http.DefaultTransport,
		},
	}

	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		r.client.Transport = t.Clone()
	}

	r.req, _ = http.NewRequest("", "", nil)
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *Request) SetHeader(header map[string]string) {
	for k, v := range header {
		r.req.Header.Set(k, v)
	}
}

func (r *Request) SetBasicAuth(username, password string) {
	r.req.SetBasicAuth(username, password)
}

func (r *Request) SetBearerTokenAuth(token string) {
	r.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (r *Request) SetTimeout(t time.Duration) {
	r.client.Timeout = t
}

func (r *Request) SetTransport(rt http.RoundTripper) {
	r.client.Transport = rt
}

func (r *Request) SetClient(client *http.Client) {
	r.client = client
}

func (r *Request) SetSkipTLS() {
	if r.client.Transport == nil {
		r.client.Transport = &http.Transport{}
	}

	transport, ok := r.client.Transport.(*http.Transport)
	if !ok {
		return
	}

	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	transport.TLSClientConfig.InsecureSkipVerify = true // #nosec
}

func (r *Request) SetContext(ctx context.Context) {
	r.req = r.req.WithContext(ctx)
}

func (r *Request) parseURL(originURL string) error {
	var err error
	r.req.URL, err = url.Parse(originURL)
	return err
}

func (r *Request) DownloadToWriter(originURL string, w io.Writer) (int64, error) {
	r.req.Method = http.MethodGet
	err := r.parseURL(originURL)
	if err != nil {
		return 0, err
	}
	// #nosec G704 -- URL scheme is validated, and risk is acknowledged
	resp, err := r.client.Do(r.req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	r.resp = resp
	return io.Copy(w, r.resp.Body)
}

func (r *Request) Download(originURL, filePath string) (int64, error) {
	f, err := file.SafeCreate(filepath.Dir(filePath), filePath, file.FileModeNewFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return r.DownloadToWriter(originURL, f)
}

func (r Request) Status() (int, string) {
	return r.resp.StatusCode, r.resp.Status
}

func (r Request) Response() *http.Response {
	return r.resp
}

func (r Request) Request() *http.Request {
	return r.req
}

func Download(originURL, filePath string) (int64, error) {
	return NewRequest().Download(originURL, filePath)
}

func DownloadCtx(ctx context.Context, originURL, filePath string) (int64, error) {
	return NewRequest(WithContext(ctx)).Download(originURL, filePath)
}

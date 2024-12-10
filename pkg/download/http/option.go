package http

import (
	"context"
	"net/http"
	"time"
)

type Option func(*Request)

func WithTimeout(t time.Duration) Option {
	return func(r *Request) {
		r.SetTimeout(t)
	}
}

func WithHeader(headers ...map[string]string) Option {
	return func(r *Request) {
		for _, header := range headers {
			r.SetHeader(header)
		}
	}
}

func WithSkipTLS() Option {
	return func(r *Request) {
		r.SetSkipTLS()
	}
}

func WithBasicAuth(username, password string) Option {
	return func(r *Request) {
		r.SetBasicAuth(username, password)
	}
}

func WithBearerTokenAuth(token string) Option {
	return func(r *Request) {
		r.SetBearerTokenAuth(token)
	}
}

func WithTransport(tr http.RoundTripper) Option {
	return func(r *Request) {
		r.SetTransport(tr)
	}
}

func WithClient(client *http.Client) Option {
	return func(r *Request) {
		r.SetClient(client)
	}
}

func WithContext(ctx context.Context) Option {
	return func(r *Request) {
		r.SetContext(ctx)
	}
}

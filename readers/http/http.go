package http

import (
	"bytes"
	"errors"
	"io"
	nethttp "net/http"
)

type Option func(*nethttp.Request)

func New(url string, opts ...Option) (*http, error) {
	client := nethttp.Client{}

	request, err := nethttp.NewRequest(nethttp.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(request)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.New("http error")
	}

	rawBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &http{reader: bytes.NewReader(rawBytes)}, nil
}

type http struct {
	reader io.Reader
}

func (f *http) GetReader() io.Reader {
	return f.reader
}

func WithHeaders(headers map[string][]string) Option {
	return func(r *nethttp.Request) {
		for name, values := range headers {
			for _, value := range values {
				r.Header.Add(name, value)
			}
		}
	}
}

func WithQueryParams(params map[string][]string) Option {
	return func(r *nethttp.Request) {
		query := r.URL.Query()

		for name, values := range params {
			for _, value := range values {
				query.Add(name, value)
			}
		}

		r.URL.RawQuery = query.Encode()
	}
}

func WithBasicAuth(username, password string) Option {
	return func(r *nethttp.Request) {
		r.SetBasicAuth(username, password)
	}
}

package {{.Params.package}}

import (
    "net/http"
    "errors"
)

type option func(*options)

type options struct {
    httpClient http.Client
    baseURL string
}

func WithHTTPClient(c http.Client) option {
    return func(o *options) {
        o.httpClient = c
    }
}

func WithBaseURL(url string) option {
    return func(o *options) {
        o.baseURL = url
    }
}

func New(opts ...option) *Client {
    options := options{}

    for _, opt := range opts {
        opt(&options)
    }

    return &Client{
        options: options,
    }
}

type Client struct {
    options options 
}

{{ range $k, $v := .Operations }}
func (c *Client) {{$k}}({{if $v.Input }}in {{ ref $v.Input }}{{ end }}) (*{{ ref $v.Output }}, error) {
    return nil, errors.New("")
}
{{ end }}
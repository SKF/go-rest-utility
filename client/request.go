package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-http-utils/headers"
	"github.com/jtacoma/uritemplates"
)

type Request struct {
	uriTemplate  string
	uriVariables map[string]interface{}

	method string
	header http.Header
	body   io.Reader
}

func NewRequest(method, uriTemplate string) *Request {
	return &Request{
		uriTemplate:  uriTemplate,
		uriVariables: make(map[string]interface{}),

		method: method,
		header: make(http.Header),
		body:   http.NoBody,
	}
}

func Get(uriTemplate string) *Request {
	return NewRequest(http.MethodGet, uriTemplate)
}

func Head(uriTemplate string) *Request {
	return NewRequest(http.MethodHead, uriTemplate)
}

func Post(uriTemplate string) *Request {
	return NewRequest(http.MethodPost, uriTemplate)
}

func Put(uriTemplate string) *Request {
	return NewRequest(http.MethodPut, uriTemplate)
}

func Delete(uriTemplate string) *Request {
	return NewRequest(http.MethodDelete, uriTemplate)
}

func Connect(uriTemplate string) *Request {
	return NewRequest(http.MethodConnect, uriTemplate)
}

func Options(uriTemplate string) *Request {
	return NewRequest(http.MethodOptions, uriTemplate)
}

func Trace(uriTemplate string) *Request {
	return NewRequest(http.MethodTrace, uriTemplate)
}

func Patch(uriTemplate string) *Request {
	return NewRequest(http.MethodPatch, uriTemplate)
}

func (r *Request) Assign(variable string, value interface{}) *Request {
	r.uriVariables[variable] = value

	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.header.Set(key, value)

	return r
}

type jsonPayload struct {
	payload interface{}
	buffer  io.Reader
}

func (jp *jsonPayload) Read(p []byte) (n int, err error) {
	if jp.buffer == nil {
		switch payload := jp.payload.(type) {
		case []byte:
			jp.buffer = bytes.NewBuffer(payload)
		case string:
			jp.buffer = bytes.NewBufferString(payload)
		default:
			buf := new(bytes.Buffer)

			if err = json.NewEncoder(buf).Encode(payload); err != nil {
				return
			}

			jp.buffer = buf
		}
	}

	return jp.Read(p)
}

func (r *Request) WithJSONPayload(payload interface{}) *Request {
	r.header.Set(headers.ContentType, "application/json")
	r.body = &jsonPayload{payload: payload}

	return r
}

func (r *Request) toURL(baseURL *url.URL) (*url.URL, error) {
	template, err := uritemplates.Parse(r.uriTemplate)
	if err != nil {
		return nil, fmt.Errorf("unable to parse uri template: %w", err)
	}

	expandedTemplate, err := template.Expand(r.uriVariables)
	if err != nil {
		return nil, fmt.Errorf("unable to expand uri template: %w", err)
	}

	templateURL, err := url.Parse(expandedTemplate)
	if err != nil {
		return nil, fmt.Errorf("unable to parse expanded uri template: %w", err)
	}

	if baseURL == nil {
		return templateURL, nil
	}

	return baseURL.ResolveReference(templateURL), nil
}

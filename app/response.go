package main

import "fmt"

type Response struct {
	HttpVersion string
	StatusCode  int
	StatusStr   string
	Header      map[string]string
	Body        string
}

func NewResponse() Response {
	return Response{
		HttpVersion: "HTTP/1.1",
		StatusCode:  0,
		StatusStr:   "",
		Header:      map[string]string{},
		Body:        "",
	}
}

func (r *Response) AddHeader(key, val string) {
	r.Header[key] = val
}

func (r *Response) Compose() string {
	res := fmt.Sprintf("%s %d %s\r\n", r.HttpVersion, r.StatusCode, r.StatusStr)
	for key, val := range r.Header {
		res += fmt.Sprintf("%s: %s\r\n", key, val)
	}
	res += "\r\n"
	res += r.Body
	return res
}

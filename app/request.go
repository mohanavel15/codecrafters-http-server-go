package main

import "strings"

type Request struct {
	Method      string
	Path        string
	HttpVersion string
	Headers     map[string]string
}

func ParseRequest(request_str string) Request {
	request := Request{}

	lines := strings.Split(request_str, "\r\n")
	status_line := strings.Split(lines[0], " ")

	request.Method = status_line[0]
	request.Path = status_line[1]
	request.HttpVersion = status_line[2]

	request.Headers = map[string]string{}

	idx := 1
	for lines[idx] != "" {
		header := strings.Split(lines[idx], ":")
		key := strings.TrimSpace(header[0])
		val := strings.TrimSpace(header[1])
		request.Headers[key] = val
		idx += 1
	}

	return request
}

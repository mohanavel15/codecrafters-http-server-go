package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	Handler(conn)
}

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

func Handler(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf = buf[:n]

	request_str := string(buf)
	lines := strings.Split(request_str, "\r\n")
	status_line := strings.Split(lines[0], " ")

	path := status_line[1]

	res := NewResponse()

	if path == "/" {
		res.StatusCode = 200
		res.StatusStr = "OK"
	} else if strings.Contains(path, "/echo/") {
		res.StatusCode = 200
		res.StatusStr = "OK"

		word := strings.Split(path, "/")[2]

		res.AddHeader("Content-Type", "text/plain")
		res.AddHeader("Content-Length", fmt.Sprint(len(word)))
		res.Body = word

	} else {
		res.StatusCode = 404
		res.StatusStr = "Not Found"
	}

	conn.Write([]byte(res.Compose()))
}

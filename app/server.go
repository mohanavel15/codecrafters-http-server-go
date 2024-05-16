package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	buf = buf[:n]

	request := ParseRequest(string(buf))

	res := NewResponse()

	if request.Path == "/" {
		res.StatusCode = 200
		res.StatusStr = "OK"
	} else if strings.Contains(request.Path, "/echo/") {
		res.StatusCode = 200
		res.StatusStr = "OK"

		shouldEncode := false
		if encode, ok := request.Headers["Accept-Encoding"]; ok {
			if strings.Contains(encode, "gzip") {
				shouldEncode = true
				res.AddHeader("Content-Encoding", "gzip")
			}
		}

		word := strings.Split(request.Path, "/")[2]

		res.AddHeader("Content-Type", "text/plain")

		if shouldEncode {
			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			w.Write([]byte("hello, world\n"))
			w.Close()

			word = b.String()
		}

		res.AddHeader("Content-Length", fmt.Sprint(len(word)))
		res.Body = word

	} else if strings.Contains(request.Path, "/files/") {
		if len(os.Args) < 3 {
			res.StatusCode = 404
			res.StatusStr = "Not Found"
		} else {
			dir := os.Args[2]
			filename := strings.Split(request.Path, "/")[2]
			filepath := fmt.Sprintf("%s/%s", dir, filename)

			if request.Method == "GET" {
				raw, err := os.ReadFile(filepath)
				if err != nil {
					res.StatusCode = 404
					res.StatusStr = "Not Found"
				} else {
					res.StatusCode = 200
					res.StatusStr = "OK"

					res.AddHeader("Content-Type", "application/octet-stream")
					res.AddHeader("Content-Length", fmt.Sprint(len(raw)))
					res.Body = string(raw)
				}
			} else if request.Method == "POST" {
				err := os.WriteFile(filepath, []byte(request.Body), 0644)
				if err != nil {
					res.StatusCode = 404
					res.StatusStr = "Not Found"
				} else {
					res.StatusCode = 201
					res.StatusStr = "Created"
				}
			}
		}

	} else if request.Path == "/user-agent" {
		res.StatusCode = 200
		res.StatusStr = "OK"

		userAgent := request.Headers["User-Agent"]

		res.AddHeader("Content-Type", "text/plain")
		res.AddHeader("Content-Length", fmt.Sprint(len(userAgent)))
		res.Body = userAgent
	} else {
		res.StatusCode = 404
		res.StatusStr = "Not Found"
	}

	conn.Write([]byte(res.Compose()))
}

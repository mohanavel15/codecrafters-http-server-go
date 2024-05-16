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

	dir := "./"
	if len(os.Args) >= 3 && os.Args[1] == "--directory" {
		dir = os.Args[2]
	}

	r := NewRouter()
	r.AddRoute("GET", "/", func(req Request, res *Response) {
		// Nothing
	})

	r.AddRoute("GET", "/user-agent", func(req Request, res *Response) {
		userAgent := req.Headers["User-Agent"]

		res.AddHeader("Content-Type", "text/plain")
		res.AddHeader("Content-Length", fmt.Sprint(len(userAgent)))
		res.Body = userAgent
	})

	r.AddRoute("GET", "/echo/{any}", func(req Request, res *Response) {
		word := strings.Split(req.Path, "/")[2]
		res.Body = word
	})

	r.AddRoute("GET", "/files/{any}", func(req Request, res *Response) {
		filename := strings.Split(req.Path, "/")[2]
		filepath := fmt.Sprintf("%s/%s", dir, filename)

		raw, err := os.ReadFile(filepath)
		if err != nil {
			res.StatusCode = 404
			res.StatusStr = "Not Found"
		} else {
			res.AddHeader("Content-Type", "application/octet-stream")
			res.Body = string(raw)
		}
	})

	r.AddRoute("POST", "/files/{any}", func(req Request, res *Response) {
		filename := strings.Split(req.Path, "/")[2]
		filepath := fmt.Sprintf("%s/%s", dir, filename)

		res.StatusCode = 201
		res.StatusStr = "Created"

		err := os.WriteFile(filepath, []byte(req.Body), 0644)
		if err != nil {
			res.StatusCode = 500
			res.StatusStr = "Internal Server Error"
		}
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go Handler(conn, &r)
	}
}

func Handler(conn net.Conn, router *Router) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	buf = buf[:n]

	req := ParseRequest(string(buf))

	res := NewResponse()
	res.StatusCode = 404
	res.StatusStr = "Not Found"
	res.AddHeader("Content-Type", "text/plain")

	router.Handle(req, &res)

	if encode, ok := req.Headers["Accept-Encoding"]; ok {
		if strings.Contains(encode, "gzip") {
			res.AddHeader("Content-Encoding", "gzip")
			var buf bytes.Buffer

			w := gzip.NewWriter(&buf)
			w.Write([]byte(res.Body))
			w.Close()

			res.Body = buf.String()
		}
	}

	res.AddHeader("Content-Length", fmt.Sprint(len(res.Body)))

	conn.Write([]byte(res.Compose()))
}

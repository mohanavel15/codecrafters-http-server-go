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

		word := strings.Split(request.Path, "/")[2]

		res.AddHeader("Content-Type", "text/plain")
		res.AddHeader("Content-Length", fmt.Sprint(len(word)))
		res.Body = word

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

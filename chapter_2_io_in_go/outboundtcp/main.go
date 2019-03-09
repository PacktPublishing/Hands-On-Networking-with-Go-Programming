package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		fmt.Println("error connecting", err)
		os.Exit(1)
	}
	conn.Write([]byte("GET / HTTP/1.0\n\n"))
	data, err := readAll(conn)
	fmt.Println(string(data), err)
}

func readAll(r io.Reader) (data []byte, err error) {
	var b bytes.Buffer
	buf := make([]byte, 32*1024)
	var read int
	for {
		read, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if read == 0 {
			break
		}
		b.Write(buf[:read])
	}
	err = nil
	data = b.Bytes()
	return
}

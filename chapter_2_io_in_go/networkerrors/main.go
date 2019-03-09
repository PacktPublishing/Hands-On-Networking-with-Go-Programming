package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9999", time.Second*5)
	if err != nil {
		fmt.Println("failed to connect")
		return
	}
	err = conn.SetDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		fmt.Println("error setting deadline", err)
		return
	}
	start := time.Now()
	fmt.Println("Connected")
	defer func() {
		fmt.Println("Completed in", time.Now().Sub(start))
	}()
	expectedToRead := 2
	data := make([]byte, 2)
	read, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if read != expectedToRead {
		fmt.Println("expected to read 2 bytes, read", read)
		return
	}
}

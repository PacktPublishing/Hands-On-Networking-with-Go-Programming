package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	id := 0
	for {
		go connect(id)
		time.Sleep(time.Second)
		id++
	}
}

func connect(id int) {
	conn, err := net.Dial("tcp", "localhost:8050")
	if err != nil {
		return
	}
	defer conn.Close()
	data := make([]byte, 8)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	if n != 8 {
		fmt.Println("Unexpected data received")
		return
	}
	if isConnectionRejected(data) {
		fmt.Printf("Connection %d was rejected by server...\n", id)
		return
	}
	for {
		n, err := conn.Read(data)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading:", err)
			return
		}
		if n != 8 {
			fmt.Println("Unexpected data received")
			return
		}
		fmt.Printf("%d: Received data\n", id)
	}
}

func isConnectionRejected(data []byte) bool {
	for _, d := range data {
		if d != 1 {
			return true
		}
	}
	return false
}

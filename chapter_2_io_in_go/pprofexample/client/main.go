package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	for {
		sayHello()
		time.Sleep(time.Second)
	}
}

func sayHello() {
	conn, err := net.Dial("tcp", "localhost:8181")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte("Hello"))
}

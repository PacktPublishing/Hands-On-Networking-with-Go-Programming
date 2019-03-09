package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	var id int
	for i := 0; i < 100000; i++ {
		go connect(&id)
		time.Sleep(time.Millisecond * 1)
	}
}

func connect(id *int) {
	conn, err := net.Dial("tcp", "localhost:8181")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	*id++
	fmt.Println("Connected:", *id)
	for {
		conn.Write([]byte("a"))
		time.Sleep(time.Minute)
	}
}

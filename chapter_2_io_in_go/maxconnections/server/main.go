package main

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 8181})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var connections uint64
	for {
		conn, err := l.AcceptTCP()
		if err != nil && err.Error() != "use of closed network connection" {
			fmt.Println("Error:", err)
			continue
		}
		count := atomic.AddUint64(&connections, 1)
		data := make([]byte, 1)
		go func(c uint64) {
			defer conn.Close()
			fmt.Println("Accepted connection", c)
			conn.Read(data)
			time.Sleep(time.Hour)
		}(count)
	}
}

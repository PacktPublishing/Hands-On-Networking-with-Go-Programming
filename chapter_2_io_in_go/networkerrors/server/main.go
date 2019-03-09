package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("error listening", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("error accepting", err)
				continue
			}
			go func(c net.Conn) {
				fmt.Println("Connection received, waiting 10 seconds")
				time.Sleep(time.Second * 10)
				c.Write([]byte("OK"))
				fmt.Println("Closing connection")
				c.Close()
			}(conn)
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	// conn.SetKeepAlive()
	// conn.SetKeepAlivePeriod()
	// conn.SetLinger()
	// conn.SetNoDelay()
	// conn.SetReadBuffer()
	// conn.SetReadDeadline()
}

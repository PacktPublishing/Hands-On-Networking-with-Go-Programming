package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var ports = []int{
	22,   // SSH
	23,   // Telnet
	80,   // HTTP
	443,  // HTTPS
	8080, // HTTP
	53,   // Domain
	1433, // SQL Server
	3306, // MySQL
}

var targetFlag = flag.String("target", "192.168.0.1", "The target to scan")

func main() {
	flag.Parse()

	target := net.ParseIP(*targetFlag)
	if target == nil {
		fmt.Println("Failed to parse target address, is it formatted correctly?")
		return
	}

	var wg sync.WaitGroup
	for _, port := range ports {
		port := port
		go func() {
			wg.Add(1)
			defer wg.Done()
			addr := &net.TCPAddr{
				IP:   target,
				Port: port,
			}
			if portActive(addr.String()) {
				fmt.Println(addr.String())
			}
		}()
	}
	wg.Wait()
}

func portActive(addr string) bool {
	var d net.Dialer
	// Wait for 500 ms per connection.
	d.Timeout = time.Millisecond * 500
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

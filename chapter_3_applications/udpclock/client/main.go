package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 5),
		Port: 3002,
	})
	if err != nil {
		fmt.Printf("error starting to listen: %v\n", err)
		os.Exit(1)
	}

	go func() {
		prefix := []byte("time:")
		for {
			data := make([]byte, 1024)
			fmt.Println("Reading from socket")
			_, remoteAddr, err := conn.ReadFromUDP(data)
			fmt.Println("Read socket")
			if err != nil {
				fmt.Printf("error reading from UDP: %v\n", err)
				continue
			}
			// Check that the data is something we're interested in.
			if bytes.HasPrefix(data, prefix) {
				// Get the time from the server.
				// Skip the prefix ("time:")
				suffix := data[len(prefix):]
				// Read the int64
				ns := int64(binary.LittleEndian.Uint64(suffix))
				fmt.Println(remoteAddr, time.Unix(0, ns), ns)
			}
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("Stopped.")
}

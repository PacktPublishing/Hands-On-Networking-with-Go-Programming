package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Start the clock ticking.
	clock := make(chan time.Time)
	go func() {
		for {
			clock <- time.Now()
			time.Sleep(time.Second)
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-sigs:
			fmt.Println("shutting down...")
			break loop
		case t := <-clock:
			conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   net.IPv4(224, 0, 0, 5), // OSPF routing protocol address.
				Port: 3002,
			})
			if err != nil {
				fmt.Printf("error creating broadcast socket: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Writing to socket")
			// We need 5 bytes for "time:" and then another 8 bytes for the int64.
			data := []byte("time:")
			ns := make([]byte, 8)
			binary.LittleEndian.PutUint64(ns, uint64(t.UnixNano()))
			_, err = conn.Write(append(data, ns...))
			if err != nil {
				fmt.Printf("error writing header to the socket: %v\n", err)
			}
			fmt.Println("Closing socket")
			conn.Close()
		}
	}
	fmt.Println("Complete.")
}

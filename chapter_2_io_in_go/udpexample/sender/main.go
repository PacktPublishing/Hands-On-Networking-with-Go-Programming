package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	laddr := &net.UDPAddr{
		Port: 3022,
	}
	raddr := &net.UDPAddr{
		Port: 3002,
	}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Printf("error starting to dial: %v\n", err)
		os.Exit(1)
	}
	go func() {
		response := make([]byte, 8)
		for {
			_, err := conn.Read(response)
			if err != nil {
				fmt.Printf("error reading data: %v\n", err)
				continue
			}
			fmt.Println("Received", string(response))
		}
	}()
	bytes := make([]byte, 8)
	for i := int64(0); i < 255; i++ {
		binary.LittleEndian.PutUint64(bytes, uint64(i))
		_, err := conn.Write(bytes)
		if err != nil {
			fmt.Printf("error sending data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(i)
	}
	fmt.Println("Sent data.")
	time.Sleep(time.Second * 5)
}

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	laddr := &net.UDPAddr{
		Port: 3002,
	}
	raddr := &net.UDPAddr{
		Port: 3022,
	}
	l, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Printf("error starting to listen: %v\n", err)
		os.Exit(1)
	}
	data := make([]byte, 8)
	for {
		_, err := l.Read(data)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(hex(data))
		_, err = l.Write([]byte("      OK"))
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}

func hex(v []byte) string {
	s := make([]string, len(v))
	for i, v := range v {
		s[i] = fmt.Sprintf("0x%02x", v)
	}
	return strings.Join(s, " ")
}

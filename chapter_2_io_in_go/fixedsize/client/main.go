package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9092")
	if err != nil {
		fmt.Println("error connecting", err)
		os.Exit(1)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("github.com/a-h/timeseries/v1"))
	if err != nil {
		fmt.Println("error sending preamble", err)
		os.Exit(1)
	}
	ok := make([]byte, 2)
	read, err := conn.Read(ok)
	if err != nil {
		fmt.Println("error reading preamble", err)
		os.Exit(1)
	}
	if read != len(ok) || !isOK(ok) {
		fmt.Println("invalid preamble returned")
		os.Exit(1)
	}
	// Send some samples.
	for x := 0.0; x < 2.0; x += 0.1 {
		d := Data{
			SeriesID:  1,
			Timestamp: time.Now().UnixNano(),
			Value:     math.Sin(10.0 * x),
		}
		err = binary.Write(conn, binary.LittleEndian, d)
		if err != nil {
			fmt.Println("failed to write to server", err)
			os.Exit(1)
		}
	}
}

type Data struct {
	SeriesID  byte
	Timestamp int64
	Value     float64
}

func isOK(ok []byte) bool {
	return ok[0] == 'O' && ok[1] == 'K'
}

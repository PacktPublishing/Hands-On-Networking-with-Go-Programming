package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	list, err := net.Listen("tcp", "127.0.0.1:9092")
	if err != nil {
		fmt.Println("error creating listener", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	go func() {
		for {
			conn, err := list.Accept()
			if err != nil {
				fmt.Println("error accepting connection", err)
				continue
			}
			go handleConnection(ctx, &wg, conn)
		}
	}()

	<-sigs
	cancel()
	wg.Wait()
}

var expectedPreamble = "github.com/a-h/timeseries/v1"

func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	wg.Add(1)
	defer wg.Done()
	defer conn.Close()

	// Handshake with the client.
	preamble := make([]byte, len(expectedPreamble))
	read, err := conn.Read(preamble)
	if err != nil {
		fmt.Println("error reading from connection")
		return
	}
	if read != len(expectedPreamble) {
		fmt.Println("client sent unexpected message")
		return
	}
	_, err = conn.Write([]byte("OK"))
	if err != nil {
		fmt.Println("error during handshake")
		return
	}

	// Start reading normally.
	dataSize := 1 + 8 + 8 // one byte for the series, one for the time (int64) and one for the data (float64)
	data := make([]byte, dataSize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if read, err = conn.Read(data); err != nil || read != dataSize {
				fmt.Println("communication stopped")
				return
			}
			// Read the data
			rdr := bytes.NewReader(data)
			var d Data
			err = binary.Read(rdr, binary.LittleEndian, &d)
			if err != nil {
				fmt.Println("failed to read into Data", err)
				return
			}
			fmt.Printf("%s: %d %v %v\n", conn.RemoteAddr().String(),
				int(d.SeriesID),
				time.Unix(0, d.Timestamp),
				d.Value)
		}
	}
}

type Data struct {
	SeriesID  byte
	Timestamp int64
	Value     float64
}

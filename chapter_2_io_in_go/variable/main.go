package main

import (
	"bytes"
	"fmt"
)

func main() {
	var payload bytes.Buffer
	data := []byte{1, 2, 3, 4, 5}
	err := payload.WriteByte(byte(len(data)))
	if err != nil {
		fmt.Println("failed to write length of data")
		return
	}
	_, err = payload.Write(data)
	if err != nil {
		fmt.Println("failed to write data to buffer")
		return
	}
	// Write the data to the connection
	// _, err = conn.Write(payload)
}

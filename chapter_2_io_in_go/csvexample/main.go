package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

const data = `A1,B1
A2,B2
A3,B3
A4,B4
A5,B5`

func main() {
	r := csv.NewReader(bytes.NewReader([]byte(data)))
	var record []string
	var err error
	for record, err = r.Read(); err != io.EOF; record, err = r.Read() {
		fmt.Println(record[0] + "|" + record[1])
	}
	if err != io.EOF {
		fmt.Println("failed to read CSV data", err)
		os.Exit(1)
	}
}

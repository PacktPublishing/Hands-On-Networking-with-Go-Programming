package main

import (
	"bufio"
	"bytes"
	"fmt"
)

const text = `A|B|C|D|E|F|G`

func main() {
	pipeSplit := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		// Split at the pipe.
		if i := bytes.IndexByte(data, '|'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, return everything up to that.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
	s := bufio.NewScanner(bytes.NewBufferString(text))
	s.Split(pipeSplit)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	if s.Err() != nil {
		fmt.Println("error scanning data", s.Err())
	}
}

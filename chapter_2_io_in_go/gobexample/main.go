package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

type book struct {
	Name  string `gob:"n"`
	Stars int    `gob:"s"`
}

func main() {
	// The io.Writer to write to.
	var w bytes.Buffer
	b := book{
		Name:  "Go Network Programming",
		Stars: 5,
	}
	e := gob.NewEncoder(&w)
	err := e.Encode(b)
	if err != nil {
		fmt.Println("error encoding Gob:", err)
		return
	}
	fmt.Println("Bytes written:", w.Len())
	fmt.Println(string(w.Bytes()))
}

func hex(v []byte) string {
	s := make([]string, len(v))
	for i, v := range v {
		s[i] = fmt.Sprintf("0x%02x", v)
	}
	return strings.Join(s, " ")
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type book struct {
	Name  string `json:"name"`
	Stars int    `json:"stars"`
}

func main() {
	// The io.Writer to write to.
	var w bytes.Buffer
	b := book{
		Name:  "Go Network Programming",
		Stars: 5,
	}
	e := json.NewEncoder(&w)
	err := e.Encode(b)
	if err != nil {
		fmt.Println("error encoding JSON:", err)
		return
	}
	fmt.Println("Bytes written:", w.Len())
}

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
	j := `{
		"name": "Go Network Programming",
		"stars": 5
}`
	buf := bytes.NewBufferString(j)

	var b book
	d := json.NewDecoder(buf)
	err := d.Decode(&b)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
		return
	}
	fmt.Println("Book Name:", b.Name)
	fmt.Println("Book Stars:", b.Stars)
}

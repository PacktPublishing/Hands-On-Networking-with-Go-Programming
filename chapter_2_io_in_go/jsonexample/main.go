package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Book struct {
	Name  string `json:"name"`
	Stars int    `json:"stars"`
}

func main() {
	b := Book{
		Name:  "Book",
		Stars: 2,
	}
	j, err := json.Marshal(b)
	if err != nil {
		fmt.Println("Error marshalling Book data into JSON", err)
		os.Exit(1)
	}
	fmt.Println(string(j))

	var bk Book
	bookJSON := `{"name":"other book","stars":3}`
	err = json.Unmarshal([]byte(bookJSON), &bk)
	if err != nil {
		fmt.Println("Error unmarshalling Book data", err)
		os.Exit(1)
	}
	fmt.Printf("%s", bk.Name)
}

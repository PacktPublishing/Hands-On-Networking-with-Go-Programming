package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	createFile()
	readFile()
	editFile()
	deleteFile()
}

func createFile() {
	f, err := os.Create("example.txt")
	if err != nil {
		log.Println("failed to create file", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString("File contents")
	if err != nil {
		log.Println("failed to write to file", err)
	}
}

func readFile() {
	f, err := os.Open("example.txt")
	if err != nil {
		log.Println("failed to open file", err)
		return
	}
	defer f.Close()
	data := make([]byte, 1024)
	n, err := f.Read(data)
	if err != nil {
		log.Println("failed to read data", err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", n, string(data[:n]))
}

func editFile() {
	f, err := os.OpenFile("example.txt", os.O_RDWR, 0)
	if err != nil {
		log.Println("failed to open file", err)
		return
	}
	defer f.Close()
	data := make([]byte, 1024)
	n, err := f.Read(data)
	if err != nil {
		log.Println("failed to read data", err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", n, string(data[:n]))
	_, err = f.WriteString("\nRead the data, but also added some")
	if err != nil {
		log.Println("failed to write data", err)
		return
	}
}

func deleteFile() {
	if err := os.Remove("example.txt"); err != nil {
		log.Println("failed to delete file", err)
	}
}

func pathJoinSplit() {
	joined := path.Join("/users/me", "more/paths", "filename.txt")
	dir, filename := path.Split(joined)
	fmt.Println(dir)
	fmt.Println(filename)
}

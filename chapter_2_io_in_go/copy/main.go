package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var sourceFlag = flag.String("source", "", "The source path")
var targetFlag = flag.String("target", "", "The target path")

func main() {
	flag.Parse()
	if *sourceFlag == "" {
		fmt.Println("Missing source flag")
		os.Exit(1)
	}
	if *targetFlag == "" {
		fmt.Println("Missing target flag")
		os.Exit(1)
	}

	// Check that the source file exists.
	if _, err := os.Stat(*sourceFlag); err == os.ErrNotExist {
		fmt.Println("Source file not found")
		os.Exit(1)
	}

	copy(*sourceFlag, *targetFlag)
}

func copy(from, to string) {
	src, err := os.Open(from)
	if err != nil {
		fmt.Println("couldn't open source", err)
		os.Exit(1)
	}
	defer src.Close()
	dst, err := os.Create(to)
	if err != nil {
		log.Println("couldn't create target", err)
		os.Exit(1)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println("couldn't copy data", err)
		os.Exit(1)
	}
}

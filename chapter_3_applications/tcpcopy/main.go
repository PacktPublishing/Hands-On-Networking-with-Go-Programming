package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/tcpcopy/receive"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/tcpcopy/send"
)

var modeFlag = flag.String("m", "recv", "Mode (send or recv)")
var targetFlag = flag.String("t", "", "Target host of the operation")
var keyFlag = flag.String("key", "", "4 words used as the encryption key")
var fileFlag = flag.String("f", "", "The path to the file to send")

func main() {
	flag.Parse()
	if *modeFlag == "send" {
		// Send
		fmt.Printf("Sending file %s to %s\n", *fileFlag, *targetFlag)
		err := send.Start(*fileFlag, *targetFlag, *keyFlag)
		exit(err)
	}
	// Receive.
	exit(receive.Start())
}

func exit(err error) {
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	fmt.Println("Complete.")
	os.Exit(0)
}

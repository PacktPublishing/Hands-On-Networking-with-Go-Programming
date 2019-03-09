package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"

	"golang.org/x/text/encoding/japanese"
)

func main() {
	convertToShiftJIS()
	convertFromGB2312UTF8()
}

func convertToShiftJIS() {
	// Convert from UTF-8 to Shift JIS
	src := bytes.NewBuffer([]byte{0xe3, 0x82, 0xa2, 0xe3, 0x82, 0xbf, 0xe3, 0x83, 0xaa}) // アタリ in UTF-8 encoding
	// Create an encoder.
	e := japanese.ShiftJIS.NewEncoder()
	// Wrap the output with the encoder.
	dst := new(bytes.Buffer)
	// Copy from the source to the destination, via the encoder.
	_, err := io.Copy(e.Writer(dst), src)
	if err != nil {
		fmt.Println("encoding error:", err)
		os.Exit(1)
	}
	// Print out the hex.
	for _, b := range dst.Bytes() {
		fmt.Printf("0x%2x ", b)
	}
	fmt.Println()
	// Outputs: 0x83 0x41 0x83 0x5e 0x83 0x8a
}

func convertFromGB2312UTF8() {
	src := bytes.NewBuffer([]byte{0xd6, 0xd0, 0xb9, 0xfa}) // 中国 in GB2312 encoding
	// Create a decoder.
	e := simplifiedchinese.GBK.NewDecoder()
	// Copy from the source to Stdout, via the decoder.
	_, err := io.Copy(os.Stdout, e.Reader(src))
	if err != nil {
		fmt.Println("decoding error:", err)
		os.Exit(1)
	}
	// Outputs: 中国
}

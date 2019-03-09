package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

func main() {
	fmt.Println("int64")
	ints := []int64{-2048, -1048, -512, -256, -128, -1, 0, 1, 128, 256, 512, 1024, 2048}
	for _, i := range ints {
		var b bytes.Buffer
		binary.Write(&b, binary.LittleEndian, i)
		bs := b.Bytes()
		fmt.Printf("%6d | %s | %s\n", i, hex(bs), bits(bs))
	}
	fmt.Println()

	fmt.Println("float64")
	floats := []float64{math.NaN(), math.Inf(-1), -128, -2.5, 2, -1, -0.5, 0, 0.5, 1, 2, 2.5, 128, math.Inf(1)}
	for _, f := range floats {
		var b bytes.Buffer
		binary.Write(&b, binary.LittleEndian, f)
		bs := b.Bytes()
		fmt.Printf("%7.2f | %s | %s\n", f, hex(bs), bits(bs))
	}
	fmt.Println()

	fmt.Println("bool")
	bools := []bool{false, true}
	for _, bl := range bools {
		var b bytes.Buffer
		binary.Write(&b, binary.LittleEndian, bl)
		bs := b.Bytes()
		fmt.Printf("%7v | %s | %s\n", bl, hex(bs), bits(bs))
	}
	fmt.Println()
}

func hex(v []byte) string {
	s := make([]string, len(v))
	for i, v := range v {
		s[i] = fmt.Sprintf("0x%02x", v)
	}
	return strings.Join(s, " ")
}

func bits(v []byte) string {
	s := make([]string, len(v))
	for i, v := range v {
		s[i] = fmt.Sprintf("%08b", v)
	}
	return strings.Join(s, " ")
}

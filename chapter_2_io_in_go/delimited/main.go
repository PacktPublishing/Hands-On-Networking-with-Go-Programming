package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, VariableWidth{
		Fixed: 1,
		// Variable: "this is a string",
		Variable: 3,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type VariableWidth struct {
	Fixed    int64
	Variable int64
}

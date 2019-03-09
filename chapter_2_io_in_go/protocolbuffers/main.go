package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a-h/gonp/02_io_in_go/protocolbuffers/messages"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func main() {
	var r messages.SearchResponse
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("unable to create timestamp", err)
	}
	r.Modified = ts
	r.Results = []*messages.SearchResult{
		&messages.SearchResult{
			Id:          1,
			Description: "description",
			Url:         "https://example.com/description/1",
		},
	}
	bytes, err := proto.Marshal(&r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Binary format.
	// 0x0a 0x32 0x08 0x01 0x12 0x21 0x68 0x74 0x74 0x70 0x73 0x3a 0x2f 0x2f 0x65 0x78 0x61 0x6d 0x70 0x6c
	// 0x65 0x2e 0x63 0x6f 0x6d 0x2f 0x64 0x65 0x73 0x63 0x72 0x69 0x70 0x74 0x69 0x6f 0x6e 0x2f 0x31 0x1a
	// 0x0b 0x64 0x65 0x73 0x63 0x72 0x69 0x70 0x74 0x69 0x6f 0x6e 0x12 0x0b 0x08 0xbe 0xe6 0x8a 0xe4 0x05
	// 0x10 0xb8 0x97 0xc6 0x3c
	fmt.Println(bytes)
	// Readable version
	// results: <
	//   id: 1
	//   url: "https://example.com/description/1"
	//   description: "description"
	// >
	// modified: <
	//   seconds: 1552069438
	//   nanos: 126979000
	// >
	proto.MarshalText(os.Stdout, &r)

	// Unmarshal back.
	var unmarshalled messages.SearchResponse
	err = proto.Unmarshal(bytes, &unmarshalled)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// {[id:1 url:"https://example.com/description/1" description:"description" ] seconds:1552069617 nanos:315187000  {} [] 0}
	fmt.Println(unmarshalled)
}

func hex(v []byte) string {
	s := make([]string, len(v))
	for i, v := range v {
		s[i] = fmt.Sprintf("0x%02x", v)
	}
	return strings.Join(s, " ")
}

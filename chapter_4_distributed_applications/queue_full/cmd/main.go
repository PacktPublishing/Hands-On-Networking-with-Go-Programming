package main

import (
	"fmt"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_4_distributed_applications/queue_full"
)

func main() {
	q := queue_full.New(10)

	for i := 0; i < 10; i++ {
		if err := q.Enqueue(i); err != nil {
			fmt.Println("got error:", err)
			return
		}
	}

	for {
		d, ok := q.Dequeue()
		if !ok {
			break
		}
		fmt.Println(d)
	}
}

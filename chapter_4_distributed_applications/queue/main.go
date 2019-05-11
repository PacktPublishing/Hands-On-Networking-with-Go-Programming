package main

import (
	"fmt"
	"strconv"
	"sync"
)

type data struct {
	Name string
}

func main() {
	// Make a queue that can hold up to 10 items before blocking.
	q := make(chan data, 10)

	// Post 10,000 items to the queue.
	go func() {
		for i := 0; i < 100000; i++ {
			q <- data{Name: strconv.Itoa(i)}
		}
	}()

	// Define a processing function.
	processor := func(processor int) {
		for {
			select {
			case d := <-q:
				fmt.Println(processor, "Received data", d)
			default:
				fmt.Println(processor, "Nothing to collect")
			}
		}
	}
	// Start up 2 processors.
	for i := 0; i < 2; i++ {
		go processor(i)
	}

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

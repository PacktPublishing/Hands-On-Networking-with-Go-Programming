package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

const data = `Name,Age
A1,10
A2,20
A3,30
A4,40
A5,50`

func main() {
	// You can get an io.Reader from many sources in Go. For example, you'll get a reader back
	// from a HTTP request or from S3.
	br := bytes.NewReader([]byte(data))

	// The io.Reader can be passed straight to the csv.NewReader function. This means that you
	// don't need to read all of the data into RAM in one go, and can start processing records
	// before you've finished transferring all of the data.
	r := csv.NewReader(br)
	var record []string
	var err error
	// You might want to skip the first record if it contains headings.
	// Simplest way is to just do `r.Read` once.
	_, err = r.Read()
	if err != nil {
		fmt.Println("failed to read CSV headings", err)
		// It's not usually a good idea to do this, since it stops your entire process server, but for an
		// example program, it's OK.
		os.Exit(1)
	}

	// A typical pattern here is to start one or more worker routines to process the records, using
	// a channel to marshal the communication. A channel is basically queue. This one will block on input
	// i.e. code will wait for the queue to empty before pushing another item into the queue.
	records := make(chan Record)

	// The WaitGroup can be used to wait for processes to complete. It's basically a thread-safe counter.
	workerCount := 2
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		// The go statement starts this function as a new goroutine which will run asynchronously.
		// The range over the channel means that it will run until the channel is closed.
		// Go can run hundreds of thousands of goroutines concurrently, but for making calls out to
		// databases etc., we'll probably want to limit concurrency. One simple way is to only start
		// a few workers, but there's also a semaphore package.

		// Note that the anonymous function that we're starting as a goroutine can have parameters passed
		// to it, but can also take closures (access variables in the parent scope). It's important to
		// use variables in the parent scope in a thread-safe way, so typically Go's concurrency features
		// are used to do this. From highest level to lowest level there are channels, mutexes and atomic
		// operations.
		go func(workerIndex int) {
			// The defer statement means that whatever path your code takes, wg.Done
			// will always be called when this function ends. It's going to end when the
			// channel is closed by the CSV reading code coming up next.
			defer wg.Done()
			for r := range records {
				//TODO: Process the record, I'll just print it.
				fmt.Printf("Index: %d: Name: %s, Age: %d\n", workerIndex, r.Name, r.Age)
			}
		}(i)
	}

	// Start processing the CSV data.
	for record, err = r.Read(); err != io.EOF; record, err = r.Read() {
		// Push the record to the channel and keep reading. One of the workers will collect it.
		age, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			// TODO: Decide what to do with the fact that we can't parse this record, e.g. skip it.
		}
		records <- NewRecord(record[0], int(age))
	}
	// Close the channel to let the workers know to exit their work.
	close(records)
	// At the end, you always get an io.EOF "error" to tell you you've reached the end of the
	// file.
	if err != io.EOF {
		// Do something if it's not EOF.
		fmt.Println("failed to read CSV data", err)
	}
	// Wait for the routines to finish.
	wg.Wait()
	fmt.Println("Done")
}

// NewRecord creates a new Record with required fields populated.
func NewRecord(name string, age int) Record {
	return Record{
		Name: name,
		Age:  age,
	}
}

// Record within the CSV.
type Record struct {
	Name string
	Age  int
}

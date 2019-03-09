package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	var s http.Server
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure we wait for all requests to complete.
		wg.Add(1)
		defer wg.Done()
		// This is a really slow HTTP request, it takes 1 minute to complete.
		// In the real world, you might have to make multiple calls in sequence.
		// If, after the first call, the client is no longer connected,
		// why bother making the second?
		fmt.Println("Server: simulating first API call...")
		<-time.After(time.Second * 10)

		if r.Context().Err() == context.Canceled {
			fmt.Println("Server: skipping second API call because the client is no longer connected...")
			w.WriteHeader(http.StatusTeapot)
			w.Write([]byte("Server: request cancelled by client."))
			return
		}

		fmt.Println("Server: simulating second API call...")
		<-time.After(time.Millisecond * 500)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server: complete!"))
	})
	s.Addr = "localhost:8090"

	// Start serving.
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.ListenAndServe()
	}()

	fmt.Println("Client: waiting for server to start")
	time.Sleep(time.Second * 5)
	fmt.Println("Client: making request with 1 second timeout")
	req, err := http.NewRequest("GET", "http://localhost:8090", nil)
	if err != nil {
		fmt.Println("Client:", err)
		return
	}
	var c http.Client
	c.Timeout = time.Second * 1
	_, err = c.Do(req)
	if err != nil {
		// It should have timed out.
		fmt.Println("Client:", err)
	}
	if err := s.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Waiting for server to gracefully shutdown")
	wg.Wait()
	fmt.Println("Complete")
}

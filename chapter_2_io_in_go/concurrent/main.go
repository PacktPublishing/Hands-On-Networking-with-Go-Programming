package main

import (
	"fmt"
	"net/http"
	"time"
)

type result struct {
	URL       string
	Status    int
	TimeTaken time.Duration
	Err       error
}

func main() {
	urls := []string{"http://google.com", "http://wikipedia.com"}
	results := make(chan result)
	for _, u := range urls {
		u := u
		go func() {
			start := time.Now()
			resp, err := http.Get(u)
			r := result{
				URL:       u,
				Err:       err,
				TimeTaken: time.Now().Sub(start),
			}
			if err == nil {
				r.Status = resp.StatusCode
			}
			results <- r
		}()
	}
	for i := 0; i < len(urls); i++ {
		r := <-results
		fmt.Println(r)
	}
}

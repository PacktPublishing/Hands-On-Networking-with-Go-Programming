package main

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
)

func main() {
	var attempt int
	lastAttempt := time.Now()
	operation := func() error {
		attempt++
		fmt.Println("Attempt", attempt, time.Now().Sub(lastAttempt))
		lastAttempt = time.Now()
		return fmt.Errorf("failed again")
	}
	bo := backoff.WithMaxTries(backoff.NewExponentialBackOff(), 10)
	err := backoff.Retry(operation, bo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}

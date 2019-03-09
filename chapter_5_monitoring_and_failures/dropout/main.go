package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		wg.Add(1)
		defer wg.Done()
		countUntilCancelled(ctx)
	}()
	fmt.Println("Press enter to cancel....")
	bufio.NewReader(os.Stdin).ReadString('\n')
	cancel()
	fmt.Println("Cancel received, waiting for graceful shutdown")
	wg.Wait()
	fmt.Println("Graceful shutdown complete")
}

func countUntilCancelled(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var i int64
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker received cancellation")
			return
		case <-ticker.C:
			i++
			fmt.Println(i)
		}
	}
}

package example5

import (
	"context"
	"fmt"
	"time"
)

// DoWork simulates a job that takes 50ms to complete, but can be cancelled via context
func DoWork(ch chan<- string) {
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "work done"
	}()
}

// Example: Wait for work to finish, but cancel if it takes longer than 10ms
func ExampleWorkWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ch := make(chan string)

	DoWork(ch)

	select {
	case result := <-ch:
		fmt.Println("[Main] result:", result)
	case <-ctx.Done():
		fmt.Println("[Main] result: work cancelled")
	}
}

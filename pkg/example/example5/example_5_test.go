package example5

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkDoWorkUnbuffered(b *testing.B) {
	before := runtime.NumGoroutine()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		ch := make(chan string, 1)
		DoWork(ch)
		select {
		case <-ch:
			// work done
		case <-ctx.Done():
			// work cancelled
		}
		cancel()
	}
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Println("\n\nGoroutine sayısı (before):", before)
	fmt.Println("Goroutine sayısı (after):", after)
}

func BenchmarkDoWorkBuffered(b *testing.B) {
	before := runtime.NumGoroutine()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		ch := make(chan string, 1)
		DoWork(ch)
		select {
		case <-ch:
			// work done
		case <-ctx.Done():
			// work cancelled
		}
		cancel()
	}
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Println("\n\nGoroutine sayısı (before):", before)
	fmt.Println("Goroutine sayısı (after):", after)
}

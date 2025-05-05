package example5

import (
	"context"
	"testing"
	"time"
)

func BenchmarkDoWorkUnbuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		ch := make(chan string)
		DoWork(ch)
		select {
		case <-ch:
			// work done
		case <-ctx.Done():
			// work cancelled
		}
		cancel()
	}
}

func BenchmarkDoWorkBuffered(b *testing.B) {
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
}

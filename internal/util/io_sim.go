package util

import (
	"fmt"
	"math/rand"
	"time"
)

// SimulateIO simulates an IO-bound wait with a random delay between 0-40ms
func SimulateIO(ms int) {
	delay := time.Duration(rand.Intn(ms+1)) * time.Millisecond
	ch := make(chan struct{})
	select {
	case <-time.After(delay):
		// IO completed (timeout simulates IO response)
	case <-ch:
		// This will never happen, just for completeness
	}
}

// SimulateError returns an error with the given probability (0.0-1.0). If no error, returns nil.
// Example: probability=0.05 means 5% chance to return error.
func SimulateError(probability float64, errMsg string) error {
	if rand.Float64() < probability {
		return fmt.Errorf(errMsg)
	}
	return nil
}

// Fast dot product for []float64 slices (manual, not SIMD)
// If you want to use SIMD, you can use github.com/minio/simd but it only supports float32.
func FastDot(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("FastDot: slice lengths do not match")
	}
	var sum float64
	for i := 0; i < len(a); i++ {
		sum += a[i] * b[i]
	}
	return sum
}

// Example for SIMD usage (float32 only, requires github.com/minio/simd)
// func FastDotSIMD(a, b []float32) float32 {
// 	return simd.Dot32(a, b)
// }

// SimulateCPUBoundFor simulates a CPU-bound workload for the given duration (in milliseconds).
func SimulateCPUBoundFor(ms int) {
	end := time.Now().Add(time.Duration(ms) * time.Millisecond)
	var x float64 = 1.0001
	for time.Now().Before(end) {
		x = x * 1.0001
		x = x / 1.0001
	}
}

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

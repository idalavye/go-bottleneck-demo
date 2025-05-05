package example4

import (
	"testing"
)

// Benchmark for value slice
func BenchmarkSumScoresValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SumScoresValue()
	}
}

// Benchmark for pointer slice
func BenchmarkSumScoresPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SumScoresPointer()
	}
}

func BenchmarkSumFixedArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SumFixedArray()
	}
}

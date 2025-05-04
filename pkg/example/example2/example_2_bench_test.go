package example2

import (
	"math/rand"
	"testing"
)

// generateMatrix generates a 2D slice with given dimensions and random values
func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(100)
		}
	}
	return matrix
}

/*
go test -bench=BenchmarkParallelRowBasedSum -run=^$ ./pkg/example/example2 -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof ./mem.prof
go tool trace trace.out
*/
func BenchmarkParallelRowBasedSum(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelRowBasedSum(matrix)
	}
}

/*
go test -bench=BenchmarkWorkerPoolRowBasedSum_1Worker -run=^$ ./pkg/example/example2 -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof ./mem.prof
go tool trace trace.out
*/
func BenchmarkWorkerPoolRowBasedSum_1Worker(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSum(matrix, 1)
	}
}

func BenchmarkWorkerPoolRowBasedSum_4Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSum(matrix, 4)
	}
}

func BenchmarkWorkerPoolRowBasedSum_8Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSum(matrix, 8)
	}
}

func BenchmarkWorkerPoolRowBasedSum_16Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSum(matrix, 16)
	}
}

/*
go test -bench=BenchmarkWorkerPoolRowBasedSumNoLock_1Worker -run=^$ ./pkg/example/example2 -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof ./cpu.prof
go tool trace trace.out
*/
func BenchmarkWorkerPoolRowBasedSumNoLock_1Worker(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSumNoLock(matrix, 1)
	}
}

func BenchmarkWorkerPoolRowBasedSumNoLock_4Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSumNoLock(matrix, 4)
	}
}

func BenchmarkWorkerPoolRowBasedSumNoLock_8Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSumNoLock(matrix, 8)
	}
}

func BenchmarkWorkerPoolRowBasedSumNoLock_16Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WorkerPoolRowBasedSumNoLock(matrix, 16)
	}
}

/*
go test -bench=BenchmarkChunkedRowBasedSum_1Worker -run=^$ ./pkg/example/example2 -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof ./mem.prof
go tool trace trace.out
*/
func BenchmarkChunkedRowBasedSum_1Worker(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChunkedRowBasedSum(matrix, 1)
	}
}

func BenchmarkChunkedRowBasedSum_4Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChunkedRowBasedSum(matrix, 4)
	}
}

func BenchmarkChunkedRowBasedSum_8Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChunkedRowBasedSum(matrix, 8)
	}
}

func BenchmarkChunkedRowBasedSum_16Workers(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChunkedRowBasedSum(matrix, 16)
	}
}

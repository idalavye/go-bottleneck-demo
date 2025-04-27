package search

import (
	"testing"
)

/*
goos: darwin
goarch: arm64
pkg: github.com/idagdelen/go-bottlenecks/internal/search
=== RUN   BenchmarkSumRowMajor
BenchmarkSumRowMajor
BenchmarkSumRowMajor-10             4507            252860 ns/op               0 B/op              0 allocs/op
PASS
ok      github.com/idagdelen/go-bottlenecks/internal/search     1.887s
*/
func BenchmarkSumRowMajor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SumRowMajor(productVectors)
	}
}

/*
goos: darwin
goarch: arm64
pkg: github.com/idagdelen/go-bottlenecks/internal/search
BenchmarkSumColMajor-10    	    3849	    322698 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/idagdelen/go-bottlenecks/internal/search	2.005s
*/
func BenchmarkSumColMajor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SumColMajor(productVectors)
	}
}

/*
goos: darwin
goarch: arm64
pkg: github.com/idagdelen/go-bottlenecks/internal/search
BenchmarkParallelSumRowMajor-10    	      25	  46874900 ns/op	 7265352 B/op	  200142 allocs/op
PASS
ok  	github.com/idagdelen/go-bottlenecks/internal/search	2.323s
*/
/*
GOGC=off go test -bench=BenchmarkParallelRowSumWithMerge -run=^$ ./internal/search -cpu 1 -benchtime 3s
GOGC=off go test -bench=BenchmarkParallelRowSumWithMerge -run=^$ ./internal/search -cpuprofile=cpu.prof -memprofile=mem.prof -trace=trace.out -cpu 1 -benchtime 3s
go tool pprof -http=:8080 ./cpu.prof
go tool pprof -http=:8080 ./mem.prof
go tool trace trace.out
*/
func BenchmarkParallelRowSumWithMerge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParallelRowSumWithMerge(productVectors)
	}
}

/*
goos: darwin
goarch: arm64
pkg: github.com/idagdelen/go-bottlenecks/internal/search
BenchmarkParallelRowSumWithWorkers-10    	      98	  12119224 ns/op	 3203459 B/op	      13 allocs/op
PASS
ok  	github.com/idagdelen/go-bottlenecks/internal/search	2.778s
*/
/*
GOGC=off go test -bench=BenchmarkParallelRowSumWithWorkers -run=^$ ./internal/search -cpu 1 -benchtime 3s
GOGC=off go test -bench=BenchmarkParallelRowSumWithWorkers -run=^$ ./internal/search -cpuprofile=cpu.prof -memprofile=mem.prof -trace=trace.out -cpu 1 -benchtime 3s
go tool pprof -http=:8080 ./cpu.prof
go tool pprof -http=:8080 ./mem.prof
go tool trace trace.out
*/
func BenchmarkParallelRowSumWithWorkers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParallelRowSumWithWorkers(productVectors)
	}
}

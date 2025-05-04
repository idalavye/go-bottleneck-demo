package example2

import (
	"sync"
)

// ParallelRowBasedSum sums the matrix row by row, but each row is processed in a separate goroutine (unlimited goroutines)
func ParallelRowBasedSum(matrix [][]int) int {
	sum := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < len(matrix); i++ {
		wg.Add(1)
		go func(row []int) {
			defer wg.Done()
			rowSum := 0
			for _, v := range row {
				rowSum += v
			}
			mu.Lock()
			sum += rowSum
			mu.Unlock()
		}(matrix[i])
	}

	wg.Wait()
	return sum
}

// WorkerPoolRowBasedSum sums the matrix row by row using a worker pool
func WorkerPoolRowBasedSum(matrix [][]int, workerCount int) int {
	sum := 0
	var mu sync.Mutex
	rowCh := make(chan []int)
	var wg sync.WaitGroup

	// Worker function
	worker := func() {
		defer wg.Done()
		for row := range rowCh {
			rowSum := 0
			for _, v := range row {
				rowSum += v
			}
			mu.Lock()
			sum += rowSum
			mu.Unlock()
		}
	}

	// Start workers
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	// Send rows to workers
	for i := 0; i < len(matrix); i++ {
		rowCh <- matrix[i]
	}
	close(rowCh)

	wg.Wait()
	return sum
}

// WorkerPoolRowBasedSumNoLock sums the matrix row by row using a worker pool without any lock
// Each worker keeps its own local sum and sends it to the result channel
func WorkerPoolRowBasedSumNoLock(matrix [][]int, workerCount int) int {
	rowCh := make(chan []int)
	resultCh := make(chan int, workerCount)
	var wg sync.WaitGroup

	worker := func() {
		localSum := 0
		for row := range rowCh {
			for _, v := range row {
				localSum += v
			}
		}
		resultCh <- localSum
		wg.Done()
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	for i := 0; i < len(matrix); i++ {
		rowCh <- matrix[i]
	}
	close(rowCh)

	wg.Wait()
	close(resultCh)

	total := 0
	for s := range resultCh {
		total += s
	}
	return total
}

// ChunkedRowBasedSum sums the matrix row by row by dividing the work into chunks for each worker
// No lock or channel is used, each worker processes a range of rows and writes its result to a slice
func ChunkedRowBasedSum(matrix [][]int, workerCount int) int {
	chunkSize := (len(matrix) + workerCount - 1) / workerCount
	results := make([]int, workerCount)
	var wg sync.WaitGroup

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func(workerIdx int) {
			defer wg.Done()
			start := workerIdx * chunkSize
			end := start + chunkSize
			if end > len(matrix) {
				end = len(matrix)
			}
			localSum := 0
			for i := start; i < end; i++ {
				for _, v := range matrix[i] {
					localSum += v
				}
			}
			results[workerIdx] = localSum
		}(w)
	}
	wg.Wait()

	total := 0
	for _, s := range results {
		total += s
	}
	return total
}

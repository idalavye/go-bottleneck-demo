package example1

// Row-based traversal: iterate row by row
func RowBasedSum(matrix [][]int) int {
	sum := 0
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

// Column-based traversal: iterate column by column
func ColumnBasedSum(matrix [][]int) int {
	sum := 0
	if len(matrix) == 0 {
		return 0
	}
	for j := 0; j < len(matrix[0]); j++ {
		for i := 0; i < len(matrix); i++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

package main

import (
	"fmt"
	"hash/fnv"
)

// Product struct for metadata
// This struct is used to represent product metadata
// (ID and Name fields)
type Product struct {
	ID   int
	Name string
}

// generateVector generates a deterministic vector for a given product name
func generateVector(name string) []float64 {
	h := fnv.New64a()
	h.Write([]byte(name))
	hash := h.Sum64()
	return []float64{
		float64(int64(hash>>0)&0xFFFF)/65536.0*2 - 1,
		float64(int64(hash>>16)&0xFFFF)/65536.0*2 - 1,
		float64(int64(hash>>32)&0xFFFF)/65536.0*2 - 1,
		float64(int64(hash>>48)&0xFFFF)/65536.0*2 - 1,
	}
}

func mock_data_gen() {
	const productCount = 100_000

	fmt.Println("// Bu dosya otomatik olarak scripts/mock_data_gen.go scripti ile üretilmiştir.")
	fmt.Println("// Elle düzenlemeyiniz.")
	fmt.Println()
	fmt.Println("package search")
	fmt.Println()
	fmt.Println("var productVectors = [][]float64{")
	for i := 1; i <= productCount; i++ {
		name := fmt.Sprintf("Otomatik Ürün %d", i)
		vec := generateVector(name)
		fmt.Printf("\t{%.2f, %.2f, %.2f, %.2f},\n", vec[0], vec[1], vec[2], vec[3])
	}
	fmt.Println("}")
	fmt.Println()
	fmt.Println("var productMetadata = []Product{")
	for i := 1; i <= productCount; i++ {
		fmt.Printf("\t{ID: %d, Name: \"Otomatik Ürün %d\"},\n", i, i)
	}
	fmt.Println("}")
}

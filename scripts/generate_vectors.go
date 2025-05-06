package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type QdrantPoint struct {
	ID      int         `json:"id"`
	Vector  []float32   `json:"vector"`
	Payload interface{} `json:"payload"`
}

type QdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	productCount := 100_000
	vectorSize := 4

	points := make([]QdrantPoint, 0, productCount)
	for i := 1; i <= productCount; i++ {
		vec := make([]float32, vectorSize)
		for j := 0; j < vectorSize; j++ {
			vec[j] = rand.Float32()
		}
		points = append(points, QdrantPoint{
			ID:     i,
			Vector: vec,
			Payload: map[string]interface{}{
				"name": fmt.Sprintf("Product %d", i),
				"desc": fmt.Sprintf("Description for product %d", i),
			},
		})
	}

	upsert := QdrantUpsertRequest{Points: points}

	file, err := os.Create("products_vectors.json")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(upsert); err != nil {
		fmt.Println("Failed to encode JSON:", err)
		os.Exit(1)
	}

	fmt.Println("Generated products_vectors.json with", productCount, "products in Qdrant upsert format.")
}

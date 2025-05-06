package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	const (
		filename  = "products_vectors.json"
		batchSize = 1000
		qdrantURL = "http://localhost:6333/collections/products/points?wait=true"
	)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		os.Exit(1)
	}

	var upsert QdrantUpsertRequest
	if err := json.Unmarshal(data, &upsert); err != nil {
		fmt.Println("Invalid JSON format:", err)
		os.Exit(1)
	}

	total := len(upsert.Points)
	fmt.Printf("Uploading %d products to Qdrant in batches of %d...\n", total, batchSize)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := upsert.Points[i:end]
		batchReq := QdrantUpsertRequest{Points: batch}
		body, err := json.Marshal(batchReq)
		if err != nil {
			fmt.Println("Failed to marshal batch:", err)
			os.Exit(1)
		}
		req, err := http.NewRequest("PUT", qdrantURL, bytes.NewReader(body))
		if err != nil {
			fmt.Printf("Batch %d-%d: Request creation error: %v\n", i+1, end, err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Batch %d-%d: HTTP error: %v\n", i+1, end, err)
			os.Exit(1)
		}
		resp.Body.Close()
		fmt.Printf("Batch %d-%d uploaded\n", i+1, end)
	}

	fmt.Println("All products uploaded to Qdrant!")
}

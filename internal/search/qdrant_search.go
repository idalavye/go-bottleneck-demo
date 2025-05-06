package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type QdrantSearchResult struct {
	ID    int     `json:"id"`
	Score float32 `json:"score"`
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
}

type qdrantSearchResponse struct {
	Result []struct {
		ID      int     `json:"id"`
		Score   float32 `json:"score"`
		Payload struct {
			Name string `json:"name"`
			Desc string `json:"desc"`
		} `json:"payload"`
	} `json:"result"`
}

func searchProductsQdrant(vector []float32, top int) ([]QdrantSearchResult, error) {
	url := "http://localhost:6333/collections/products/points/search"
	reqBody := map[string]interface{}{
		"vector": vector,
		"top":    top,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http post error: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	var qResp qdrantSearchResponse
	if err := json.Unmarshal(respBody, &qResp); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	results := make([]QdrantSearchResult, 0, len(qResp.Result))
	for _, r := range qResp.Result {
		results = append(results, QdrantSearchResult{
			ID:    r.ID,
			Score: r.Score,
			Name:  r.Payload.Name,
			Desc:  r.Payload.Desc,
		})
	}
	return results, nil
}

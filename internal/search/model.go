package search

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ScoredProduct struct {
	*Product
	Score float64 `json:"score"`
}

package search

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ScoredProduct struct {
	Product
	Score float64 `json:"score"`
}

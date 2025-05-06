package api

import (
	"github.com/idagdelen/go-bottlenecks/internal/search"
)

type searchResult struct {
	products []search.ScoredProduct
	totalSum float64
}

type enrichResult struct {
	enrichedProducts  []EnrichedProduct
	recommendedAdResp *EnrichedProduct
}

type job struct {
	idx  int
	prod search.ScoredProduct
}
type result struct {
	idx  int
	item *EnrichedProduct
}

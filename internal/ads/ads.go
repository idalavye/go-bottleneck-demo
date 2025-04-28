package ads

// RecommendProductByIDs simulates recommending a single product based on a list of IDs.
// It fetches product and stock info from simulated services and returns the first available product.
// This is a PoC for demonstrating bottlenecks in sequential vs concurrent code in Go.

import (
	"math/rand"

	"github.com/idagdelen/go-bottlenecks/internal/product"
	"github.com/idagdelen/go-bottlenecks/internal/stock"
	"github.com/idagdelen/go-bottlenecks/internal/util"
)

// RecommendedProduct holds product and stock info for recommendation
// This struct is used to return combined info from both services.
type RecommendedProduct struct {
	Product *product.Product
	Stock   *stock.Stock
}

// AdsService provides product recommendations based on IDs
// It uses product and stock services to fetch data.
type AdsService struct {
	ProductService product.ProductService
	StockService   stock.StockService
}

// NewAdsService creates a new AdsService with given product and stock services
func NewAdsService(productService product.ProductService, stockService stock.StockService) *AdsService {
	return &AdsService{
		ProductService: productService,
		StockService:   stockService,
	}
}

// RecommendProductByIDs fetches product and stock info for each ID and returns the first available product
func (a *AdsService) RecommendProductByIDs(ids []int) (*RecommendedProduct, error) {

	id, err := a.GetRandomRecommendedID(ids)
	if err != nil {
		return nil, err
	}
	// 5% error simulation (generic)
	if err := util.SimulateError(0.05, "network error: failed to fetch product"); err != nil {
		return nil, err
	}

	prod, err := a.ProductService.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	stk, err := a.StockService.GetStockByProductID(id)
	if err != nil {
		return nil, err
	}

	// Only recommend if stock is available
	if stk.Quantity > 0 {
		return &RecommendedProduct{
			Product: prod,
			Stock:   stk,
		}, nil
	}
	return nil, nil // No available product found
}

// GetRandomRecommendedID simulates a network call and returns a random ID from the list
func (s *AdsService) GetRandomRecommendedID(ids []int) (int, error) {
	if len(ids) == 0 {
		return 0, nil // or error if you prefer
	}
	util.SimulateIO(40)

	idx := rand.Intn(len(ids))
	return ids[idx], nil
}

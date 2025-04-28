package stock

import (
	"math/rand"
	"time"

	"github.com/idagdelen/go-bottlenecks/internal/util"
)

// Stock represents the stock information for a product
// (In a real scenario, this could include more fields)
type Stock struct {
	ProductID int
	Quantity  int
}

// StockService defines the interface for fetching stock details
// as if from an external service.
type StockService interface {
	GetStockByProductID(id int) (*Stock, error)
}

// SimulatedStockService simulates an external stock service.
type SimulatedStockService struct{}

// GetStockByProductID simulates a network call by sleeping for a random duration
// and returns a randomly generated stock quantity for the given product ID.
func (s *SimulatedStockService) GetStockByProductID(id int) (*Stock, error) {
	// Simulate network latency between 0ms and 40ms
	util.SimulateIO(40)

	// 5% error simulation (generic)
	if err := util.SimulateError(0.05, "network error: failed to fetch product"); err != nil {
		return nil, err
	}

	// Generate random stock quantity between 0 and 100
	stock := &Stock{
		ProductID: id,
		Quantity:  rand.Intn(101),
	}
	return stock, nil
}

// NewSimulatedStockService returns a new instance of SimulatedStockService
func NewSimulatedStockService() StockService {
	// Seed the random number generator (should be done once in main, but for PoC it's ok here)
	rand.Seed(time.Now().UnixNano())
	return &SimulatedStockService{}
}

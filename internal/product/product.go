package product

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/idagdelen/go-bottlenecks/internal/util"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
}

// ProductService defines the interface for fetching product details
// as if from an external service.
type ProductService interface {
	GetProductByID(id int) (*Product, error)
}

// SimulatedProductService simulates an external product service.
type SimulatedProductService struct{}

// GetProductByID simulates a network call by sleeping for a random duration
// and returns a randomly generated product. 5% of the time, returns an error.
func (s *SimulatedProductService) GetProductByID(id int) (*Product, error) {
	// Simulate network latency between 0ms and 40ms
	util.SimulateIO(40)

	// 5% error simulation (generic)
	if err := util.SimulateError(0.05, "network error: failed to fetch product"); err != nil {
		return nil, err
	}

	// Generate random product data
	product := &Product{
		ID:          id,
		Name:        fmt.Sprintf("Product-%d", rand.Intn(1000)),
		Description: "This is a randomly generated product.",
		Price:       rand.Float64()*100 + 1, // Price between 1 and 100
	}
	return product, nil
}

// NewSimulatedProductService returns a new instance of SimulatedProductService
func NewSimulatedProductService() ProductService {
	// Seed the random number generator (should be done once in main, but for PoC it's ok here)
	rand.Seed(time.Now().UnixNano())
	return &SimulatedProductService{}
}

// FormatPrice returns the product price as a string with two decimals and ₺ suffix
func (p *Product) FormatPrice() string {
	return fmt.Sprintf("%.2f₺", p.Price)
}

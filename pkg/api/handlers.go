package api

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/trace"
	"strconv"
	"sync"
	"time"

	"github.com/idagdelen/go-bottlenecks/internal/ads"
	"github.com/idagdelen/go-bottlenecks/internal/product"
	"github.com/idagdelen/go-bottlenecks/internal/search"
	"github.com/idagdelen/go-bottlenecks/internal/stock"
)

// ProductService instance (should be injected in real apps, global for PoC)
var prodService = product.NewSimulatedProductService()

// StockService instance (should be injected in real apps, global for PoC)
var stockService = stock.NewSimulatedStockService()

// AdsService instance (should be injected in real apps, global for PoC)
var adsService = ads.NewAdsService(prodService, stockService)

// HandleSearch handles search
// @Summary Search Demo
// @Description Searches the vectorized database with the given keyword and enriches the results with external services
// @Tags search
// @Param term query string true "Keyword to search for"
// @Param itemCount query int false "Number of products to return (default: 10)"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/search [get]
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	// Prevent crash: recover from any panic in handler
	defer func() {
		if err := recover(); err != nil {
			// Log the panic (in real app use logger)
			println("Recovered from panic in HandleSearch:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Internal server error (panic)",
			})
		}
	}()

	ctx, task := trace.NewTask(r.Context(), "HandleSearch")
	defer task.End()
	trace.WithRegion(ctx, "HandleSearch-Request", func() {
		searchTerm := r.URL.Query().Get("term")
		itemCount := r.URL.Query().Get("itemCount")

		parsedItemCount, err := strconv.Atoi(itemCount)
		if err != nil {
			parsedItemCount = 10 // Default value
		}

		products, totalSum := search.SearchProducts(searchTerm, parsedItemCount)

		// Use enrichment and ad logic only if demoMode=on
		withEnrichment := true
		/* _, recommendedAdResp := enrichProductsWithDetailsAndAd(products, withEnrichment) */
		enrichedProducts, recommendedAdResp := enrichProductsWithDetailsAndAd(products, withEnrichment)

		resp := Response{
			Success: true,
			Message: "Search completed successfully",
			Data: map[string]interface{}{
				"result":        enrichedProducts,
				"totalSum":      totalSum,
				"recommendedAd": recommendedAdResp,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

// HandleHealthCheck handles health check requests
// @Summary Health Check
// @Description Returns API health status
// @Tags system
// @Produce json
// @Success 200 {object} Response
// @Router /api/health [get]
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Prevent crash: recover from any panic in handler
	defer func() {
		if err := recover(); err != nil {
			println("Recovered from panic in HandleHealthCheck:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Internal server error (panic)",
			})
		}
	}()

	resp := Response{
		Success: true,
		Message: "API is healthy",
		Data: map[string]interface{}{
			"status": "up",
			"time":   time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// enrichProductsWithDetailsAndAd enriches products with details and optionally adds a recommended ad
// If withEnrichment is false, returns empty enrichedProducts and nil recommendedAdResp
func enrichProductsWithDetailsAndAd(products []search.ScoredProduct, withEnrichment bool) ([]EnrichedProduct, *EnrichedProduct) {
	if !withEnrichment {
		return nil, nil
	}
	var enrichedProducts []EnrichedProduct
	var idList []int
	for _, p := range products {
		idList = append(idList, p.ID)
		prod, err := prodService.GetProductByID(p.ID)
		if err == nil && prod != nil {
			stk, _ := stockService.GetStockByProductID(p.ID) // Ignore stock error for PoC
			enrichedProducts = append(enrichedProducts, EnrichedProduct{
				ID:          prod.ID,
				Name:        prod.Name,
				Description: prod.Description,
				Price:       prod.FormatPrice(),
				Score:       p.Score, // Add the score from search
				Stock:       0,
			})
			if stk != nil {
				enrichedProducts[len(enrichedProducts)-1].Stock = stk.Quantity
			}
		}
	}

	// Get a recommended ad product using AdsService and the id list
	recommendedAd, _ := adsService.RecommendProductByIDs(idList)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0, // Not relevant for ad
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return enrichedProducts, recommendedAdResp
}

// enrichProductsWithDetailsAndAdParallel enriches products concurrently, intentionally causing memory inefficiency and potential goroutine leaks for demo purposes
// This function demonstrates:
// - Parallel enrichment with goroutines
// - Excessive memory allocation
// - Not cleaning up nil entries in maps
// - Potential goroutine leaks
//
// To observe GC and STW pauses, run the app with:
//
//	GODEBUG=gctrace=1 go run ./cmd/main.go
func enrichProductsWithDetailsAndAdParallel(products []search.ScoredProduct, withEnrichment bool) ([]EnrichedProduct, *EnrichedProduct) {
	if !withEnrichment {
		return nil, nil
	}

	type result struct {
		idx  int
		item *EnrichedProduct
	}

	var (
		enrichedProducts = make([]EnrichedProduct, len(products))
		idList           = make([]int, 0, len(products))
		resultsCh        = make(chan result, len(products))
		// Intentionally create a large map and do not clean up nil entries
		/* memLeakMap = make(map[int]*EnrichedProduct, len(products)*10) */
	)

	// Start a goroutine for each product
	for i, p := range products {
		idList = append(idList, p.ID)
		go func(idx int, prodID int, score float64) {
			// Intentionally allocate a large slice (memory waste)
			_ = make([]byte, 1024*1024*2) // 1MB per goroutine (was 1MB)

			prod, err := prodService.GetProductByID(prodID)
			if err == nil && prod != nil {
				stk, _ := stockService.GetStockByProductID(prodID)
				item := &EnrichedProduct{
					ID:          prod.ID,
					Name:        prod.Name,
					Description: prod.Description,
					Price:       prod.FormatPrice(),
					Score:       score,
					Stock:       0,
				}
				if stk != nil {
					item.Stock = stk.Quantity
				}
				// Intentionally do not clean up nil entries in the map
				/* memLeakMap[prodID] = item */
				resultsCh <- result{idx: idx, item: item}
			} else {
				println("goroutine error", idx, err)
				// Hatalı durumda da kanala nil gönder
				/* resultsCh <- result{idx: idx, item: nil} */
			}
		}(i, p.ID, p.Score)
	}

	// Collect results (some goroutines may leak and never send)
	collected := 0
	timeout := time.After(1 * time.Second)

loop:
	for collected < len(products) {
		select {
		case res := <-resultsCh:
			if res.item != nil {
				enrichedProducts[res.idx] = *res.item
			}
			collected++
		case <-timeout:
			// Timeout oldu, kalanları beklemeden çık
			println("Timeout: not all goroutines responded")
			break loop
		}
	}

	// Get a recommended ad product using AdsService and the id list
	recommendedAd, _ := adsService.RecommendProductByIDs(idList)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0, // Not relevant for ad
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return enrichedProducts, recommendedAdResp
}

// enrichProductsWithDetailsAndAdParallel enriches products concurrently, demonstrating how improper goroutine and memory management can cause leaks.
// This version fixes goroutine and memory leaks using WaitGroup and context for timeout.
func enrichProductsWithDetailsAndAdParallelImprovement(products []search.ScoredProduct, withEnrichment bool) ([]EnrichedProduct, *EnrichedProduct) {
	if !withEnrichment {
		return nil, nil
	}

	type result struct {
		idx  int
		item *EnrichedProduct
	}

	var (
		enrichedProducts = make([]EnrichedProduct, len(products))
		idList           = make([]int, 0, len(products))
		resultsCh        = make(chan result, len(products))
	)

	// Use context for timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Start a goroutine for each product
	for i, p := range products {
		idList = append(idList, p.ID)
		wg.Add(1)
		go func(ctx context.Context, idx int, prodID int, score float64) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// Context cancelled, do not proceed
				return
			default:
				// continue
			}

			prod, err := prodService.GetProductByID(prodID)
			if err == nil && prod != nil {
				stk, _ := stockService.GetStockByProductID(prodID)
				item := &EnrichedProduct{
					ID:          prod.ID,
					Name:        prod.Name,
					Description: prod.Description,
					Price:       prod.FormatPrice(),
					Score:       score,
					Stock:       0,
				}
				if stk != nil {
					item.Stock = stk.Quantity
				}
				resultsCh <- result{idx: idx, item: item}
			} else {
				println("goroutine error", idx, err)
				resultsCh <- result{idx: idx, item: nil}
			}
		}(ctx, i, p.ID, p.Score)
	}

	// Goroutine to close the channel after all workers are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	collected := 0
	for res := range resultsCh {
		if res.item != nil {
			enrichedProducts[res.idx] = *res.item
		}
		collected++
		if collected == len(products) {
			break
		}
	}

	// Get a recommended ad product using AdsService and the id list
	recommendedAd, _ := adsService.RecommendProductByIDs(idList)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0, // Not relevant for ad
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return enrichedProducts, recommendedAdResp
}

package api

import (
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
	ctx, task := trace.NewTask(r.Context(), "HandleSearch")
	defer task.End()

	TraceRegion(ctx, "HandleSearch-Request", func() {
		searchTerm := r.URL.Query().Get("term")
		itemCount := r.URL.Query().Get("itemCount")

		parsedItemCount, err := strconv.Atoi(itemCount)
		if err != nil {
			parsedItemCount = 10 // Default value
		}

		searchRes := TraceRegionWithResult(ctx, "SearchProducts-Region", func() searchResult {
			products, totalSum := search.SearchProductsHeapOptimized(searchTerm, parsedItemCount)
			return searchResult{products: products, totalSum: totalSum}
		})
		products := searchRes.products
		totalSum := searchRes.totalSum

		result := TraceRegionWithResult(ctx, "EnrichProducts-Region", func() enrichResult {
			eps, ad := enrichProductsWithDetailsAndAdWorkerPool(products)
			return enrichResult{enrichedProducts: eps, recommendedAdResp: ad}
		})

		enrichedProducts := result.enrichedProducts
		recommendedAdResp := result.recommendedAdResp

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
func enrichProductsWithDetailsAndAd(products []search.ScoredProduct) ([]EnrichedProduct, *EnrichedProduct) {
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

// sync.Pool for EnrichedProduct reuse across requests
var enrichedProductPool = sync.Pool{
	New: func() interface{} { return new(EnrichedProduct) },
}

func enrichProductsWithDetailsAndAdWorkerPool(products []search.ScoredProduct) ([]EnrichedProduct, *EnrichedProduct) {
	numWorkers := 10
	jobs := make(chan job, len(products))
	results := make(chan result, len(products))

	var enrichedProducts = make([]EnrichedProduct, 0, len(products))
	idList := make([]int, len(products))

	worker := func() {
		for j := range jobs {
			prod, err := prodService.GetProductByID(j.prod.ID)
			if err == nil && prod != nil {
				stk, _ := stockService.GetStockByProductID(j.prod.ID)
				item := enrichedProductPool.Get().(*EnrichedProduct)
				item.ID = prod.ID
				item.Name = prod.Name
				item.Description = prod.Description
				item.Price = prod.FormatPrice()
				item.Score = j.prod.Score
				item.Stock = 0
				if stk != nil {
					item.Stock = stk.Quantity
				}
				results <- result{idx: j.idx, item: item}
			} else {
				results <- result{idx: j.idx, item: nil}
			}
		}
	}

	for w := 0; w < numWorkers; w++ {
		go worker()
	}

	for i, p := range products {
		idList[i] = p.ID
		jobs <- job{idx: i, prod: p}
	}
	close(jobs)

	for i := 0; i < len(products); i++ {
		res := <-results
		if res.item != nil {
			enrichedProducts = append(enrichedProducts, *res.item)
		}
	}
	close(results)

	recommendedAd, _ := adsService.RecommendProductByIDs(idList)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0,
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return enrichedProducts, recommendedAdResp
}

func enrichProductsWithDetailsParallelAndIndexBased(products []search.ScoredProduct) ([]EnrichedProduct, *EnrichedProduct) {
	var wg sync.WaitGroup
	results := make([]*EnrichedProduct, len(products))
	idList := make([]int, len(products))

	for i, p := range products {
		idList[i] = p.ID
		wg.Add(1)
		go func(i int, p search.ScoredProduct) {
			defer wg.Done()

			prod, err := prodService.GetProductByID(p.ID)
			if err == nil && prod != nil {
				stk, _ := stockService.GetStockByProductID(p.ID)
				item := enrichedProductPool.Get().(*EnrichedProduct)
				item.ID = prod.ID
				item.Name = prod.Name
				item.Description = prod.Description
				item.Price = prod.FormatPrice()
				item.Score = p.Score
				item.Stock = 0
				if stk != nil {
					item.Stock = stk.Quantity
				}
				results[i] = item
			}
		}(i, p)
	}
	wg.Wait()

	finalResults := make([]EnrichedProduct, 0, len(products))
	for _, item := range results {
		if item != nil {
			finalResults = append(finalResults, *item)
		}
	}

	recommendedAd, _ := adsService.RecommendProductByIDs(idList)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0,
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return finalResults, recommendedAdResp
}

// enrichProductsWithDetailsAndAdIndexBased enriches products concurrently, each goroutine writes to its own index, then nils are cleaned up
func enrichProductsWithDetailsParallelAndIndexBased_v2(products []search.ScoredProduct) ([]EnrichedProduct, *EnrichedProduct) {
	var wg sync.WaitGroup
	results := make([]*EnrichedProduct, len(products))
	idList := make([]int, len(products))

	recommendedAdCh := make(chan *ads.RecommendedProduct, 1)
	go func(ids []int) {
		recommendedAd, _ := adsService.RecommendProductByIDs(ids)
		recommendedAdCh <- recommendedAd
	}(func() []int {
		ids := make([]int, len(products))
		for i, p := range products {
			ids[i] = p.ID
		}
		return ids
	}())

	for i, p := range products {
		idList[i] = p.ID
		wg.Add(1)

		go func(i int, p search.ScoredProduct) {
			defer wg.Done()

			var prod *product.Product
			var stk *stock.Stock
			var prodErr error
			/* var stockWg sync.WaitGroup */

			prod, prodErr = prodService.GetProductByID(p.ID)
			stk, _ = stockService.GetStockByProductID(p.ID)

			/* 	stockWg.Add(2)
			// Product goroutine
			go func() {
				defer stockWg.Done()
			}()
			// Stock goroutine
			go func() {
				defer stockWg.Done()
			}()
			stockWg.Wait() */

			if prodErr == nil && prod != nil {
				item := enrichedProductPool.Get().(*EnrichedProduct)
				item.ID = prod.ID
				item.Name = prod.Name
				item.Description = prod.Description
				item.Price = prod.FormatPrice()
				item.Score = p.Score
				item.Stock = 0
				if stk != nil {
					item.Stock = stk.Quantity
				}
				results[i] = item
			}
		}(i, p)
	}
	wg.Wait()

	finalResults := make([]EnrichedProduct, 0, len(products))
	for _, item := range results {
		if item != nil {
			finalResults = append(finalResults, *item)
		}
	}

	recommendedAd := <-recommendedAdCh
	close(recommendedAdCh)
	var recommendedAdResp *EnrichedProduct
	if recommendedAd != nil {
		recommendedAdResp = &EnrichedProduct{
			ID:          recommendedAd.Product.ID,
			Name:        recommendedAd.Product.Name,
			Description: recommendedAd.Product.Description,
			Price:       recommendedAd.Product.FormatPrice(),
			Score:       0,
			Stock:       recommendedAd.Stock.Quantity,
		}
	}
	return finalResults, recommendedAdResp
}

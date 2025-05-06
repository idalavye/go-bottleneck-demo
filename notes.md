Step 1 - Let's execute our load test to evaluate the performance

```sh
hey -z 10s -c 50 "http://localhost:8080/api/search?term=örnekKelime&itemCount=10"
```

Step 2 - Let's profile sync code

Step 3 - Let's parallel

```go
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
```

Step 4 - Lets add cancel message

```
else {
  /* resultsCh <- result{idx: idx, item: nil} */
}
```

```
 Block time (GC mark assist wait for work) | GC (Çöp Toplayıcı) işaretleme fazında, goroutine’in GC’ye yardım etmek için beklediği süre. |
| Block time (network) | Ağ (network) işlemleri sırasında (ör. veri beklerken) bloklanan süredir. |
| Block time (preempted) | Goroutine’in çalışma hakkı elinden alındığında (preempt) beklediği süredir. |
| Block time (select) | select ifadesinde (kanallardan veri beklerken) bloklanan süredir. |
| Block time (sync) | sync paketindeki kilitler (ör. Mutex, RWMutex) nedeniyle beklenen süredir. |
| Block time (sync.(Cond).Wait) | sync.Cond ile koşul değişkeni beklerken geçen süredir. |
| Block time (syscall) | Sistem çağrısı (ör. dosya, ağ, vs.) sırasında beklenen süredir. |
| Sched wait time | Goroutine’in scheduler tarafından çalıştırılmak için beklediği süredir. |
| Syscall execution time | Sistem çağrısının (ör. dosya okuma/yazma, ağ işlemi) gerçek çalıştırılma süresidir. |
| Unknown time | Sebebi belirlenemeyen veya sınıflandırılamayan bekleme süresidir. |


curl http://localhost:6060/debug/pprof/goroutine\?debug\=2 > goroutine.txt
```

```
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
			println("Timeout: not all goroutines responded")
			break loop
		}
	}
```

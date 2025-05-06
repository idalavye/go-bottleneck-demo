Index based search

```go
func enrichProductsWithDetailsAndAdIndexBased(products []search.ScoredProduct) ([]*EnrichedProduct, *EnrichedProduct) {
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

	// Temizlik: nil olanları çıkar
	finalResults := make([]*EnrichedProduct, 0, len(products))
	for _, item := range results {
		if item != nil {
			finalResults = append(finalResults, item)
		}
	}

	// Reklam ürünü
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
```

- En son herşey paralleştirelim;

```go
// enrichProductsWithDetailsAndAdIndexBased enriches products concurrently, each goroutine writes to its own index, then nils are cleaned up
func enrichProductsWithDetailsAndAdIndexBased(products []search.ScoredProduct) ([]*EnrichedProduct, *EnrichedProduct) {
	var wg sync.WaitGroup
	results := make([]*EnrichedProduct, len(products))
	idList := make([]int, len(products))

	recommendedAdCh := make(chan *ads.RecommendedProduct, 1)
	go func(ids []int) {
		recommendedAd, _ := adsService.RecommendProductByIDs(ids)
		recommendedAdCh <- recommendedAd
	}(func() []int { // idList henüz dolmadığı için kopyasını gönderiyoruz
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
			var stockWg sync.WaitGroup

			stockWg.Add(2)
			// Product goroutine
			go func() {
				defer stockWg.Done()
				prod, prodErr = prodService.GetProductByID(p.ID)
			}()
			// Stock goroutine
			go func() {
				defer stockWg.Done()
				stk, _ = stockService.GetStockByProductID(p.ID)
			}()
			stockWg.Wait()

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

	// Temizlik: nil olanları çıkar
	finalResults := make([]*EnrichedProduct, 0, len(products))
	for _, item := range results {
		if item != nil {
			finalResults = append(finalResults, item)
		}
	}

	// Reklam ürünü sonucunu bekle
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
```

```
Neden Bu Kadar Fark Var?
Worker pool yaklaşımı, CPU-bound işler için idealdir. Çünkü CPU sayısı kadar worker ile en verimli şekilde işlem yapılır.
Goroutine başına iş yaklaşımı, IO-bound işler için daha hızlıdır. Çünkü bekleyen çok sayıda iş varken, Go runtime bunları verimli şekilde planlar ve bekleyenler varken diğerlerini çalıştırır.
Senin örneğinde, ürün ve stok servisleri muhtemelen IO-bound (ör: sleep, network, vs.) olduğu için, binlerce goroutine açmak çok daha hızlı sonuç veriyor.
```

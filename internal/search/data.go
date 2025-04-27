package search

import (
	"hash/fnv"
)

func getEmbedding(text string) []float64 {
	h := fnv.New64a()
	h.Write([]byte(text))
	hash := h.Sum64()

	// Her arama terimi için aynı vektör, farklı terimler için farklı vektör üretir
	return []float64{
		float64(int64(hash>>0)&0xFFFF)/65536.0*2 - 1,  // -1 ile 1 arası
		float64(int64(hash>>16)&0xFFFF)/65536.0*2 - 1, // -1 ile 1 arası
		float64(int64(hash>>32)&0xFFFF)/65536.0*2 - 1, // -1 ile 1 arası
		float64(int64(hash>>48)&0xFFFF)/65536.0*2 - 1, // -1 ile 1 arası
	}
}

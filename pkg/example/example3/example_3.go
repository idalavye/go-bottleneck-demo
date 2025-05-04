package example3

import (
	"encoding/json"
	"sync"
)

type Product struct {
	ID     int                    `json:"id"`
	Name   string                 `json:"name"`
	Email  string                 `json:"email"`
	Tags   []string               `json:"tags"`
	Active bool                   `json:"active"`
	Score  float64                `json:"score"`
	Meta   map[string]interface{} `json:"meta"`
}

func parseJSONWithoutPool(data []byte) *Product {
	var s Product
	_ = json.Unmarshal(data, &s)
	return &s
}

var productPool = sync.Pool{
	New: func() interface{} {
		return new(Product)
	},
}

func parseJSONWithPool(data []byte) *Product {
	s := productPool.Get().(*Product)
	_ = json.Unmarshal(data, s)
	productPool.Put(s)
	return s
}

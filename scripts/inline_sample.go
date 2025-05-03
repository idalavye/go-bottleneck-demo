/*
Go'da inlining ve escape analysis örneği.
Aşağıdaki örnekte;
- Pointer döndüren ve value döndüren iki factory fonksiyonu var.
- Interface ile sarmaladığımızda değişkenin heap'e kaçtığını görebilirsin.
Interface ile sarmalanan value ve pointer örneklerinde, interface kullanımı nedeniyle ilgili değerler heap'e kaçıyor. Çünkü interface'ler altında dinamik olarak veri tutuluyor ve bu veri genellikle heap'te saklanıyor.
- escape analysis çıktısı için: go run -gcflags="-m" scripts/inline_sample.go
*/

package main

import "fmt"

// Basit bir struct
type Data struct {
	value int
}

// Value semantic factory
func NewDataValue(v int) Data {
	return Data{value: v} // Value döndürür, çoğunlukla stack'te kalır
}

// Pointer semantic factory
func NewDataPointer(v int) *Data {
	return &Data{value: v} // Pointer döndürür, inlining ile bazen stack'te kalır
}

// Interface ile sarmalama
type DataIface interface {
	Get() int
}

func (d Data) Get() int { return d.value }

func inline_sample() {
	// Value semantic
	d1 := NewDataValue(10)
	fmt.Println("Value semantic:", d1.Get())

	// Pointer semantic
	d2 := NewDataPointer(20)
	fmt.Println("Pointer semantic:", d2.Get())

	// Interface ile value
	var i1 DataIface = NewDataValue(30)
	fmt.Println("Interface + value semantic:", i1.Get())

	// Interface ile pointer
	var i2 DataIface = NewDataPointer(40)
	fmt.Println("Interface + pointer semantic:", i2.Get())
}

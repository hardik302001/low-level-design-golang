package main

type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

func NewProduct(id int, name string, price float64, quantity int) *Product {
	return &Product{ID: id, Name: name, Price: price, Quantity: quantity}
}

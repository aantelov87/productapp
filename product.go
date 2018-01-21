package exampleapp

import (
	"time"
)

// Product describes a product.
type Product struct {
	ID         string
	Name       string
	Desc       string
	Attributes []string
	Images     []string
	Created    time.Time
}

// ProductStore provides access to a product storage
type ProductStore interface {
	// List returns a list of produts.
	List() ([]*Product, error)

	// Lookup retrieves a book by its ID.
	Lookup(id string) (*Product, error)

	// Create saves a given book, assigning it a new ID.
	Create(o *Product) error
}

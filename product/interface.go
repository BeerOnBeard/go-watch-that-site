package product

// Product contains the name of the product and the direct URI to the product page
type Product struct {
	Name     string
	URI      string
	ImageURI string
}

// Finder finds products and returns them over the channel as they are found
// Use the completeChannel to notify consumers when search is complete
type Finder interface {
	Find(productChannel chan *Product, errorChannel chan *error, completeChannel chan bool)
}

// Storer provides a way to store and retrieve products
type Storer interface {
	Get() ([]*Product, error)
	Put([]*Product) error
}

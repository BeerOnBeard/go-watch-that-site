package scrapers

// Product contains the name of the product and the direct URI to the product page
type Product struct {
	Name string
	URI  string
}

// ProductFinder finds products and returns them over the channel as they are found
// Use the completeChannel to notify consumers when search is complete
type ProductFinder interface {
	FindProducts(productChannel chan *Product, errorChannel chan *error, completeChannel chan bool)
}

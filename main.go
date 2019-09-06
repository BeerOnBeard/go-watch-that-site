package main

import (
	"fmt"
	"net/http"

	"github.com/beeronbeard/go-watch-that-site/product"
	"github.com/beeronbeard/go-watch-that-site/product/finders/airborneoutlet"
	"github.com/beeronbeard/go-watch-that-site/product/finders/canyonoutlet"
	"github.com/beeronbeard/go-watch-that-site/product/storer"
)

const (
	productsFilePath  = "bikes"
	canyonOutletURI   = "https://www.canyon.com/en-us/outlet/mountain-bikes/"
	airborneOutletURI = "https://airbornebicycles.com/collections/outlet"
)

func main() {
	finders := []product.Finder{
		&canyonoutlet.CanyonOutlet{Client: http.DefaultClient, URI: canyonOutletURI},
		&airborneoutlet.AirborneOutlet{Client: http.DefaultClient, URI: airborneOutletURI},
	}

	productChannel := make(chan *product.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)

	for _, finder := range finders {
		go finder.Find(productChannel, errorChannel, completeChannel)
	}

	completeCount := 0

	var products []*product.Product
loop:
	for {
		select {
		case product := <-productChannel:
			products = append(products, product)
		case err := <-errorChannel:
			panic(err)
		case <-completeChannel:
			completeCount++
			if completeCount == len(finders) {
				break loop
			}
		}
	}

	storer := storer.New(productsFilePath)

	storedProducts, err := storer.Get()
	if err != nil {
		panic(err)
	}

	var newProducts []*product.Product
newProductsLoop:
	for i := 0; i < len(products); i++ {
		for j := 0; j < len(storedProducts); j++ {
			if products[i].Name == storedProducts[j].Name {
				continue newProductsLoop
			}
		}

		newProducts = append(newProducts, products[i])
	}

	var removedProducts []*product.Product
removedProductsLoop:
	for i := 0; i < len(storedProducts); i++ {
		for j := 0; j < len(products); j++ {
			if storedProducts[i].Name == products[j].Name {
				continue removedProductsLoop
			}
		}

		removedProducts = append(removedProducts, storedProducts[i])
	}

	if len(newProducts) > 0 || len(removedProducts) > 0 {
		fmt.Printf("New: %v", newProducts)
		fmt.Printf("Removed: %v", removedProducts)
	}

	err = storer.Put(products)
	if err != nil {
		panic(err)
	}
}

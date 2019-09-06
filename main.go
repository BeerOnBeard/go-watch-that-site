package main

import (
	"fmt"
	"net/http"

	"github.com/beeronbeard/go-watch-that-site/scrapers"
	"github.com/beeronbeard/go-watch-that-site/scrapers/airborneoutlet"
	"github.com/beeronbeard/go-watch-that-site/scrapers/canyonoutlet"
)

const (
	canyonOutletURI   = "https://www.canyon.com/en-us/outlet/mountain-bikes/"
	airborneOutletURI = "https://airbornebicycles.com/collections/outlet"
)

func main() {
	finders := []scrapers.ProductFinder{
		&canyonoutlet.CanyonOutlet{Client: http.DefaultClient, URI: canyonOutletURI},
		&airborneoutlet.AirborneOutlet{Client: http.DefaultClient, URI: airborneOutletURI},
	}

	productChannel := make(chan *scrapers.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)

	for _, finder := range finders {
		go finder.FindProducts(productChannel, errorChannel, completeChannel)
	}

	completeCount := 0

loop:
	for {
		select {
		case product := <-productChannel:
			fmt.Println(product)
		case err := <-errorChannel:
			fmt.Println(err)
		case <-completeChannel:
			completeCount++
			if completeCount == len(finders) {
				break loop
			}
		}
	}
}

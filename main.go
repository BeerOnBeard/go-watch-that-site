package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beeronbeard/go-watch-that-site/scrapers/airborneoutlet"
	"github.com/beeronbeard/go-watch-that-site/scrapers/canyonoutlet"
)

const (
	canyonOutletURI   = "https://www.canyon.com/en-us/outlet/mountain-bikes/"
	airborneOutletURI = "https://airbornebicycles.com/collections/outlet"
)

func main() {
	c := canyonoutlet.CanyonOutlet{Client: http.DefaultClient, URI: canyonOutletURI}
	a := airborneoutlet.AirborneOutlet{Client: http.DefaultClient, URI: airborneOutletURI}
	canyonProducts, err := c.FindProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	airborneProducts, err := a.FindProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Print(canyonProducts)
	fmt.Print(airborneProducts)
}

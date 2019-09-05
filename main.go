package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beeronbeard/go-watch-that-site/scrapers"
)

const canyonOutletURI = "https://www.canyon.com/en-us/outlet/mountain-bikes/"

func main() {
	c := scrapers.CanyonOutlet{Client: http.DefaultClient, URI: canyonOutletURI}
	products, err := c.FindProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Print(products)
}

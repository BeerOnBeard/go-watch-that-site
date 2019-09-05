package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beeronbeard/go-watch-that-site/scrapers"
	"golang.org/x/net/html"
)

const canyonOutletURI = "https://www.canyon.com/en-us/outlet/mountain-bikes/"

func main() {
	response, err := http.Get(canyonOutletURI)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	c := scrapers.CanyonOutlet{}
	products, err := c.FindProducts(doc)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Print(products)
}

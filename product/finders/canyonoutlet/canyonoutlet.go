package canyonoutlet

import (
	"net/http"
	"strings"

	"github.com/beeronbeard/go-watch-that-site/product"
	"golang.org/x/net/html"
)

// CanyonOutlet is a product finder for the Canyon Outlet
type CanyonOutlet struct {
	Client *http.Client
	URI    string
}

// Find products in the Canyon Outlet
func (finder *CanyonOutlet) Find(productChannel chan *product.Product, errorChannel chan *error, completeChannel chan bool) {
	response, err := finder.Client.Get(finder.URI)
	if err != nil {
		errorChannel <- &err
		completeChannel <- true
		return
	}

	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		errorChannel <- &err
		completeChannel <- true
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "productTile__link" {
					product := getProductInfo(n)
					productChannel <- &product
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	completeChannel <- true
}

func getProductInfo(node *html.Node) product.Product {
	var name, uri, imageURI string
	for _, attr := range node.Attr {
		if attr.Key == "aria-label" {
			name = attr.Val
		}

		if attr.Key == "href" {
			uri = attr.Val
		}
	}

	imageNode := node.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.LastChild.PrevSibling
	for _, attr := range imageNode.Attr {
		if attr.Key == "data-src" {
			imageURI = strings.Split(attr.Val, "?")[0]
		}
	}

	return product.Product{Name: name, URI: uri, ImageURI: imageURI}
}

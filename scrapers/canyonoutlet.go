package scrapers

import (
	"golang.org/x/net/html"
)

// CanyonOutlet is a scraper for the Canyon Outlet
type CanyonOutlet struct{}

// FindProducts searches the doc for Canyon Outlet products
func (CanyonOutlet) FindProducts(doc *html.Node) ([]Product, error) {
	var products []Product

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "productTile__link" {
					products = append(products, getProductInfo(n))
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return products, nil
}

func getProductInfo(node *html.Node) Product {
	var name, uri string
	for _, attr := range node.Attr {
		if attr.Key == "aria-label" {
			name = attr.Val
		}

		if attr.Key == "href" {
			uri = attr.Val
		}
	}

	return Product{name, uri}
}

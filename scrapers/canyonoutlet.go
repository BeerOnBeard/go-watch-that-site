package scrapers

import (
	"net/http"

	"golang.org/x/net/html"
)

// CanyonOutlet is a scraper for the Canyon Outlet
type CanyonOutlet struct {
	Client *http.Client
	URI    string
}

// FindProducts searches for Canyon Outlet products
func (c *CanyonOutlet) FindProducts() ([]Product, error) {
	var products []Product

	response, err := c.Client.Get(c.URI)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

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

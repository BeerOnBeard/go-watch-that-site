package airborneoutlet

import (
	"net/http"
	"strings"

	"github.com/beeronbeard/go-watch-that-site/scrapers"
	"golang.org/x/net/html"
)

// AirborneOutlet is a scraper for the Airborne Outlet
type AirborneOutlet struct {
	Client *http.Client
	URI    string
}

// FindProducts in the Airborne Outlet
func (a *AirborneOutlet) FindProducts() ([]scrapers.Product, error) {
	var products []scrapers.Product

	response, err := a.Client.Get(a.URI)
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
				if attr.Key == "class" && strings.Contains(attr.Val, "product-info__caption") {
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

/*
<a class="product-info__caption hidden" href="/collections/outlet/products/griffin-27-5-demo-bike" itemprop="url">
	<div class="product-details">
    <span class="title" itemprop="name">Griffin 27.5+ (Demo Bike)</span>
    <span class="shopify-product-reviews-badge" data-id="1351312572474"></span>
    <span class="price sale">
      <span class="money">$ 899.95</span>
      <span class="was_price">
        <span class="money">$ 1,350.95</span>
      </span>
    </span>
	</div>
</a>
*/
func getProductInfo(node *html.Node) scrapers.Product {
	var name, uri string

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			uri = attr.Val
			break
		}
	}

	titleNode := node.LastChild.FirstChild.NextSibling.FirstChild
	priceNode := node.LastChild.LastChild.PrevSibling.FirstChild.NextSibling.FirstChild
	name = titleNode.Data + " " + priceNode.Data

	return scrapers.Product{Name: name, URI: uri}
}

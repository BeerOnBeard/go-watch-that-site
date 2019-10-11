package airborneoutlet

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beeronbeard/go-watch-that-site/product"
)

var expectedProducts = []product.Product{
	product.Product{
		Name:     "Airborne Seeker 29",
		URI:      "https://airbornebicycles.com/collections/outlet/products/seeker",
		ImageURI: "https://cdn.shopify.com/s/files/1/1067/3142/products/airborne-seeker_900x.png?v=1519945441"},
	product.Product{
		Name:     "Airborne Griffin 27.5+",
		URI:      "https://airbornebicycles.com/collections/outlet/products/griffin-27-5-demo-bike",
		ImageURI: "https://cdn.shopify.com/s/files/1/1067/3142/products/airborne-griffin_c7099c00-f1dc-4f08-b756-70f89a3982b8_900x.png?v=1527618327"},
	product.Product{
		Name:     "Airborne Goblin 29",
		URI:      "https://airbornebicycles.com/collections/outlet/products/goblin-29-demo-bike",
		ImageURI: "https://cdn.shopify.com/s/files/1/1067/3142/products/GOBLIN29-Profile_3d5ad1a8-5b4e-451e-8f53-241a64dc3e76_900x.png?v=1557838909"},
	product.Product{
		Name:     "Airborne Goblin EVO 27.5",
		URI:      "https://airbornebicycles.com/collections/outlet/products/goblin-evo-27-5-demo-bike",
		ImageURI: "https://cdn.shopify.com/s/files/1/1067/3142/products/airborne-goblin-evo_017ccc6b-6af1-4b47-894f-9d78ca2b6fda_900x.png?v=1542738441"},
}

func TestFindProducts(t *testing.T) {
	data, err := ioutil.ReadFile("airborneoutlet.html")
	if err != nil {
		t.Fatalf("Could not read file. %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(data)
	}))
	defer server.Close()

	finder := AirborneOutlet{server.Client(), server.URL}
	productChannel := make(chan *product.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)
	go finder.Find(productChannel, errorChannel, completeChannel)

	var products []product.Product

loop:
	for {
		select {
		case product := <-productChannel:
			products = append(products, *product)
		case err := <-errorChannel:
			t.Fatal(err)
		case <-completeChannel:
			break loop
		}
	}

	for i, expectedProduct := range expectedProducts {
		if products[i] != expectedProduct {
			fmt.Println(products[i].Name)
			fmt.Println(expectedProduct.Name)
			fmt.Println(products[i].URI)
			fmt.Println(expectedProduct.URI)
			fmt.Println(products[i].ImageURI)
			fmt.Println(expectedProduct.ImageURI)
			t.Fatal(expectedProduct.Name + " not parsed correctly")
		}
	}
}

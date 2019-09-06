package airborneoutlet

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beeronbeard/go-watch-that-site/product"
)

var expectedProducts = []product.Product{
	product.Product{Name: "Seeker 29 $ 729.95", URI: "/collections/outlet/products/seeker"},
	product.Product{Name: "Griffin 27.5+ (Demo Bike) $ 899.95", URI: "/collections/outlet/products/griffin-27-5-demo-bike"},
	product.Product{Name: "Goblin 29 (Demo Bike) $ 999.95", URI: "/collections/outlet/products/goblin-29-demo-bike"},
	product.Product{Name: "Goblin EVO 27.5 (Demo Bike) $ 1,199.95", URI: "/collections/outlet/products/goblin-evo-27-5-demo-bike"},
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

	a := AirborneOutlet{server.Client(), server.URL}
	productChannel := make(chan *product.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)
	go a.FindProducts(productChannel, errorChannel, completeChannel)

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

	if products[0] != expectedProducts[0] {
		t.Fatal("Seeker not parsed correctly")
	}

	if products[1] != expectedProducts[1] {
		t.Fatal("Griffin not parsed correctly")
	}

	if products[2] != expectedProducts[2] {
		t.Fatal("Goblin not parsed correctly")
	}

	if products[3] != expectedProducts[3] {
		t.Fatal("Goblin EVO not parsed correctly")
	}
}

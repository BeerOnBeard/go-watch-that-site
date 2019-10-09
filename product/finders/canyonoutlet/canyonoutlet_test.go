package canyonoutlet

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beeronbeard/go-watch-that-site/product"
)

var expectedProducts = []product.Product{
	product.Product{Name: "Exceed CF SL 6.9 Pro Race  Price: $2,159.99", URI: "https://www.canyon.com/en-us/outlet/mountain-bikes/exceed-cf-sl-6.9-pro-race/1321.html?dwvar_1321_pv_rahmenfarbe=BK&dwvar_1321_pv_rahmengroesse=XS", ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwf3517e28/images/full/full_2017_/2017/full_2017_exceed-cf-sl-69-pro-race_c1023.png"},
	product.Product{Name: "Spectral AL 6.0  Price: $2,470.00", URI: "/en-us/outlet/mountain-bikes/spectral-al-6.0/50008701_M05301B18H0750-3.html", ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwe64f7b78/images/full/full_spect/2019/full_spectral-al-6_c1328.png"},
	product.Product{Name: "Torque CF 7.0  Price: $3,500.00", URI: "https://www.canyon.com/en-us/outlet/mountain-bikes/torque-cf-7.0/2541.html?dwvar_2541_pv_rahmengroesse=XL&dwvar_2541_pv_rahmenfarbe=BK%2FBU", ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwdb056baa/images/full/full_us_20/2018/full_us_2018_torque-cf-7-us_c1277.png"},
}

func TestFindProducts(t *testing.T) {
	data, err := ioutil.ReadFile("canyonoutlet.html")
	if err != nil {
		t.Fatalf("Could not read file. %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(data)
	}))
	defer server.Close()

	productChannel := make(chan *product.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)

	finder := CanyonOutlet{server.Client(), server.URL}
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

	if products[0] != expectedProducts[0] {
		fmt.Println(products[0].Name)
		fmt.Println(expectedProducts[0].Name)
		fmt.Println(products[0].URI)
		fmt.Println(expectedProducts[0].URI)
		fmt.Println(products[0].ImageURI)
		fmt.Println(expectedProducts[0].ImageURI)
		t.Fatal("Exceed not parsed correctly")
	}

	if products[1] != expectedProducts[1] {
		t.Fatal("Spectral not parsed correctly")
	}

	if products[2] != expectedProducts[2] {
		t.Fatal("Torque not parsed correctly")
	}
}

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
	product.Product{
		Name:     "Exceed CF SL 6.9 Pro Race  Price: $2,159.99",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/exceed-cf-sl-6.9-pro-race/1321.html?dwvar_1321_pv_rahmenfarbe=BK&dwvar_1321_pv_rahmengroesse=XS",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwf3517e28/images/full/full_2017_/2017/full_2017_exceed-cf-sl-69-pro-race_c1023.png"},
	product.Product{
		Name:     "Spectral AL 6.0  Price: $2,530.00",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/spectral-al-6.0/50008357_M05301B18M0322-5.html",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwe64f7b78/images/full/full_spect/2019/full_spectral-al-6_c1328.png"},
	product.Product{
		Name:     "Exceed CF SL 6.0 Pro Race  Price: $1,980.00",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/exceed-cf-sl-6.0-pro-race/50007942_M10001H18H0435.html",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dw10fc4e23/images/full/full_excee/2019/full_exceed-cf-sl-6-pro-race_c1024.png"},
	product.Product{
		Name:     "Strive CFR 9.0 Team  Price: $5,220.00",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/strive-cfr-9.0-team/50012063_M06601V18J0192.html",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dw61ae4f8b/images/full/full_striv/2019/full_strive-cfr-9-team_c1321.png"},
	product.Product{
		Name:     "Spectral WMN AL 5.0  Price: $2,180.00",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/spectral-wmn-al-5.0/50008562_M05001B18H0187-3.html",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dwb9bd6e40/images/full/full_spect/2019/full_spectral-al-5-wmn_c1267.png"},
	product.Product{
		Name:     "Neuron AL 6.0 WMN  Price: $1,740.00",
		URI:      "https://www.canyon.com/en-us/outlet/mountain-bikes/neuron-al-6.0-wmn/50012564_M07001B18E0649.html",
		ImageURI: "https://www.canyon.com/dw/image/v2/BCML_PRD/on/demandware.static/-/Sites-canyon-master/default/dw7330b946/images/full/full_neuro/2019/full_neuron-al-6-wmn-us_c1315.png"},
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

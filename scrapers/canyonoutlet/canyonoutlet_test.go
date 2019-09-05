package canyonoutlet

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beeronbeard/go-watch-that-site/scrapers"
)

var expectedProducts = []scrapers.Product{
	scrapers.Product{Name: "Exceed CF SL 6.9 Pro Race  Price: $2,159.99", URI: "https://www.canyon.com/en-us/outlet/mountain-bikes/exceed-cf-sl-6.9-pro-race/1321.html?dwvar_1321_pv_rahmenfarbe=BK&dwvar_1321_pv_rahmengroesse=XS"},
	scrapers.Product{Name: "Spectral AL 6.0  Price: $2,470.00", URI: "/en-us/outlet/mountain-bikes/spectral-al-6.0/50008701_M05301B18H0750-3.html"},
	scrapers.Product{Name: "Torque CF 7.0  Price: $3,500.00", URI: "https://www.canyon.com/en-us/outlet/mountain-bikes/torque-cf-7.0/2541.html?dwvar_2541_pv_rahmengroesse=XL&dwvar_2541_pv_rahmenfarbe=BK%2FBU"},
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

	c := CanyonOutlet{server.Client(), server.URL}
	products, err := c.FindProducts()
	if err != nil {
		t.Fatalf("FindProducts failed. %v", err)
	}

	if products[0] != expectedProducts[0] {
		t.Fatal("Exceed not parsed correctly")
	}

	if products[1] != expectedProducts[1] {
		t.Fatal("Spectral not parsed correctly")
	}

	if products[2] != expectedProducts[2] {
		t.Fatal("Torque not parsed correctly")
	}
}

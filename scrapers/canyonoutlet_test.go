package scrapers

import (
	"bytes"
	"io/ioutil"
	"testing"

	"golang.org/x/net/html"
)

var expectedProducts = []Product{
	Product{"Exceed CF SL 6.9 Pro Race  Price: $2,159.99", "https://www.canyon.com/en-us/outlet/mountain-bikes/exceed-cf-sl-6.9-pro-race/1321.html?dwvar_1321_pv_rahmenfarbe=BK&dwvar_1321_pv_rahmengroesse=XS"},
	Product{"Spectral AL 6.0  Price: $2,470.00", "/en-us/outlet/mountain-bikes/spectral-al-6.0/50008701_M05301B18H0750-3.html"},
	Product{"Torque CF 7.0  Price: $3,500.00", "https://www.canyon.com/en-us/outlet/mountain-bikes/torque-cf-7.0/2541.html?dwvar_2541_pv_rahmengroesse=XL&dwvar_2541_pv_rahmenfarbe=BK%2FBU"},
}

func TestFindProducts(t *testing.T) {
	data, err := ioutil.ReadFile("canyonoutlet.html")
	if err != nil {
		t.Fatalf("Could not read file. %v", err)
	}

	reader := bytes.NewReader(data)
	doc, err := html.Parse(reader)
	if err != nil {
		t.Fatalf("Could not parse html. %v", err)
	}

	c := CanyonOutlet{}
	products, err := c.FindProducts(doc)
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

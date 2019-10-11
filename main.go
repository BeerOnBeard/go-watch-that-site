package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/beeronbeard/go-watch-that-site/product"
	"github.com/beeronbeard/go-watch-that-site/product/finders/airborneoutlet"
	"github.com/beeronbeard/go-watch-that-site/product/finders/canyonoutlet"
	"github.com/beeronbeard/go-watch-that-site/product/storer"
)

const (
	productsFilePath  = "bikes"
	canyonOutletURI   = "https://www.canyon.com/en-us/outlet/mountain-bikes/"
	airborneOutletURI = "https://airbornebicycles.com/collections/outlet"
	emailSubject      = "Outlet Bikes Update!"
)

func main() {
	envEmailTo := os.Getenv("GWTS_EMAIL_TO")
	envEmailFrom := os.Getenv("GWTS_EMAIL_FROM")
	envEmailPassword := os.Getenv("GWTS_EMAIL_PASSWORD")

	emailTo := flag.String("emailTo", envEmailTo, "Where to send email to")
	emailFrom := flag.String("emailFrom", envEmailFrom, "Where to send email from")
	emailPassword := flag.String("emailPassword", envEmailPassword, "Password for email auth")

	flag.Parse()

	if *emailTo == "" || *emailFrom == "" || *emailPassword == "" {
		usage()
		os.Exit(1)
	}

	finders := []product.Finder{
		&canyonoutlet.CanyonOutlet{Client: http.DefaultClient, URI: canyonOutletURI},
		&airborneoutlet.AirborneOutlet{Client: http.DefaultClient, URI: airborneOutletURI},
	}

	productChannel := make(chan *product.Product)
	errorChannel := make(chan *error)
	completeChannel := make(chan bool)

	for _, finder := range finders {
		go finder.Find(productChannel, errorChannel, completeChannel)
	}

	completeCount := 0

	var products []*product.Product
loop:
	for {
		select {
		case product := <-productChannel:
			products = append(products, product)
		case err := <-errorChannel:
			panic(err)
		case <-completeChannel:
			completeCount++
			if completeCount == len(finders) {
				break loop
			}
		}
	}

	storer := storer.New(productsFilePath)

	storedProducts, err := storer.Get()
	if err != nil {
		panic(err)
	}

	var newProducts []*product.Product
newProductsLoop:
	for i := 0; i < len(products); i++ {
		for j := 0; j < len(storedProducts); j++ {
			if products[i].Name == storedProducts[j].Name {
				continue newProductsLoop
			}
		}

		newProducts = append(newProducts, products[i])
	}

	var removedProducts []*product.Product
removedProductsLoop:
	for i := 0; i < len(storedProducts); i++ {
		for j := 0; j < len(products); j++ {
			if storedProducts[i].Name == products[j].Name {
				continue removedProductsLoop
			}
		}

		removedProducts = append(removedProducts, storedProducts[i])
	}

	if len(newProducts) > 0 || len(removedProducts) > 0 {
		mailer := GMailer{*emailFrom, *emailPassword}
		mailer.Send(*emailTo, emailSubject, generateEmailBody(newProducts, removedProducts))
	}

	err = storer.Put(products)
	if err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("go-watch-that-site -emailTo=email@email.test -emailFrom=email@email.test -emailPassword=superDuperSecrectH@x0rPa$$w0rd")
}

func generateEmailBody(newProducts, removedProducts []*product.Product) string {
	body := "<html><body><h1>Outlet Bikes Update</h1><h2>New Bikes</h2>"

	for _, product := range newProducts {
		body += generateProductHTML(product)
	}

	body += "<h2>Removed Bikes</h2>"

	for _, product := range removedProducts {
		body += generateProductHTML(product)
	}

	body += "</tbody></table></body></html>"
	return body
}

func generateProductHTML(p *product.Product) string {
	return fmt.Sprintf("<a href=\"%s\"><div><h3>%s</h3><img style=\"width:100%%; max-width:800px\" src=\"%s\" /></div></a>", p.URI, p.Name, p.ImageURI)
}

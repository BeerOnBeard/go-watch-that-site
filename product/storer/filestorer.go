package storer

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/beeronbeard/go-watch-that-site/product"
)

type fileStorer struct {
	FilePath string
}

// New creates a new instance of a product storer backed by a file at path
func New(path string) product.Storer {
	return &fileStorer{FilePath: path}
}

// Get retrieves the products stored
func (storer *fileStorer) Get() ([]*product.Product, error) {
	_, err := os.Stat(storer.FilePath)
	if os.IsNotExist(err) {
		return []*product.Product{}, nil
	}

	content, err := ioutil.ReadFile(storer.FilePath)
	if err != nil {
		return nil, err
	}

	var products []*product.Product
	err = json.Unmarshal(content, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// Put replaces the stored products
func (storer *fileStorer) Put(products []*product.Product) error {
	file, err := os.OpenFile(storer.FilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	p, err := json.Marshal(products)
	if err != nil {
		return err
	}

	_, err = file.Write(p)
	return err
}

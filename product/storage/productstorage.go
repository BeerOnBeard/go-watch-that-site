package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/beeronbeard/go-watch-that-site/product"
)

type fileStorage struct {
	FilePath string
}

// NewFile creates a new instance of ProductStorage backed by a file at path
func NewFile(path string) product.Storage {
	return &fileStorage{FilePath: path}
}

// Get retrieves the products stored
func (s *fileStorage) Get() ([]*product.Product, error) {
	_, err := os.Stat(s.FilePath)
	if os.IsNotExist(err) {
		return []*product.Product{}, nil
	}

	content, err := ioutil.ReadFile(s.FilePath)
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
func (s *fileStorage) Put(products []*product.Product) error {
	file, err := os.OpenFile(s.FilePath, os.O_WRONLY|os.O_CREATE, 0600)
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

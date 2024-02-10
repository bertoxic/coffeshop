package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetProductList() Products {
	return productList
}
func AddPrduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)

}

func getNextID() int {
	lp := productList[len(productList)-1]

	return lp.ID + 1

}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(rd io.Reader) error {
	e := json.NewDecoder(rd)
	return e.Decode(p)
}
func UpdateProduct(id int, product *Product) error{
	_ , pos, err := findProduct(id)
	if err != nil {
		return err
	}
	product.ID = id
	productList[pos] = product;
		return nil
}

var ErroProductNotFound = fmt.Errorf("product was not found")

func findProduct(id int) (*Product, int, error) {

	for i, p := range productList {

		if id == p.ID {
			return p,i ,nil
		}

	}
	return nil,0, ErroProductNotFound
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "latte",
		Description: "Frothy cold not tasty black",
		Price:       2.38,
		SKU:         "5467hj",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	&Product{
		ID:          2,
		Name:        "edonta",
		Description: "Frothy cold not tasty black",
		Price:       12.38,
		SKU:         "t647hj",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}

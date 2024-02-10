
// schemes: http
// Host: localhost
// BasePath: /v2
// Version: 0.0.1

// consumes:
// -applicatuin/json
// -application/json



package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/bertoxic/coffeshop/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		if p == nil {
			log.Println("Error: Products receiver is nil")
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		p.AddProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		//expect the id in the uri
		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		//path := regex.FindStringSubmatch(r.URL.Path)
		if len(g) != 1 {
			http.Error(rw, "Invalid uri..x", http.StatusBadRequest)
			return

		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid uri..xx", http.StatusBadRequest)
			return

		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid uri..xxx", http.StatusBadRequest)

			return
		}
		p.l.Println("got id", id)

		p.UpdateProduct(rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProductList()

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to print jsonData", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle post service addproduct handler")
	vars := mux.Vars(r)
	p.l.Println("prices isssss", vars["price"])
	p.l.Println(r.Body)
	p.l.Println("name isssss", vars["name"])
	//price, err := strconv.Atoi(vars["price"])
	// if err != nil {
	// 	http.Error(rw, "cannot convert string to int in price", http.StatusBadRequest)
	// 	return
	// }
	// prod := &data.Product{
	// 	Name:        vars["name"],
	// 	Price:       float32(price),
	// 	Description: vars["description"],
	// 	SKU:         vars["sku"],
	// }
	//prod:= &data.Product{}
	prod:= r.Context().Value(Keyproduct{}).(*data.Product)
	err := prod.FromJSON(r.Body)
	p.l.Println(prod)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to unmarshall json %v",err), http.StatusBadRequest)
	}

	p.l.Printf("Prod:%v", prod)
	data.AddPrduct(prod)

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	p.l.Println("name isssss", vars["name"])
	if err != nil {
		http.Error(rw, "unable to unmarshal id for PUT product", http.StatusBadRequest)
		return
	}
	//prod := &data.Product{}

	//err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json for PUT product", http.StatusBadRequest)
		return
	}
	p.l.Println("Hanle PUT product", id)
	prod:= r.Context().Value(Keyproduct{}).(*data.Product)
	err = data.UpdateProduct(id, prod)
	p.l.Println("product has just been recently updated")
	if err == data.ErroProductNotFound {
		http.Error(rw, "product not found", http.StatusNotFound)
	}

}

func (p *Products) getm(rw http.ResponseWriter, r *http.Request) {

}
 type Keyproduct struct {}
func (p *Products) MiddleWareProductValidation(next http.Handler) http.Handler{

		return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request){
			prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json for PUT product", http.StatusBadRequest)
		return
	}

	ctx:= context.WithValue(r.Context(),Keyproduct{},prod)
	r = r.WithContext(ctx)
	next.ServeHTTP(rw, r)
		})
		
}



//curl  localhost:9090/ -XPOST -d '{"name":"buyete","description":"grey hot  not tasty black","price":2.38,"sku":"5467hj"}'
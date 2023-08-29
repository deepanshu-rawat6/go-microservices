package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/deepanshu-rawat6/go-microservices/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns the products from the data store
// By changing just the name from getProducts to GetProducts, our function becomes public
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to encode into json", http.StatusInternalServerError)
	}
}

// AddProducts add a product to the data store
func (p *Products) AddProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadGateway)
	}

	p.l.Printf("Prod: %#v", prod)

	data.AddProduct(&prod)
}

// UpdateProducts puts/updates a specific product present in the data store
func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product", id)

	prod := &data.Product{}

	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshall json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/deepanshu-rawat6/go-microservices/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.AddProducts(w, r)
		return
	}

	if r.Method == http.MethodPut {
		// expecting the id in the URI
		rgx := regexp.MustCompile(`/([0-9]+)`)
		g := rgx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		p.l.Println("Got ID", id)
	}

	// catch all(request methods other than GET, will get 405 error)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to encode into json", http.StatusInternalServerError)
	}
}

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

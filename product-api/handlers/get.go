package handlers

import (
	"net/http"

	"github.com/deepanshu-rawat6/go-microservices/product-api/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// GetProducts returns the products from the data store
// By changing just the name from getProducts to GetProducts, our function becomes public
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

	prods := data.GetProducts()
	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}

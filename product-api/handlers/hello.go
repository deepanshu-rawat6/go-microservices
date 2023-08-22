package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello -> simple handler
type Hello struct {
	l *log.Logger
}

// Creates a new hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHTTP implements the go http.Hanlder interface
//
//	https://pkg.go.dev/net/http#Handler
func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello, World!")

	// read the body
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Println("Error reading body", err)

		http.Error(w, "Oops", http.StatusBadRequest)
		return
	}

	// Writing response
	fmt.Fprintf(w, "Hello, %s", d)
}

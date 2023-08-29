package main

import (
	"context" // https://pkg.go.dev/context
	"log"     // https://pkg.go.dev/log
	"net/http"
	"os"        // https://pkg.go.dev/os
	"os/signal" // https://pkg.go.dev/os/signal
	"time"      // https://pkg.go.dev/time

	"github.com/deepanshu-rawat6/go-microservices/product-api/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	env.Parse()

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Handlers similar to routes to trigger funcs when hitting that end point
	ph := handlers.NewProducts(l)

	// Router from gorilla mux package
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProducts)

	// sm.Handle("/products", ph)

	// custom server, avoid default server
	s := &http.Server{
		// Addr:         ":" + portString,
		Addr:         *bindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// using gorutines, another thread executes the server logic
	go func() {
		l.Printf("Starting server on port %v", portString)

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// For graceful termination, we use os.signal, to send signals like interrupt(if we are closing the server)
	// Now, if there are no requests in queue, then we send the kill signal to stop the server gracefully
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received.
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig) // graceful termination

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	tcx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tcx)
}

package main

import (
	"context" // https://pkg.go.dev/context
	"log"     // https://pkg.go.dev/log
	"net/http"
	"os"        // https://pkg.go.dev/os
	"os/signal" // https://pkg.go.dev/os/signal
	"time"      // https://pkg.go.dev/time

	"github.com/deepanshu-rawat6/go-microservices/product-api/data"
	"github.com/deepanshu-rawat6/go-microservices/product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicholasjackson/env"

	gohandlers "github.com/gorilla/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	env.Parse()

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	// Handlers similar to routes to trigger funcs when hitting that end point
	ph := handlers.NewProducts(l, v)

	// Router from gorilla mux package
	sm := mux.NewRouter()

	// Hanlders for API
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
	}

	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS Handlers
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// sm.Handle("/products", ph)

	// custom server, avoid default server
	s := &http.Server{
		// Addr:         ":" + portString,
		Addr:         *bindAddress,
		Handler:      ch(sm),
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig) // graceful termination

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

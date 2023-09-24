# Microservices in Go

## Part 7: Documenting RESTful APIs with OpenAPI and Swagger

## Part 8: Autogenerating HTTP clients from swagger files

### Generating HTTP Clients

Using `swagger` to generate HTTP clients, using the command: 

```bash
swagger generate client <path to swagger.yaml> -A <name>
```

#### Good Pracitce:

Try to structure auto-generated snippets, by `go-swagger` to a folder like `sdk`, or maybe in other repository.

### Creating tests for clients

```go
func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod.GetPayload()[0])
}
```

#### Debugging: 

Adding the correct `content-type` headers for our `swagger spec`. Debugging auto-generated code can be really painful.

```go
rw.Header().Add("Content-Type","application/json")
```

## Part 9: Handling CORS(Cross-Origin Resoruce Sharing)

First understanding CORS, we can use this blog post on medium [here](https://medium.com/@baphemot/understanding-cors-18ad6b478e2b)

Importing `gorilla/handlers` and using the **CORS** method, to **grant** `allowed origins` for `localhost:3000`

Here:
```go
ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
```

Wrapping the `sm` router from the mux package from gorilla into the **CORS** handlers `sm`

## Part 10: Handling files using the Go Standard Library

More insights on [http](https://pkg.go.dev/net/http) package
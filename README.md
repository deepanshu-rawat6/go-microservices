# Microservices in Go

## Part 8: Autogenerating HTTP clients from swagger files

### Generating HTTP Clients

Using `swagger` to generate HTTP clients, using the command: 

```bash
swagger generate client <path to swagger.yaml> -A <name>
```

#### Good Pracitce:

Try to structure auto-generated snippets, by `go-swagger` to a folder like `sdk`, or maybe in other repository.

### Creating tests for clients




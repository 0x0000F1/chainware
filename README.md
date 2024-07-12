# Middleware Chaining Library for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/0x0000F1/chainware.svg)](https://pkg.go.dev/github.com/0x0000F1/chainware)

This library provides a way to manage and chain HTTP middlewares using a hash table with separate chaining for Go's `net/http` package. It allows you to easily add, retrieve, and chain middlewares for different routes in a thread-safe manner.

## Features

- Efficient and thread-safe middleware management.
- Supports chaining of multiple middlewares for different routes.
- Easy to integrate with existing Go `net/http` applications.

## Installation

To install the library, you can use `go get`:

```sh
go get github.com/0x0000F1/chainware
```

## Usage
Importing the Package

```go

import (
    "github.com/0x0000F1/chainware"
    "net/http"
)
```

## Example

Here's a complete example demonstrating how to use the middleware chaining library:

```go

package main

import (
    "fmt"
    "github.com/0x0000F1/chainware"
    "net/http"
)

func main() {
    // Create a new middleware table
    mt := middleware.NewMiddlewareTable()

    // Define some example middlewares
    loggingMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Println("Logging middleware")
            next.ServeHTTP(w, r)
        })
    }

    authMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Println("Auth middleware")
            next.ServeHTTP(w, r)
        })
    }

    // Add middlewares to the table for a specific route
    mt.AddMiddleware("/example", loggingMiddleware)
    mt.AddMiddleware("/example", authMiddleware)

    // Define the final handler
    finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })

    // Create the server
    http.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
        chain := mt.GetChain("/example")
        if chain != nil {
            handler := chain.ChainMiddlewares(finalHandler)
            handler.ServeHTTP(w, r)
        } else {
            finalHandler.ServeHTTP(w, r)
        }
    })

    http.ListenAndServe(":8080", nil)
}
```

## Explanation

    1. **Middleware Type**: A Middleware is a function that wraps an http.Handler and returns a new http.Handler.
    2. **Chain**: A Chain represents a sequence of middlewares. It supports adding middlewares and chaining them to create a final http.Handler.
    3. **MiddlewareTable**: A MiddlewareTable manages multiple Chain instances, each associated with a different key (e.g., route path).
    4. **Adding Middlewares**: Use AddMiddleware to add a middleware to a chain associated with a specific key.
    5. **Chaining Middlewares**: Use ChainMiddlewares to wrap the final http.Handler with the middlewares in the chain.

## Documentation
**`pkg/chain.go`**

   - **Middleware**: A type definition for middleware functions.
   - **Chain**: A struct representing a chain of middlewares.
       - **NewChain(capacity int) *Chain**: Creates a new Chain with a given initial capacity.
       - **AddMiddleware(m Middleware)**: Adds a middleware to the chain.
       - **GetMiddlewares() []Middleware**: Returns the list of middlewares in the chain.
       - **ChainMiddlewares(finalHandler http.Handler) http.Handler**: Chains the middlewares and returns the final http.Handler.

   - **MiddlewareTable**: A struct that holds multiple chains of middlewares.
       - **NewMiddlewareTable() *MiddlewareTable**: Creates a new MiddlewareTable.
       - **AddMiddleware(key string, m Middleware)**: Adds a middleware to the chain associated with the given key.
       - **GetChain(key string) *Chain**: Returns the chain of middlewares associated with the given key.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue on GitHub.

## License

This project is licensed under the MIT License.

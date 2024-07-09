package main

import (
    "github.com/0x0000F1/chainware/pkg"

	"fmt"
	"net/http"
)

// Sample middleware that logs the request path.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request path: %s\n", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Sample middleware that adds a header.
func addHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Custom-Header", "MiddlewareHeader")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Final handler that handles the request.
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create a new middleware chain with an initial capacity and sample middlewares.
	chainware := chain.NewChain(2, loggingMiddleware, addHeaderMiddleware)

	// Use the chain with the final handler.
	http.Handle("/", chainware.Handler(finalHandler))

	// Start the server.
	http.ListenAndServe(":8080", nil)
}


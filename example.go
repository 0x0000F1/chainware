package main

import (
    "fmt"
    middleware "github.com/0x0000F1/chainware/pkg"
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


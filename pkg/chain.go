package chain

import (
	"net/http"
)

// Middleware defines a function to process middleware.
type Middleware func(http.Handler) http.Handler

// Chain represents a stack of middlewares.
type Chain struct {
	middlewares []Middleware
}

// NewChain initializes a new middleware chain with a given capacity.
func NewChain(initialCapacity int, middlewares ...Middleware) *Chain {
	chain := &Chain{
		middlewares: make([]Middleware, 0, initialCapacity),
	}
	chain.middlewares = append(chain.middlewares, middlewares...)
	return chain
}

// Add appends a new middleware to the chain.
func (c *Chain) Add(middleware Middleware) *Chain {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Handler applies the middleware chain to a final handler.
func (c *Chain) Handler(finalHandler http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		finalHandler = c.middlewares[i](finalHandler)
	}
	return finalHandler
}

// HandlerFunc applies the middleware chain to a final handler function.
func (c *Chain) HandlerFunc(finalHandlerFunc http.HandlerFunc) http.Handler {
	return c.Handler(finalHandlerFunc)
}

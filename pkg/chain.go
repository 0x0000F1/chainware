package chain

import (
    "net/http"
    "sync"
)

// Middleware is a function that wraps an http.Handler.
// It takes an http.Handler as input and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain represents a chain of middlewares.
// It holds a slice of Middleware functions and a mutex for concurrent access.
type Chain struct {
    middlewares []Middleware
    mutex       sync.RWMutex
}

// MiddlewareTable holds the chains of middlewares for different keys.
// It uses a map to associate keys (e.g., route paths) with their corresponding middleware chains.
type MiddlewareTable struct {
    table map[string]*Chain
    mutex sync.RWMutex
}

// NewChain creates and returns a new Chain.
// It pre-allocates a slice with the given capacity to optimize memory usage.
func NewChain(capacity int) *Chain {
    return &Chain{
        middlewares: make([]Middleware, 0, capacity),
    }
}

// AddMiddleware adds a middleware to the chain.
// It locks the Chain to ensure thread safety during the append operation.
func (c *Chain) AddMiddleware(m Middleware) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.middlewares = append(c.middlewares, m)
}

// GetMiddlewares returns the list of middlewares in the chain.
// It uses a read lock to allow concurrent read access.
func (c *Chain) GetMiddlewares() []Middleware {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.middlewares
}

// ChainMiddlewares chains the middlewares and returns a final http.Handler.
// It starts with the finalHandler and wraps it with each middleware in reverse order.
func (c *Chain) ChainMiddlewares(finalHandler http.Handler) http.Handler {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    handler := finalHandler
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        handler = c.middlewares[i](handler)
    }
    return handler
}

// NewMiddlewareTable creates and returns a new MiddlewareTable.
// It initializes the map to hold the chains of middlewares.
func NewMiddlewareTable() *MiddlewareTable {
    return &MiddlewareTable{
        table: make(map[string]*Chain),
    }
}

// AddMiddleware adds a middleware to the chain for a given key.
// It locks the table to ensure thread safety during the update operation.
func (mt *MiddlewareTable) AddMiddleware(key string, m Middleware) {
    mt.mutex.Lock()
    defer mt.mutex.Unlock()
    // If the chain for the given key doesn't exist, create a new one with initial capacity of 1.
    if _, exists := mt.table[key]; !exists {
        mt.table[key] = NewChain(1) // Start with capacity of 1
    }
    // Add the middleware to the chain.
    mt.table[key].AddMiddleware(m)
}


// GetChain returns the chain of middlewares for a given key.
// It uses a read lock to allow concurrent read access.
func (mt *MiddlewareTable) GetChain(key string) *Chain {
    mt.mutex.RLock()
    defer mt.mutex.RUnlock()
    return mt.table[key]
}


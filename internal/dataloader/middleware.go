package dataloader

import (
	"context"
	"net/http"

	"github.com/example/dpc-dataloader-demo/internal/core/ports"
)

// contextKey is an unexported type to prevent context key collisions.
// Using a struct type (instead of a bare string) ensures no other package
// can accidentally produce the same key.
type contextKey struct{}

// loadersKey is the context key for storing/retrieving Loaders.
var loadersKey = contextKey{}

// Middleware returns an HTTP middleware that creates a fresh Loaders instance
// per request and stores it in the request context.
//
// Usage in server setup:
//
//	router.Use(dataloader.Middleware(redemptionRepo))
//
// This MUST be applied before the GraphQL handler so that field resolvers
// can retrieve the loaders via dataloader.For(ctx).
func Middleware(redemptionRepo ports.RedemptionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(redemptionRepo)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For retrieves the Loaders from the given context.
// Returns nil if no Loaders are found (e.g., if the middleware was not applied).
// Field resolvers should call this to get the loader instance.
func For(ctx context.Context) *Loaders {
	loaders, _ := ctx.Value(loadersKey).(*Loaders)
	return loaders
}

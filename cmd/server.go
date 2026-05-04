package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"

	"github.com/example/dpc-dataloader-demo/graph"
	"github.com/example/dpc-dataloader-demo/internal/core/service"
	"github.com/example/dpc-dataloader-demo/internal/dataloader"
	"github.com/example/dpc-dataloader-demo/internal/repositories"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ─── Dependency Injection ──────────────────────────────────────
	// Same pattern as the original project's cmd/server.go

	// 1. Repository layer
	redemptionRepo := repositories.NewRedemptionPort()

	// 2. Service layer
	redemptionSvc := service.NewRedemptionSvc(redemptionRepo)

	// 3. Router setup (chi, same as original project)
	r := chi.NewRouter()
	r.Use(cors.Default().Handler)

	// 4. DataLoader middleware — MUST be applied BEFORE the GraphQL handler.
	//    This creates a fresh Loaders instance per request and stores it in context.
	//    The repository is passed in so the loaders can call batch-fetch methods.
	r.Use(dataloader.Middleware(redemptionRepo))

	// 5. GraphQL server (gqlgen, same as original project)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			RedemptionSvc: redemptionSvc,
		},
	}))

	// 6. Routes
	r.Handle("/query", srv)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

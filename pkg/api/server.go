package main

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/kacejot/resize-service/pkg/storage/db"
	"github.com/kacejot/resize-service/pkg/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kacejot/resize-service/pkg/api/graph"
	"github.com/kacejot/resize-service/pkg/api/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver, err := graph.NewResolver()
	utils.Unwrap(err)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("Authorization")

		ctx := context.WithValue(r.Context(), db.UserContextKey, user)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

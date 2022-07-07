package main

import (
	"net/http"
	"os"

	"CP_Discussion/graph"
	"CP_Discussion/graph/generated"
	"CP_Discussion/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", c.Handler(srv))

	log.Info.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Error.Fatal(http.ListenAndServe(":"+port, nil))
}

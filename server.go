package main

import (
	"net/http"
	"os"

	"CP_Discussion/directive"
	"CP_Discussion/graph"
	"CP_Discussion/graph/generated"
	"CP_Discussion/log"
	"CP_Discussion/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
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

	conf := generated.Config{Resolvers: &graph.Resolver{}}
	conf.Directives.Auth = directive.AuthDirective
	conf.Directives.Admin = directive.AdminDirective

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(conf))

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middleware.AuthMiddleware)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(srv))
	router.Methods(http.MethodGet, http.MethodPost)

	log.Info.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Error.Fatal(http.ListenAndServe(":"+port, router))
}

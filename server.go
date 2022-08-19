package main

import (
	"os"

	"CP_Discussion/auth"
	"CP_Discussion/directive"
	"CP_Discussion/file/fileHandler"
	"CP_Discussion/graph"
	"CP_Discussion/graph/generated"
	"CP_Discussion/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	conf := generated.Config{Resolvers: &graph.Resolver{}}
	conf.Directives.Auth = directive.AuthDirective
	conf.Directives.Admin = directive.AdminDirective

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(conf))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{"localhost"})
	router.Use(auth.CorsMiddleware())
	router.GET("/", auth.AuthMiddleware(), playgroundHandler())
	router.POST("/query", auth.AuthMiddleware(), graphqlHandler())
	router.GET("/post/:year/:semester/:filename", fileHandler.FileHandler())
	router.GET("/avatar/:filename", fileHandler.FileHandler())
	router.POST("/refresh", auth.RefreshHandler())

	log.Info.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Error.Fatal(router.Run(":" + port))
}

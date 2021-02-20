package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/inblack67/GQLGenAPI/cache"
	"github.com/inblack67/GQLGenAPI/db"
	"github.com/inblack67/GQLGenAPI/graph"
	"github.com/inblack67/GQLGenAPI/graph/generated"
	"github.com/inblack67/GQLGenAPI/middlewares"
	"github.com/inblack67/GQLGenAPI/mysession"
	"github.com/rs/cors"
)

const defaultPort = "5000"

func main() {

	cache.StartRedis()
	db.ConnectDB()
	mysession.InitSessionStore()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()

	router.Use(
		cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowCredentials: true,
			// Debug: true,
		}).Handler,
	)

	router.Use(middlewares.AuthMiddleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphdata/graphstore"

	"github.com/99designs/gqlgen/handler"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphdata"
)

const defaultPort = "11111"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := &graphdata.Resolver{
		Store: graphstore.NewStore(),
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphdata.NewExecutableSchema(graphdata.Config{Resolvers: r})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

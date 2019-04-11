package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/grpc/interface/order"
	"google.golang.org/grpc"

	"github.com/99designs/gqlgen/handler"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphql/schema"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/client"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create the gRPC client that provides order details.
	grpcConn, err := grpc.Dial("orderservergrpc:8888", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	orderClient := order.NewOrdersClient(grpcConn)

	// Create the REST client for the company service.
	companyClient := client.NewCompanyClient("http://companyrestserver:9021")

	// Create a resolver.
	r := schema.NewResolver(orderClient, companyClient)

	// Allow mutations and resolvers to collect data from the services.
	server := handler.GraphQL(schema.NewExecutableSchema(schema.Config{
		Resolvers: r,
	}))
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	// Without data loader.
	// http.Handle("/query", server)
	// With data loader.
	// Configure the data loader to carry out additional lookups.
	http.Handle("/query", schema.WithCompanyDataloaderMiddleware(companyClient, server))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

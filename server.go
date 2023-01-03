package main

import (
	"CareXR_API/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"CareXR_API/config"
	"CareXR_API/ioutils"
)

const defaultPort = "8000"

func main() {
	settings, err := config.ReadConfig("config.json")
	ioutils.PanicOnError(err)
	println("Username ", settings.Username, "password", settings.Password)
	driver, err := config.NewDriver(settings)
	defer driver.Close()

	if driver == nil {
		os.Exit(1)
	}

	config.Neo4jDriver = driver

	// Test

	////

	defer func() {
		ioutils.PanicOnError(driver.Close())
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/api"))
	http.Handle("/api", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

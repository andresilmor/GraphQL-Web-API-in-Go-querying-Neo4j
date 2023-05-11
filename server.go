package main

import (
	"CareXR_WebService/graph"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"CareXR_WebService/config"
	"CareXR_WebService/ioutils"

	"github.com/gin-contrib/cors"
)

const defaultPort = "8000"

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func executeQuery(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)

	request, err := http.NewRequest("POST", "http://localhost:8000/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	responseRaw, err := io.ReadAll(res.Body)

	c.JSON(200, string(responseRaw))

}

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

	router := gin.Default()

	router.Use(GinContextToContextMiddleware())

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))

	router.POST("/api", graphqlHandler())

	router.GET("/view", playgroundHandler())

	router.Run(":8000")

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

}

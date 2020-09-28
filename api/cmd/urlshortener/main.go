package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"urlshortener/pkg/handlers"
	"urlshortener/pkg/mappings"
	"urlshortener/pkg/middleware"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Logger
	logger := log.New(os.Stderr, "LOG ", log.Lshortfile)

	// Web
	templates := template.Must(template.ParseFiles("./web/index.html"))

	// MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		logger.Fatalf("Can't connect to mongodb. %s", err.Error())
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatalf("MongoDB ping error. %s", err.Error())
		return
	}
	defer func(c *mongo.Client) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.Disconnect(ctx)
		if err != nil {
			logger.Fatalf("Can't disconnect from mongodb. %s", err.Error())
		}
	}(client)
	mappingsCollection := client.Database("urlshortener").Collection("mappings")

	// Repo
	urlsRepo := mappings.NewRepo(mappingsCollection)

	// Handlers
	mappingHandler := handlers.MappingHandler{
		MappingRepo: urlsRepo,
		Logger:      logger,
	}

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Execute(w, nil)
	})
	r.HandleFunc("/{short_url}", mappingHandler.Redirect).Methods("GET")
	r.HandleFunc("/create", mappingHandler.Add).Methods("POST")
	r.HandleFunc("/info/{short_url}", mappingHandler.GetMappingInfo).Methods("GET")

	handler := middleware.Log(logger, r)
	handler = middleware.Panic(handler)

	port := os.Getenv("API_PORT")
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

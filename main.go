package main

import (
	"elastic/controller"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := elasticsearch.Config{
		//CloudID:   os.Getenv("CLOUD_ID"),
		//APIKey:    os.Getenv("API_KEY"),
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
		Username:  os.Getenv("ELASTICSEARCH_USERNAME"),
		Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/search", controller.Search(es))
		r.Post("/create", controller.CreateOrUpdateIndexAndInsertData(es))
	})

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"strings"
	"time"
)

// ElasticsearchResponse represents the structure of the Elasticsearch response
type ElasticsearchResponse struct {
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}

// Shards represents the shards information in the Elasticsearch response
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// Hits represents the hits information in the Elasticsearch response
type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

// Total represents the total hits information in the Elasticsearch response
type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// Hit represents a single hit in the Elasticsearch response
type Hit struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Source  `json:"_source"`
}

// Source represents the source document in the Elasticsearch response
type Source struct {
	Date        time.Time `json:"Date"`
	Image       string    `json:"Image"`
	IsAvailable bool      `json:"IsAvailable"`
	Price       int       `json:"Price"`
	Stock       int       `json:"Stock"`
	Variant     string    `json:"Variant"`
	Content     string    `json:"content"`
	Title       string    `json:"title"`
}

func main() {
	cfg := elasticsearch.Config{
		CloudID: "********:************",
		APIKey:  "************",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// API Key should have cluster monitoring rights
	infores, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Println(infores)

	// Define a handler function to handle HTTP GET requests
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		//********** Search **********
		// Define your search query
		var (
			index  = []string{"sample-index"} // Replace "your_index" with your actual index name
			search = `{
			"query": {
				"match": {
					"content": "sample"
				}
			}
		}` // Replace "field_name" and "search_keywords" with your actual field name and search keywords
		)

		// Perform the search
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(index...),
			es.Search.WithBody(strings.NewReader(search)),
			es.Search.WithTrackTotalHits(true),
		)
		if err != nil {
			log.Fatalf("Error performing search: %s", err)
		}
		defer res.Body.Close()

		// Parse the response
		if res.IsError() {
			log.Fatalf("Error response: %s", res.String())
		}

		// Decode the response body into a struct
		var response ElasticsearchResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Fatalf("Error decoding response body: %s", err)
		}

		// Log the documents
		for _, hit := range response.Hits.Hits {
			log.Printf("Document ID: %s, Title: %s, Description: %s\n", hit.ID, hit.Source.Title, hit.Source.Content)
		}

		// Print the search results
		fmt.Println(res.String())

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Fatalf("Error encoding JSON response: %s", err)
		}
	})

	// Start the HTTP server on port 8080
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

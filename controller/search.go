package controller

import (
	"context"
	"elastic/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"strings"
)

// search function to perform the search operation in Elasticsearch

func Search(es *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Define your search query
		var (
			index  = []string{"sample-index"}                                                                                                          // Replace "your_index" with your actual index name
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
		var response model.ElasticsearchResponse
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
	}
}

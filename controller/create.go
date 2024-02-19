package controller

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"strings"
)

// Create function to perform the create operation in Elasticsearch

func Create(es *elasticsearch.Client, w http.ResponseWriter, r *http.Request) {
	// Define your create query
	var (
		index  = "sample-index" // Replace "your_index" with your actual index name
		id     = "1"            // Replace "your_id" with your actual document ID
		create = `{
			"content": "sample document"
		}` // Replace "field_name" and "field_value" with your actual field name and field value
	)

	// Perform the Create
	res, err := es.Create(
		index,                     // Index name
		id,                        // Document ID
		strings.NewReader(create), // Document body
		es.Create.WithContext(context.Background()),
		es.Create.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error performing create: %s", err)
	}
	defer res.Body.Close()

	// Parse the response
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	// Print the create results
	fmt.Println(res.String())

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Document created successfully"}`))
}

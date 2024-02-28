package controller

import (
	"context"
	"elastic/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// search function to perform the search operation in Elasticsearch

func Search(es *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var searchReq SearchRequest
		err := json.NewDecoder(req.Body).Decode(&searchReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusInternalServerError)
			return
		}

		// Pagination parameters
		page := 1
		size := 10 // Number of results per page
		if pageStr := req.URL.Query().Get("page"); pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		}
		if sizeStr := req.URL.Query().Get("size"); sizeStr != "" {
			size, _ = strconv.Atoi(sizeStr)
		}
		from := (page - 1) * size

		// Define your priority search fields
		priorityFields := []string{
			"productAttributeValue",
			"categoryNameEn",
			"categoryNameBn",
			"productNameEn",
			"productNameBn",
		}

		// Construct the "should" array based on priority fields
		var shouldClauses []string
		for _, field := range priorityFields {
			shouldClauses = append(shouldClauses, `{"match": {`+
				`"`+field+`": "`+searchReq.Keyword+`"}}`)
		}

		// Combine should clauses into the search query
		var search = `{
			"from": ` + strconv.Itoa(from) + `,
			"size": ` + strconv.Itoa(size) + `,
			"query": {
				"bool": {
					"should": [` + strings.Join(shouldClauses, ",") + `]
				}
			}
		}`

		// Perform the search
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("sl_test_index"), // Replace "your_index" with your actual index name
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
		var response model.Response
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Fatalf("Error decoding response body: %s", err)
		}

		// Log the documents
		for _, hit := range response.Hits.Hits {
			log.Printf("Document ID: %s, Document Index: %s", hit.ID, hit.Index)
		}

		// Print the search results
		fmt.Println(res.String())

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Fatalf("Error encoding JSON response: %s", err)
		}
	}
}

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

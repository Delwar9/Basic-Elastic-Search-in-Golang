package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
	"strings"
)

// CreateOrUpdateIndexAndInsertData checks if the index already exists. If it does, it inserts the data directly. If not, it creates the index and then inserts the data.

func CreateOrUpdateIndexAndInsertData(es *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Decode the request body into the Document struct
		var doc Products
		err := json.NewDecoder(req.Body).Decode(&doc)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusInternalServerError)
			return
		}

		// Define the index name
		indexName := "sl_test_index"

		// Check if index already exists
		res, err := es.Indices.Exists([]string{indexName})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error checking if index exists: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			log.Printf("Index '%s' already exists", indexName)
		} else {
			// Create the index if it doesn't exist
			indexSettings := `{
				"settings": {
					"number_of_shards": 1,
					"number_of_replicas": 0
				},
				"mappings": {
					"properties": {
						"productId": {"type": "integer"},
						"variantId": {"type": "integer"},
						"productAttributeId": {"type": "integer"},
						"productAttributeValue": {"type": "text"},
						"categoryId": {"type": "integer"},
						"categoryNameEn": {"type": "text"},
						"categoryNameBn": {"type": "text"},
						"productNameEn": {"type": "text"},
						"productNameBn": {"type": "text"},
						"price": {"type": "float"},
						"stock": {"type": "integer"},
						"discountAmt": {"type": "float"}
					}
				}
			}`

			createReq := esapi.IndicesCreateRequest{
				Index: indexName,
				Body:  strings.NewReader(indexSettings),
			}

			createRes, err := createReq.Do(context.Background(), es)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating index: %s", err), http.StatusInternalServerError)
				return
			}
			defer createRes.Body.Close()

			if createRes.IsError() {
				http.Error(w, fmt.Sprintf("Error creating index: %s", createRes.String()), createRes.StatusCode)
				return
			}

			log.Println("Index created successfully")
		}

		// Index data into the index
		docJSON, err := json.Marshal(doc)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling document: %s", err), http.StatusInternalServerError)
			return
		}

		indexReq := esapi.IndexRequest{
			Index:   indexName,
			Body:    strings.NewReader(string(docJSON)),
			Refresh: "true",
		}

		_, err = indexReq.Do(context.Background(), es)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error indexing document: %s", err), http.StatusInternalServerError)
			return
		}

		log.Println("Document indexed successfully")

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Document indexed successfully"})
	}
}

type Document struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type Products struct {
	ProductId             int64   `json:"productId"`
	VariantId             int64   `json:"variantId"`
	ProductAttributeId    int64   `json:"productAttributeId"`
	ProductAttributeValue string  `json:"productAttributeValue"`
	CategoryId            int64   `json:"categoryId"`
	CategoryNameEn        string  `json:"categoryNameEn"`
	CategoryNameBn        string  `json:"categoryNameBn"`
	ProductNameEn         string  `json:"productNameEn"`
	ProductNameBn         string  `json:"productNameBn"`
	Price                 float64 `json:"price"`
	Stock                 int     `json:"stock"`
	DiscountAmt           float64 `json:"discountAmt"`
}

package model

import "time"

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

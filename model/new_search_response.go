package model

type Hits struct {
	Total struct {
		Relation string `json:"relation"`
		Value    int    `json:"value"`
	} `json:"total"`
	Hits     []HitSourceData `json:"hits"`
	MaxScore float64         `json:"max_score"`
}

type HitSourceData struct {
	Index  string     `json:"_index"`
	ID     string     `json:"_id"`
	Score  float64    `json:"_score"`
	Source SourceData `json:"_source"`
}

type SourceData struct {
	//City string `json:"city"`
	//Name string `json:"name"`
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

type Response struct {
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
	TimedOut bool   `json:"timed_out"`
	Took     int    `json:"took"`
}

type Shards struct {
	Failed     int `json:"failed"`
	Skipped    int `json:"skipped"`
	Successful int `json:"successful"`
	Total      int `json:"total"`
}

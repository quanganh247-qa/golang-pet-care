package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"log"
=======
>>>>>>> e859654 (Elastic search)
=======
	"log"
>>>>>>> dc47646 (Optimize SQL query)
=======
>>>>>>> e859654 (Elastic search)

	"github.com/elastic/go-elasticsearch/v8"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ESService struct {
	client   *elasticsearch.Client
	index    string                            // Tiền tố chung, ví dụ: "petclinic"
	mappings map[string]map[string]interface{} // Mapping cho từng loại: "medicines", "diseases", etc.
}

func NewESService(config util.Config) (*ESService, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.ElasticsearchURL},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// Test kết nối
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch connection failed: %w", err)
	}
	defer res.Body.Close()

	mappings := map[string]map[string]interface{}{
		"medicines": {
			"properties": map[string]interface{}{
				"name":          map[string]interface{}{"type": "text"},
				"description":   map[string]interface{}{"type": "text"},
				"side_effects":  map[string]interface{}{"type": "text"},
				"suggest_field": map[string]interface{}{"type": "completion"},
			},
		},
		"diseases": {
			"properties": map[string]interface{}{
				"name":          map[string]interface{}{"type": "text"},
				"symptoms":      map[string]interface{}{"type": "text"},
				"suggest_field": map[string]interface{}{"type": "completion"},
			},
		},
	}

	return &ESService{client: client, index: "petclinic", mappings: mappings}, nil
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (es *ESService) CreateIndices() {
=======
func (es *ESService) CreateIndices() error {
>>>>>>> e859654 (Elastic search)
=======
func (es *ESService) CreateIndices() {
>>>>>>> dc47646 (Optimize SQL query)
=======
func (es *ESService) CreateIndices() error {
>>>>>>> e859654 (Elastic search)
	for docType, mapping := range es.mappings {
		indexName := es.index + "_" + docType

		// Chuẩn bị mapping
		body := map[string]interface{}{
			"mappings": mapping,
		}
		data, err := json.Marshal(body)
		if err != nil {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
			log.Println("error marshaling mapping for %s: %w", docType, err)
=======
			return fmt.Errorf("error marshaling mapping for %s: %w", docType, err)
>>>>>>> e859654 (Elastic search)
=======
			log.Println("error marshaling mapping for %s: %w", docType, err)
>>>>>>> dc47646 (Optimize SQL query)
=======
			return fmt.Errorf("error marshaling mapping for %s: %w", docType, err)
>>>>>>> e859654 (Elastic search)
		}

		// Tạo index
		res, err := es.client.Indices.Create(
			indexName,
			es.client.Indices.Create.WithBody(bytes.NewReader(data)),
		)
		if err != nil {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
			log.Println("error creating index %s: %w", indexName, err)
=======
			return fmt.Errorf("error creating index %s: %w", indexName, err)
>>>>>>> e859654 (Elastic search)
=======
			log.Println("error creating index %s: %w", indexName, err)
>>>>>>> dc47646 (Optimize SQL query)
=======
			return fmt.Errorf("error creating index %s: %w", indexName, err)
>>>>>>> e859654 (Elastic search)
		}
		defer res.Body.Close()

		if res.IsError() {
			var errResponse map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
				log.Println("error parsing create index error for %s: %w", indexName, err)
			}
			log.Println("error creating index %s: %v", indexName, errResponse)
		}
	}
=======
				return fmt.Errorf("error parsing create index error for %s: %w", indexName, err)
=======
				log.Println("error parsing create index error for %s: %w", indexName, err)
>>>>>>> dc47646 (Optimize SQL query)
			}
			log.Println("error creating index %s: %v", indexName, errResponse)
		}
	}
<<<<<<< HEAD
	return nil
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> dc47646 (Optimize SQL query)
=======
				return fmt.Errorf("error parsing create index error for %s: %w", indexName, err)
			}
			return fmt.Errorf("error creating index %s: %v", indexName, errResponse)
		}
	}
	return nil
>>>>>>> e859654 (Elastic search)
}

func (es *ESService) IndexDocument(docType string, id int, doc map[string]interface{}) error {
	indexName := es.index + "_" + docType

	// Convert to JSON
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling %s document: %w", docType, err)
	}

	// Index document
	res, err := es.client.Index(
		indexName,
		bytes.NewReader(data),
		es.client.Index.WithDocumentID(fmt.Sprintf("%d", id)),
		es.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("error indexing %s document: %w", docType, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("error parsing error response for %s: %w", docType, err)
		}
		return fmt.Errorf("error indexing %s document: %v", docType, errResponse)
	}

	return nil
}

// IndexMedicine sử dụng hàm chung
func (es *ESService) IndexMedicine(medicine *db.Medicine) error {
	doc := map[string]interface{}{
		"name": medicine.Name,
		"suggest_field": map[string]interface{}{
			"input":  []string{medicine.Name},
			"weight": 10,
		},
		"description":  medicine.Description.String,
		"side_effects": medicine.SideEffects,
	}
	return es.IndexDocument("medicines", int(medicine.ID), doc)
}

// IndexDisease (ví dụ)
func (es *ESService) IndexDisease(disease *db.Disease) error {
	doc := map[string]interface{}{
		"name": disease.Name,
		"suggest_field": map[string]interface{}{
			"input":  []string{disease.Name},
			"weight": 1,
		},
		"symptoms": disease.Symptoms,
	}
	return es.IndexDocument("diseases", int(disease.ID), doc)
}

func (es *ESService) GetSuggestions(docType, prefix string, size int) ([]string, error) {
	indexName := es.index + "_" + docType

	// Kiểm tra index
	exists, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		return nil, fmt.Errorf("error checking if index %s exists: %w", indexName, err)
	}
	if exists.StatusCode != 200 {
		return nil, fmt.Errorf("index %s does not exist", indexName)
	}

	// Tạo suggestion query
	suggestQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			docType + "-suggest": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field": "suggest_field",
					"size":  size,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(suggestQuery); err != nil {
		return nil, fmt.Errorf("error encoding suggest query for %s: %w", docType, err)
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex(indexName),
		es.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting %s suggestions: %w", docType, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
			return nil, fmt.Errorf("error parsing error response for %s: %w", docType, err)
		}
		return nil, fmt.Errorf("elasticsearch error for %s: %v", docType, errResponse)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing %s suggestions response: %w", docType, err)
	}

	suggestions := []string{}
	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
		if suggestData, ok := suggest[docType+"-suggest"].([]interface{}); ok && len(suggestData) > 0 {
			if firstSuggestion, ok := suggestData[0].(map[string]interface{}); ok {
				if options, ok := firstSuggestion["options"].([]interface{}); ok {
					for _, option := range options {
						if optionMap, ok := option.(map[string]interface{}); ok {
							if text, ok := optionMap["text"].(string); ok {
								suggestions = append(suggestions, text)
							}
						}
					}
				}
			}
		}
	}

	return suggestions, nil
}

// Gọi hàm cho medicines
func (es *ESService) GetMedicineSuggestions(prefix string, size int) ([]string, error) {
	return es.GetSuggestions("medicines", prefix, size)
}

// Gọi hàm cho diseases
func (es *ESService) GetDiseaseSuggestions(prefix string, size int) ([]string, error) {
	return es.GetSuggestions("diseases", prefix, size)
}

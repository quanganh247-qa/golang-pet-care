package elasticsearch

// type ESService struct {
// 	client *elasticsearch.Client
// 	index  string
// }

// // Khởi tạo Elasticsearch client
// func NewESService(config util.Config) (*ESService, error) {
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{config.ElasticsearchURL},
// 	}

// 	client, err := elasticsearch.NewClient(cfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
// 	}

// 	return &ESService{
// 		client: client,
// 		index:  "petclinic", // prefix cho các indices
// 	}, nil
// }

// // Tạo indices cho các loại dữ liệu cần search
// func (es *ESService) CreateIndices() error {
// 	indices := map[string]string{
// 		"diseases": `{
// 			"mappings": {
// 				"properties": {
// 					"name": {
// 						"type": "text",
// 						"analyzer": "standard",
// 						"fields": {
// 							"keyword": {
// 								"type": "keyword"
// 							}
// 						}
// 					},
// 					"suggest_field": {
// 						"type": "completion"
// 					},
// 					"description": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					},
// 					"symptoms": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					},
// 					"treatments": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					}
// 				}
// 			}
// 		}`,
// 		"medicines": `{
// 			"mappings": {
// 				"properties": {
// 					"name": {
// 						"type": "text",
// 						"analyzer": "standard",
// 						"fields": {
// 							"keyword": {
// 								"type": "keyword"
// 							}
// 						}
// 					},
// 					"suggest_field": {
// 						"type": "completion",
// 					},
// 					"usage": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					},
// 					"dosage": {
// 						"type": "text"
// 					},
// 					"side_effects": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					},
// 					"description": {
// 						"type": "text",
// 						"analyzer": "standard"
// 					}
// 				}
// 			}
// 		}`,
// 	}

// 	// Create new index with mapping
// 	for indexName, mapping := range indices {
// 		// Check if the index already exists
// 		if exists, err := es.client.Indices.Exists([]string{es.index + "_" + indexName}); err == nil && exists.StatusCode == 200 {
// 			// Optionally, check the mapping here and decide whether to recreate the index
// 			// If the mapping is correct, skip the creation
// 			continue // Skip creating this index if it already exists
// 		}

// 		// Create the index
// 		res, err := es.client.Indices.Create(
// 			es.index+"_"+indexName,
// 			es.client.Indices.Create.WithBody(strings.NewReader(mapping)),
// 		)
// 		if err != nil {
// 			return fmt.Errorf("error creating index %s: %w", indexName, err)
// 		}
// 		defer res.Body.Close()

// 		if res.IsError() {
// 			return fmt.Errorf("error creating index %s: %s", indexName, res.String())
// 		}
// 	}
// 	return nil
// }

// // Hàm search với highlighting và aggregations
// func (es *ESService) SearchPets(query string) (map[string]interface{}, error) {
// 	// Tạo search query
// 	searchQuery := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"multi_match": map[string]interface{}{
// 				"query":  query,
// 				"fields": []string{"name^2", "symptoms", "diagnosis"},
// 			},
// 		},
// 		// Highlight configuration
// 		"highlight": map[string]interface{}{
// 			"fields": map[string]interface{}{
// 				"name":      map[string]interface{}{},
// 				"symptoms":  map[string]interface{}{},
// 				"diagnosis": map[string]interface{}{},
// 			},
// 			"pre_tags":  []string{"<em>"},
// 			"post_tags": []string{"</em>"},
// 		},
// 		// Aggregations configuration
// 		"aggs": map[string]interface{}{
// 			"pet_types": map[string]interface{}{
// 				"terms": map[string]interface{}{
// 					"field": "type",
// 				},
// 			},
// 		},
// 	}

// 	// Convert query to JSON
// 	var buf strings.Builder
// 	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
// 		return nil, err
// 	}

// 	// Perform search
// 	res, err := es.client.Search(
// 		es.client.Search.WithIndex(es.index+"_pets"),
// 		es.client.Search.WithBody(strings.NewReader(buf.String())),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	// Parse response
// 	var result map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (es *ESService) SearchDiseases(query string, filters map[string]interface{}) (map[string]interface{}, error) {
// 	// Declare searchBody at the function level
// 	searchBody := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"bool": map[string]interface{}{
// 				"must": []map[string]interface{}{
// 					{
// 						"multi_match": map[string]interface{}{
// 							"query": query,
// 							"fields": []string{
// 								"name^3",
// 								"description^2",
// 								"symptoms",
// 							},
// 							"fuzziness": "AUTO",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		"highlight": map[string]interface{}{
// 			"fields": map[string]interface{}{
// 				"name":        map[string]interface{}{},
// 				"symptoms":    map[string]interface{}{},
// 				"description": map[string]interface{}{},
// 			},
// 			"pre_tags":  []string{"<em>"},
// 			"post_tags": []string{"</em>"},
// 		},
// 		"aggs": map[string]interface{}{
// 			"symptoms_count": map[string]interface{}{
// 				"terms": map[string]interface{}{
// 					"field": "symptoms.keyword",
// 					"size":  10,
// 				},
// 			},
// 		},
// 	}

// 	// Apply filters if provided
// 	if symptom, ok := filters["symptoms"]; ok {
// 		boolQuery := searchBody["query"].(map[string]interface{})["bool"].(map[string]interface{})
// 		boolQuery["filter"] = []map[string]interface{}{
// 			{
// 				"term": map[string]interface{}{
// 					"symptoms.keyword": symptom,
// 				},
// 			},
// 		}
// 	}

// 	// Encode searchBody to JSON
// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
// 		return nil, err
// 	}

// 	// Perform the search request
// 	res, err := es.client.Search(
// 		es.client.Search.WithIndex(es.index+"_diseases"), // Use the correct index name
// 		es.client.Search.WithBody(&buf),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	// Check for Elasticsearch errors
// 	if res.IsError() {
// 		var errRes map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
// 			return nil, fmt.Errorf("failed to parse Elasticsearch error: %v", err)
// 		}
// 		return nil, fmt.Errorf("Elasticsearch error: %v", errRes)
// 	}

// 	// Decode the response
// 	var result map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// // IndexDisease indexes a single disease document
// func (es *ESService) IndexDisease(disease *db.Disease) error {
// 	// Prepare the document
// 	doc := map[string]interface{}{
// 		"name":          disease.Name,
// 		"suggest_field": disease.Name,
// 		"description":   disease.Description.String,
// 		"symptoms":      disease.Symptoms,
// 	}

// 	// Convert to JSON
// 	data, err := json.Marshal(doc)
// 	if err != nil {
// 		return fmt.Errorf("error marshaling disease: %w", err)
// 	}

// 	// Index document
// 	res, err := es.client.Index(
// 		es.index+"_diseases",  // Index name
// 		bytes.NewReader(data), // Document body
// 		es.client.Index.WithDocumentID(fmt.Sprintf("%d", disease.ID)), // Use disease ID as document ID
// 		es.client.Index.WithRefresh("true"),                           // Refresh index immediately
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error indexing disease: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		// Read the error response
// 		var errResponse map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
// 			return fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return fmt.Errorf("error indexing document: %v", errResponse)
// 	}

// 	return nil
// }

// // IndexMedicine indexes a single medicine document
// func (es *ESService) IndexMedicine(medicine *db.Medicine) error {
// 	// Prepare the document
// 	doc := map[string]interface{}{
// 		"name": medicine.Name,
// 		"suggest_field": map[string]interface{}{
// 			"input":  []string{medicine.Name}, // Autocomplete input
// 			"weight": 1,                       // Adjust ranking
// 		},
// 		"description":  medicine.Description.String,
// 		"side_effects": medicine.SideEffects,
// 	}

// 	// Convert to JSON
// 	data, err := json.Marshal(doc)
// 	if err != nil {
// 		return fmt.Errorf("error marshaling medicine: %w", err) // Fixed "disease" -> "medicine"
// 	}

// 	// Index document
// 	res, err := es.client.Index(
// 		es.index+"_medicines", // Index name
// 		bytes.NewReader(data), // Document body
// 		es.client.Index.WithDocumentID(fmt.Sprintf("%d", medicine.ID)),
// 		es.client.Index.WithRefresh("true"),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error indexing medicine: %w", err) // Fixed "disease" -> "medicine"
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		// Read the error response
// 		var errResponse map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
// 			return fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return fmt.Errorf("error indexing document: %v", errResponse)
// 	}

// 	return nil
// }

// // DeleteDisease deletes a disease document from the index
// func (es *ESService) DeleteDisease(id int64) error {
// 	res, err := es.client.Delete(
// 		es.index+"_diseases",
// 		fmt.Sprintf("%d", id),
// 		es.client.Delete.WithRefresh("true"),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error deleting disease: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		// Ignore 404 errors (document not found)
// 		if res.StatusCode == 404 {
// 			return nil
// 		}
// 		return fmt.Errorf("error deleting document: %s", res.String())
// 	}

// 	return nil
// }

// // BulkIndexDiseases indexes multiple diseases at once
// func (es *ESService) BulkIndexDiseases(diseases []db.Disease) error {
// 	var buf bytes.Buffer

// 	for _, disease := range diseases {
// 		// Prepare the metadata for bulk indexing
// 		meta := map[string]interface{}{
// 			"index": map[string]interface{}{
// 				"_index": es.index + "_diseases",
// 				"_id":    fmt.Sprintf("%d", disease.ID),
// 			},
// 		}

// 		// Prepare the document
// 		doc := map[string]interface{}{
// 			"id":            disease.ID,
// 			"name":          disease.Name,
// 			"suggest_field": disease.Name,
// 			"description":   disease.Description.String,
// 			"symptoms":      disease.Symptoms,
// 			"created_at":    disease.CreatedAt,
// 			"updated_at":    disease.UpdatedAt,
// 		}

// 		// Add metadata and document to buffer
// 		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
// 			return fmt.Errorf("error encoding meta: %w", err)
// 		}
// 		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
// 			return fmt.Errorf("error encoding doc: %w", err)
// 		}
// 	}

// 	// Perform bulk indexing
// 	res, err := es.client.Bulk(
// 		bytes.NewReader(buf.Bytes()),
// 		es.client.Bulk.WithIndex(es.index+"_diseases"),
// 		es.client.Bulk.WithRefresh("true"),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error bulk indexing: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var raw map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
// 			return fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return fmt.Errorf("error bulk indexing: %v", raw)
// 	}

// 	return nil
// }

// // UpdateDisease updates a disease document in the index
// func (es *ESService) UpdateDisease(disease *db.Disease) error {
// 	// Prepare the update document
// 	doc := map[string]interface{}{
// 		"doc": map[string]interface{}{
// 			"name":        disease.Name,
// 			"description": disease.Description.String,
// 			"symptoms":    disease.Symptoms,
// 			"updated_at":  disease.UpdatedAt,
// 		},
// 	}

// 	data, err := json.Marshal(doc)
// 	if err != nil {
// 		return fmt.Errorf("error marshaling update doc: %w", err)
// 	}

// 	// Perform update
// 	res, err := es.client.Update(
// 		es.index+"_diseases",
// 		fmt.Sprintf("%d", disease.ID),
// 		bytes.NewReader(data),
// 		es.client.Update.WithRefresh("true"),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error updating disease: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var errResponse map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
// 			return fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return fmt.Errorf("error updating document: %v", errResponse)
// 	}

// 	return nil
// }

// // GetSuggestions returns search suggestions based on partial input
// func (es *ESService) GetDiseaseSuggestions(prefix string, size int) ([]string, error) {
// 	// Create suggestion query
// 	suggestQuery := map[string]interface{}{
// 		"suggest": map[string]interface{}{
// 			"disease-suggest": map[string]interface{}{
// 				"prefix": prefix,
// 				"completion": map[string]interface{}{
// 					"field": "suggest_field",
// 					"size":  size,
// 					"contexts": map[string]interface{}{
// 						"category": "disease",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// First verify the index exists with proper mapping
// 	indexName := es.index + "_diseases"

// 	// Check if the index exists
// 	exists, err := es.client.Indices.Exists([]string{indexName})
// 	if err != nil {
// 		return nil, fmt.Errorf("error checking if index exists: %w", err)
// 	}

// 	if exists.StatusCode != 200 {
// 		return nil, fmt.Errorf("index %s does not exist, please call CreateIndices first", indexName)
// 	}

// 	// Print the query for debugging
// 	queryJSON, _ := json.MarshalIndent(suggestQuery, "", "  ")
// 	fmt.Printf("Suggestion query: %s\n", string(queryJSON))

// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(suggestQuery); err != nil {
// 		return nil, fmt.Errorf("error encoding suggest query: %w", err)
// 	}

// 	res, err := es.client.Search(
// 		es.client.Search.WithIndex(es.index+"_diseases"),
// 		es.client.Search.WithBody(&buf),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting suggestions: %w", err)
// 	}
// 	defer res.Body.Close()

// 	responseBody, _ := io.ReadAll(res.Body)
// 	fmt.Printf("Raw response: %s\n", string(responseBody))

// 	// Reset the body for further processing
// 	res.Body = io.NopCloser(bytes.NewReader(responseBody))

// 	if res.IsError() {
// 		var errResponse map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
// 			return nil, fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return nil, fmt.Errorf("elasticsearch error: %v", errResponse)
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("error parsing suggestions response: %w", err)
// 	}

// 	suggestions := []string{}
// 	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
// 		if diseaseSuggest, ok := suggest["disease-suggest"].([]interface{}); ok && len(diseaseSuggest) > 0 {
// 			if firstSuggestion, ok := diseaseSuggest[0].(map[string]interface{}); ok {
// 				if options, ok := firstSuggestion["options"].([]interface{}); ok {
// 					for _, option := range options {
// 						if optionMap, ok := option.(map[string]interface{}); ok {
// 							if text, ok := optionMap["text"].(string); ok {
// 								suggestions = append(suggestions, text)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return suggestions, nil
// }

// // GetSuggestions returns search suggestions based on partial input
// func (es *ESService) GetMedicineSuggestions(prefix string, size int) ([]string, error) {
// 	// Create suggestion query
// 	suggestQuery := map[string]interface{}{
// 		"suggest": map[string]interface{}{
// 			"medicine-suggest": map[string]interface{}{
// 				"prefix": prefix,
// 				"completion": map[string]interface{}{
// 					"field": "suggest_field",
// 					"size":  size,
// 				},
// 			},
// 		},
// 	}

// 	// First verify the index exists with proper mapping
// 	indexName := es.index + "_medicines"

// 	// Check if the index exists
// 	exists, err := es.client.Indices.Exists([]string{indexName})
// 	if err != nil {
// 		return nil, fmt.Errorf("error checking if index exists: %w", err)
// 	}

// 	if exists.StatusCode != 200 {
// 		return nil, fmt.Errorf("index %s does not exist, please call CreateIndices first", indexName)
// 	}

// 	// Print the query for debugging
// 	queryJSON, _ := json.MarshalIndent(suggestQuery, "", "  ")
// 	fmt.Printf("Suggestion query: %s\n", string(queryJSON))

// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(suggestQuery); err != nil {
// 		return nil, fmt.Errorf("error encoding suggest query: %w", err)
// 	}

// 	res, err := es.client.Search(
// 		es.client.Search.WithIndex(es.index+"_medicines"),
// 		es.client.Search.WithBody(&buf),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting suggestions: %w", err)
// 	}
// 	defer res.Body.Close()

// 	responseBody, _ := io.ReadAll(res.Body)
// 	fmt.Printf("Raw response: %s\n", string(responseBody))

// 	// Reset the body for further processing
// 	res.Body = io.NopCloser(bytes.NewReader(responseBody))

// 	if res.IsError() {
// 		var errResponse map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
// 			return nil, fmt.Errorf("error parsing error response: %w", err)
// 		}
// 		return nil, fmt.Errorf("elasticsearch error: %v", errResponse)
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("error parsing suggestions response: %w", err)
// 	}

// 	suggestions := []string{}
// 	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
// 		if medicineSuggest, ok := suggest["medicine-suggest"].([]interface{}); ok && len(medicineSuggest) > 0 {
// 			if firstSuggestion, ok := medicineSuggest[0].(map[string]interface{}); ok {
// 				if options, ok := firstSuggestion["options"].([]interface{}); ok {
// 					for _, option := range options {
// 						if optionMap, ok := option.(map[string]interface{}); ok {
// 							if text, ok := optionMap["text"].(string); ok {
// 								suggestions = append(suggestions, text)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return suggestions, nil
// }

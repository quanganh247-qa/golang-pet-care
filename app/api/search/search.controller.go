package search

// type SearchController struct {
// 	es *elasticsearch.ESService
// }

// func NewSearchController(es *elasticsearch.ESService) *SearchController {
// 	return &SearchController{es: es}
// }

// func (c *SearchController) SearchPets(ctx *gin.Context) {
// 	query := ctx.Query("q")
// 	if query == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
// 		return
// 	}

// 	// results, err := c.es.SearchPets(query)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// ctx.JSON(http.StatusOK, results)

// 	ctx.JSON(http.StatusOK, "ok")
// }

// func (c *SearchController) SearchDiseases(ctx *gin.Context) {
// 	query := ctx.Query("q")
// 	if query == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
// 		return
// 	}

// 	filters := make(map[string]interface{})
// 	if symptom := ctx.Query("symptom"); symptom != "" {
// 		filters["symptoms"] = symptom
// 	}

// 	// results, err := c.es.SearchDiseases(query, filters)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// ctx.JSON(http.StatusOK, results)

// 	ctx.JSON(http.StatusOK, "ok")

// }

// func (c *SearchController) SearchDiseasesSuggestions(ctx *gin.Context) {
// 	query := ctx.Query("q")
// 	if query == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
// 		return
// 	}

// 	size := ctx.Query("size")
// 	if size == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "size parameter is required"})
// 		return
// 	}

// 	sizeInt, err := strconv.Atoi(size)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid size parameter"})
// 		return
// 	}

// 	results, err := c.es.GetDiseaseSuggestions(query, sizeInt)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, results)
// }

// func (c *SearchController) SearchMedicinesSuggestions(ctx *gin.Context) {
// 	query := ctx.Query("q")
// 	if query == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
// 		return
// 	}

// 	size := ctx.Query("size")
// 	if size == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "size parameter is required"})
// 		return
// 	}

// 	sizeInt, err := strconv.Atoi(size)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid size parameter"})
// 		return
// 	}

// 	results, err := c.es.GetMedicineSuggestions(query, sizeInt)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, results)
// }

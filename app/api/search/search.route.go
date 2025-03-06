package search

import (
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService) {
	search := routerGroup.RouterDefault.Group("/search")

	controller := NewSearchController(es)

	search.GET("/diseases", controller.SearchDiseases)
	search.GET("/diseases/suggestions", controller.SearchDiseasesSuggestions)
	search.GET("/medicines/suggestions", controller.SearchMedicinesSuggestions)
	// search.GET("/products", controller.SearchProducts)
}

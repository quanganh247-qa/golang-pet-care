package service_type

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type createServiceTypeRequest struct {
	Servicetypename string `json:"servicetypename"`
	Description     string `json:"description"`
	IconURL         string `json:"iconurl"`
}

type deleteServiceTypeRequest struct {
	Servicetypeid int64 `json:"servicetypeid"`
}

type ServiceTypeController struct {
	service ServiceTypeServiceInterface
}

type ServiceTypeService struct {
	storeDB db.Store
}

// route
type ServiceTypeApi struct {
	controller ServiceTypeControllerInterface
}

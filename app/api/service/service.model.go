package service

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type ServiceService struct {
	storeDB db.Store
}

type createServiceRequest struct {
	TypeID      int8    `json:"type_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Duration    int64   `json:"duration"`
	Description string  `json:"description"`
	Isavailable bool    `json:"isavailable"`
}

type createServiceResponse struct {
	ServiceID   int64   `json:"service_id"`
	TypeID      int64   `json:"type_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Duration    int64   `json:"duration"`
	Description string  `json:"description"`
	Isavailable bool    `json:"isavailable"`
}

type updateServiceRequest struct {
	ServiceID   int64   `json:"service_id"`
	TypeID      int8    `json:"type_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Duration    int64   `json:"duration"`
	Description string  `json:"description"`
	Isavailable bool    `json:"isavailable"`
}

type DeleteServiceRequest struct {
	ServiceID int64 `json:"service_id"`
}

type ServiceController struct {
	service ServiceServiceInterface
}

type ServiceApi struct {
	controller ServiceControllerInterface
}

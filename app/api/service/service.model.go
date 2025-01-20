package service

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type ServiceService struct {
	storeDB db.Store
}

type ServiceController struct {
	service ServiceServiceInterface
}

type ServiceApi struct {
	controller ServiceControllerInterface
}

type CreateServiceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Duration    int     `json:"duration" binding:"required"`
	Cost        float64 `json:"cost" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Notes       string  `json:"notes"`
}

type ServiceRepsonse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Cost        float64 `json:"cost"`
	Category    string  `json:"category"`
	Notes       string  `json:"notes"`
}

type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Cost        float64 `json:"cost"`
	Category    string  `json:"category"`
	Notes       string  `json:"notes"`
}

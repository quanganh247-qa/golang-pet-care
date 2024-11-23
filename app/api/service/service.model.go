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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
type CreateServiceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Duration    int     `json:"duration" binding:"required"`
	Cost        float64 `json:"cost" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Notes       string  `json:"notes"`
<<<<<<< HEAD
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
=======
=======
>>>>>>> cfbe865 (updated service response)
type ServiceResponse struct {
	ServiceTypeID string               `json:"service_type_id"`
	ServiceName   string               `json:"service_name"`
	Service       createServiceRequest `json:"service"`
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
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

<<<<<<< HEAD
=======
}

type ServiceTypeKey struct {
	ID       int64  `json:"id"`
	TypeName string `json:"type_name"`
}

>>>>>>> cfbe865 (updated service response)
type GroupedServiceResponse struct {
	ID       int64                   `json:"id"`
	TypeName string                  `json:"type_name"`
	Services []createServiceResponse `json:"services"`
<<<<<<< HEAD
>>>>>>> cfbe865 (updated service response)
=======
type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Cost        float64 `json:"cost"`
	Category    string  `json:"category"`
	Notes       string  `json:"notes"`
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> cfbe865 (updated service response)
}

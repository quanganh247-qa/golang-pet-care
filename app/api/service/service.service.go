package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ServiceServiceInterface interface {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
	CreateService(ctx *gin.Context, req CreateServiceRequest) (*ServiceRepsonse, error)
	GetAllServices(ctx *gin.Context, pagination *util.Pagination) ([]*ServiceRepsonse, error)
	GetServiceByID(ctx *gin.Context, id int64) (*ServiceRepsonse, error)
	UpdateService(ctx *gin.Context, id int64, req UpdateServiceRequest) (*ServiceRepsonse, error)
	DeleteService(ctx *gin.Context, id int64) error
<<<<<<< HEAD
}

func (s *ServiceService) CreateService(ctx *gin.Context, req CreateServiceRequest) (*ServiceRepsonse, error) {
=======
	createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error)
	deleteServiceService(ctx *gin.Context, serviceID int64) error
	getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error)
	updateServiceService(ctx *gin.Context, serviceid int64, req updateServiceRequest) error
	getServiceByIDService(ctx *gin.Context, serviceid int64) (*createServiceResponse, error)
	getAllServices(ctx *gin.Context, pagination *util.Pagination) ([]createServiceResponse, error)
=======
>>>>>>> b393bb9 (add service and add permission)
}

func (s *ServiceService) CreateService(ctx *gin.Context, req CreateServiceRequest) (*ServiceRepsonse, error) {

<<<<<<< HEAD
=======
	createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error)
	deleteServiceService(ctx *gin.Context, serviceID int64) error
	getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error)
	updateServiceService(ctx *gin.Context, serviceid int64, req updateServiceRequest) error
	getServiceByIDService(ctx *gin.Context, serviceid int64) (*createServiceResponse, error)
}

func (server *ServiceService) createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error) {
	var result db.Service

>>>>>>> c73e2dc (pagination function)
	if req.Name == "" || req.Price == 0 {
		return nil, fmt.Errorf("input name is empty")
	}
>>>>>>> c73e2dc (pagination function)

	var service db.Service
	var err error

<<<<<<< HEAD
=======
	var service db.Service
	var err error

>>>>>>> b393bb9 (add service and add permission)
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		service, err = q.CreateService(ctx, db.CreateServiceParams{
			Name:        pgtype.Text{String: req.Name, Valid: true},
			Description: pgtype.Text{String: req.Description, Valid: true},
			Duration:    pgtype.Int2{Int16: int16(req.Duration), Valid: true},
			Cost:        pgtype.Float8{Float64: req.Cost, Valid: true},
			Category:    pgtype.Text{String: req.Category, Valid: true},
<<<<<<< HEAD
<<<<<<< HEAD
=======
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> ada3717 (Docker file)
		})
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		return nil
<<<<<<< HEAD
	})

=======
	if err != nil {
		return nil, fmt.Errorf("service type not found")
	}

	result, err = server.storeDB.CreateService(ctx, db.CreateServiceParams{
		Typeid:      pgtype.Int8{Int64: int64(req.TypeID), Valid: true},
		Name:        req.Name,
		Price:       pgtype.Float8{Float64: req.Price, Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(req.Duration), Valid: true},
		Description: pgtype.Text{String: req.Description, Valid: true},
		Isavailable: pgtype.Bool{Bool: req.Isavailable, Valid: true},
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
<<<<<<< HEAD
=======
	})

	if err != nil {
		return nil, err
	}
>>>>>>> cfbe865 (updated service response)

	return &result, nil
}
func (server *ServiceService) deleteServiceService(ctx *gin.Context, serviceID int64) error {
	err := server.storeDB.DeleteService(ctx, serviceID)
>>>>>>> cfbe865 (updated service response)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
=======
>>>>>>> b393bb9 (add service and add permission)
	return &ServiceRepsonse{
		ID:          service.ID,
		Name:        service.Name.String,
		Description: service.Description.String,
		Duration:    int(service.Duration.Int16),
		Cost:        service.Cost.Float64,
		Category:    service.Category.String,
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ada3717 (Docker file)
	}, nil
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// Get all services
func (s *ServiceService) GetAllServices(ctx *gin.Context, pagination *util.Pagination) ([]*ServiceRepsonse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	services, err := s.storeDB.GetServices(ctx, db.GetServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
=======
func (server *ServiceService) getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]db.Service, error) {
=======
// return all services with format ServiceResponse

func (server *ServiceService) getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error) {
>>>>>>> cfbe865 (updated service response)

=======
		Notes:       service.Notes.String,
	}, nil
}

// Get all services
func (s *ServiceService) GetAllServices(ctx *gin.Context, pagination *util.Pagination) ([]*ServiceRepsonse, error) {
>>>>>>> b393bb9 (add service and add permission)
	offset := (pagination.Page - 1) * pagination.PageSize

	services, err := s.storeDB.GetServices(ctx, db.GetServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
<<<<<<< HEAD
	}
>>>>>>> c73e2dc (pagination function)
=======
	})
>>>>>>> b393bb9 (add service and add permission)
=======
func (server *ServiceService) getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]db.Service, error) {
=======
// return all services with format ServiceResponse

func (server *ServiceService) getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error) {
>>>>>>> cfbe865 (updated service response)

	offset := (pagination.Page - 1) * pagination.PageSize

	params := db.GetAllServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	}
>>>>>>> c73e2dc (pagination function)

	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}
<<<<<<< HEAD
<<<<<<< HEAD
	var serviceResponses []*ServiceRepsonse
	for _, service := range services {
		serviceResponses = append(serviceResponses, &ServiceRepsonse{
			ID:          service.ID,
			Name:        service.Name.String,
			Description: service.Description.String,
			Duration:    int(service.Duration.Int16),
			Cost:        service.Cost.Float64,
			Category:    service.Category.String,
<<<<<<< HEAD
=======

	serviceMap := make(map[ServiceTypeKey][]createServiceResponse)

=======
	var serviceResponses []*ServiceRepsonse
>>>>>>> b393bb9 (add service and add permission)
	for _, service := range services {
		serviceResponses = append(serviceResponses, &ServiceRepsonse{
			ID:          service.ID,
			Name:        service.Name.String,
			Description: service.Description.String,
<<<<<<< HEAD
			Isavailable: service.Isavailable.Bool,
		}
		serviceMap[serviceTypeKey] = append(serviceMap[serviceTypeKey], serviceResponse)
	}

<<<<<<< HEAD
=======
	serviceMap := make(map[ServiceTypeKey][]createServiceResponse)

	for _, service := range services {

		serviceType, err := server.storeDB.GetServiceType(ctx, service.Typeid.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get service type: %w", err)
		}

		serviceTypeKey := ServiceTypeKey{
			ID:       serviceType.Typeid,
			TypeName: serviceType.Servicetypename,
		}
		serviceResponse := createServiceResponse{
			ServiceID:   service.Serviceid,
			TypeID:      serviceType.Typeid,
			Name:        service.Name,
			Price:       service.Price.Float64,
			Duration:    service.Duration.Microseconds,
			Description: service.Description.String,
			Isavailable: service.Isavailable.Bool,
		}
		serviceMap[serviceTypeKey] = append(serviceMap[serviceTypeKey], serviceResponse)
	}

>>>>>>> cfbe865 (updated service response)
	var result []GroupedServiceResponse
	for serviceType, services := range serviceMap {
		result = append(result, GroupedServiceResponse{
			ID:       serviceType.ID,
			TypeName: serviceType.TypeName,
			Services: services,
<<<<<<< HEAD
>>>>>>> cfbe865 (updated service response)
=======
			Duration:    int(service.Duration.Int16),
			Cost:        service.Cost.Float64,
			Category:    service.Category.String,
			Notes:       service.Notes.String,
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> cfbe865 (updated service response)
		})
	}
	return serviceResponses, nil
}

func (s *ServiceService) GetServiceByID(ctx *gin.Context, id int64) (*ServiceRepsonse, error) {
	service, err := s.storeDB.GetServiceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	return &ServiceRepsonse{
		ID:          service.ID,
		Name:        service.Name.String,
		Description: service.Description.String,
		Duration:    int(service.Duration.Int16),
		Cost:        service.Cost.Float64,
		Category:    service.Category.String,
<<<<<<< HEAD
<<<<<<< HEAD
=======
		Notes:       service.Notes.String,
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> ada3717 (Docker file)
	}, nil
}

func (s *ServiceService) UpdateService(ctx *gin.Context, id int64, req UpdateServiceRequest) (*ServiceRepsonse, error) {

	var service db.Service
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		service, err = q.UpdateService(ctx, db.UpdateServiceParams{
			ID:          id,
			Name:        pgtype.Text{String: req.Name, Valid: true},
			Description: pgtype.Text{String: req.Description, Valid: true},
			Duration:    pgtype.Int2{Int16: int16(req.Duration), Valid: true},
			Cost:        pgtype.Float8{Float64: req.Cost, Valid: true},
			Category:    pgtype.Text{String: req.Category, Valid: true},
<<<<<<< HEAD
<<<<<<< HEAD
=======
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> ada3717 (Docker file)
		})
		if err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}
	return &ServiceRepsonse{
		ID:          service.ID,
		Name:        service.Name.String,
		Description: service.Description.String,
		Duration:    int(service.Duration.Int16),
		Cost:        service.Cost.Float64,
		Category:    service.Category.String,
<<<<<<< HEAD
<<<<<<< HEAD
=======
		Notes:       service.Notes.String,
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> ada3717 (Docker file)
	}, nil
}

func (s *ServiceService) DeleteService(ctx *gin.Context, id int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeleteService(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}
		return nil
<<<<<<< HEAD
	})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

func (server *ServiceService) getAllServices(ctx *gin.Context, pagination *util.Pagination) ([]createServiceResponse, error) {
	var services []createServiceResponse
	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := server.storeDB.GetAllServices(ctx, db.GetAllServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
=======
>>>>>>> b393bb9 (add service and add permission)
	})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

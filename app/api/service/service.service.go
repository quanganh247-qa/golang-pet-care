package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ServiceServiceInterface interface {
	CreateService(ctx *gin.Context, req CreateServiceRequest) (*ServiceRepsonse, error)
	GetAllServices(ctx *gin.Context, pagination *util.Pagination) ([]*ServiceRepsonse, error)
	GetServiceByID(ctx *gin.Context, id int64) (*ServiceRepsonse, error)
	UpdateService(ctx *gin.Context, id int64, req UpdateServiceRequest) (*ServiceRepsonse, error)
	DeleteService(ctx *gin.Context, id int64) error
}

func (s *ServiceService) CreateService(ctx *gin.Context, req CreateServiceRequest) (*ServiceRepsonse, error) {

	var service db.Service
	var err error

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		service, err = q.CreateService(ctx, db.CreateServiceParams{
			Name:        pgtype.Text{String: req.Name, Valid: true},
			Description: pgtype.Text{String: req.Description, Valid: true},
			Duration:    pgtype.Int2{Int16: int16(req.Duration), Valid: true},
			Cost:        pgtype.Float8{Float64: req.Cost, Valid: true},
			Category:    pgtype.Text{String: req.Category, Valid: true},
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	return &ServiceRepsonse{
		ID:          service.ID,
		Name:        service.Name.String,
		Description: service.Description.String,
		Duration:    int(service.Duration.Int16),
		Cost:        service.Cost.Float64,
		Category:    service.Category.String,
		Notes:       service.Notes.String,
	}, nil
}

// Get all services
func (s *ServiceService) GetAllServices(ctx *gin.Context, pagination *util.Pagination) ([]*ServiceRepsonse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	services, err := s.storeDB.GetServices(ctx, db.GetServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}
	var serviceResponses []*ServiceRepsonse
	for _, service := range services {
		serviceResponses = append(serviceResponses, &ServiceRepsonse{
			ID:          service.ID,
			Name:        service.Name.String,
			Description: service.Description.String,
			Duration:    int(service.Duration.Int16),
			Cost:        service.Cost.Float64,
			Category:    service.Category.String,
			Notes:       service.Notes.String,
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
		Notes:       service.Notes.String,
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
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
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
		Notes:       service.Notes.String,
	}, nil
}

func (s *ServiceService) DeleteService(ctx *gin.Context, id int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeleteService(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

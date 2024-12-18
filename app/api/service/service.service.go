package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ServiceServiceInterface interface {
	createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error)
	deleteServiceService(ctx *gin.Context, serviceID int64) error
	getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error)
	updateServiceService(ctx *gin.Context, serviceid int64, req updateServiceRequest) error
	getServiceByIDService(ctx *gin.Context, serviceid int64) (*createServiceResponse, error)
	getAllServices(ctx *gin.Context, pagination *util.Pagination) ([]createServiceResponse, error)
}

func (server *ServiceService) createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error) {
	var result db.Service

	if req.Name == "" || req.Price == 0 {
		return nil, fmt.Errorf("input name is empty")
	}

	_, err := server.storeDB.GetServiceType(ctx, int64(req.TypeID))

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
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
func (server *ServiceService) deleteServiceService(ctx *gin.Context, serviceID int64) error {
	err := server.storeDB.DeleteService(ctx, serviceID)
	if err != nil {
		return err
	}
	return nil
}

// return all services with format ServiceResponse

func (server *ServiceService) getAllServicesService(ctx *gin.Context, pagination *util.Pagination) ([]GroupedServiceResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	params := db.GetAllServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	}

	services, err := server.storeDB.GetAllServices(ctx, params)
	if err != nil {
		return nil, err
	}

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

	var result []GroupedServiceResponse
	for serviceType, services := range serviceMap {
		result = append(result, GroupedServiceResponse{
			ID:       serviceType.ID,
			TypeName: serviceType.TypeName,
			Services: services,
		})
	}

	return result, nil
}

func (server *ServiceService) getServiceByIDService(ctx *gin.Context, serviceID int64) (*createServiceResponse, error) {
	var service createServiceResponse
	err := server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.GetServiceByID(ctx, serviceID)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		service = createServiceResponse{
			ServiceID:   res.Serviceid,
			TypeID:      res.Typeid.Int64,
			Name:        res.Name,
			Price:       res.Price.Float64,
			Duration:    res.Duration.Microseconds,
			Description: res.Description.String,
			Isavailable: res.Isavailable.Bool,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &service, nil
}
func (server *ServiceService) updateServiceService(ctx *gin.Context, serviceid int64, req updateServiceRequest) error {
	_, err := server.storeDB.GetServiceByID(ctx, int64(serviceid))

	if err != nil {
		return fmt.Errorf("service not found")
	} else {
		params := db.UpdateServiceParams{
			Serviceid:   serviceid,
			Typeid:      pgtype.Int8{Int64: int64(req.TypeID), Valid: true},
			Name:        req.Name,
			Price:       pgtype.Float8{Float64: req.Price, Valid: true},
			Duration:    pgtype.Interval{Microseconds: int64(req.Duration), Valid: true},
			Description: pgtype.Text{String: req.Description, Valid: true},
			Isavailable: pgtype.Bool{Bool: req.Isavailable, Valid: true},
		}
		return server.storeDB.UpdateService(ctx, params)
	}
}

func (server *ServiceService) getAllServices(ctx *gin.Context, pagination *util.Pagination) ([]createServiceResponse, error) {
	var services []createServiceResponse
	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := server.storeDB.GetAllServices(ctx, db.GetAllServicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})

	for _, row := range rows {
		service := createServiceResponse{
			ServiceID:   row.Serviceid,
			TypeID:      row.Typeid.Int64,
			Name:        row.Name,
			Price:       row.Price.Float64,
			Duration:    row.Duration.Microseconds,
			Description: row.Description.String,
			Isavailable: row.Isavailable.Bool,
		}
		services = append(services, service)
	}
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return services, nil
}

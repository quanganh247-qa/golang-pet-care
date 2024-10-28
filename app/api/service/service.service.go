package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceServiceInterface interface {
	createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error)
	deleteServiceService(ctx *gin.Context, serviceID int64) error
	getAllServicesService(ctx *gin.Context) ([]db.Service, error)
	updateServiceService(ctx *gin.Context, serviceid int64, req updateServiceRequest) error
	getServiceByIDService(ctx *gin.Context, serviceid int64) (*createServiceResponse, error)
}

func (server *ServiceService) createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error) {
	var result db.Service
	if req.Name == "" || req.Price == 0 {
		return nil, fmt.Errorf("input name is empty")
	}
	fmt.Println("req.TypeID", req.TypeID)

	_, err := server.storeDB.GetServiceType(ctx, int64(req.TypeID))

	if err != nil {
		return nil, fmt.Errorf("service type not found")
	} else {
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

func (server *ServiceService) getAllServicesService(ctx *gin.Context) ([]db.Service, error) {
	services, err := server.storeDB.GetAllServices(ctx)
	if err != nil {
		return nil, err
	}
	return services, nil
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

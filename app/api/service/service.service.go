package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceServiceInterface interface {
	createServiceService(ctx *gin.Context, req createServiceRequest) (*db.Service, error)
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

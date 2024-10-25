package service_type

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceTypeServiceInterface interface {
	createServiceTypeService(ctx *gin.Context, req ServiceType) (*db.Servicetype, error)
}

func (server *ServiceTypeService) createServiceTypeService(ctx *gin.Context, req ServiceType) (*db.Servicetype, error) {

	result, err := server.storeDB.CreateServiceType(ctx, pgtype.Text{String: req.Servicetypename.String, Valid: true})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

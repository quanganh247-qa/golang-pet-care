package service_type

import (
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceTypeServiceInterface interface {
	createServiceTypeService(ctx *gin.Context, req ServiceType) (*db.Servicetype, error)
}

func (server *ServiceTypeService) createServiceTypeService(ctx *gin.Context, req ServiceType) (*db.Servicetype, error) {

	result, err := server.storeDB.CreateServiceType(ctx, req.Servicetypename.String)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

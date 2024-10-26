package service_type

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceTypeServiceInterface interface {
	createServiceTypeService(ctx *gin.Context, req createServiceTypeRequest) (*db.Servicetype, error)
	deleteServiceTypeService(ctx *gin.Context, req deleteServiceTypeRequest) error
}

func (server *ServiceTypeService) createServiceTypeService(ctx *gin.Context, req createServiceTypeRequest) (*db.Servicetype, error) {

	if req.Servicetypename == "" {
		return nil, fmt.Errorf("input servicetypename is empty")
	}

	result, err := server.storeDB.CreateServiceType(ctx, db.CreateServiceTypeParams{
		Servicetypename: req.Servicetypename,
		Description:     pgtype.Text{String: req.Description, Valid: true},
		Iconurl:         pgtype.Text{String: req.IconURL, Valid: true},
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (server *ServiceTypeService) deleteServiceTypeService(ctx *gin.Context, req deleteServiceTypeRequest) error {
	err := server.storeDB.DeleteServiceType(ctx, req.Servicetypeid)
	if err != nil {
		return err
	}
	return nil
}

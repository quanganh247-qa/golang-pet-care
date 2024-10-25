package service_type

import (
	"github.com/jackc/pgx/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ServiceType struct {
	Servicetypename pgtype.Text `json:"servicetypename"`
}

type ServiceTypeController struct {
	service ServiceTypeServiceInterface
}

type ServiceTypeService struct {
	storeDB db.Store
}

// route
type ServiceTypeApi struct {
	controller ServiceTypeControllerInterface
}

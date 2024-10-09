package component

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ComponentApi struct {
	controller ComponentControllerInterface
}
type ComponentController struct {
	service ComponetServiceInterface
}

type ComponentService struct {
	store db.Store
}

type createComponentsRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	ProjectID     int32                  `json:"project_id"`
	Content       map[string]interface{} `json:"content"`
	ComponentCode string                 `json:"component_code"`
}

type createComponentsResonse struct {
	ID            int64
	Name          string
	Description   string
	ProjectID     int32
	Content       string
	ComponentCode string
	CreatedAt     string
	UpdatedAt     string
}
type removedComponentResponse struct {
	ID            int64
	ComponentCode string
	RemovedAt     string
}

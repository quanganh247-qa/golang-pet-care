package page

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PageController struct {
	page PageServiceInterface
}

type PageService struct {
	store db.Store
}

type PageApi struct {
	controller PageControllerInterface
}

type createPageRequest struct {
	ProjectID int32                  `json:"project_id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Content   map[string]interface{} `json:"content"`
}

type createPageResponse struct {
	ID        int64
	ProjectID int32
	Name      string
	Slug      string
	CreatedAt string
	UpdatedAt string
}

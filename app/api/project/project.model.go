package project

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type ProjectService struct {
	storeDB db.Store
}

type ProjectController struct {
	service ProjectServiceInterface
}

type ProjectApi struct {
	controller ProjectControllerInterface
}

type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
type createProjectResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

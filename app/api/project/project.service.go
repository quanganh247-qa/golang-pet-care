package project

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ProjectServiceInterface interface {
	createProject(ctx *gin.Context, req createProjectRequest, username string) (*createProjectResponse, error)
}

func (s *ProjectService) createProject(ctx *gin.Context, req createProjectRequest, username string) (*createProjectResponse, error) {
	var project db.Project
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		project, err = s.storeDB.CreateProject(ctx, db.CreateProjectParams{
			Name:        req.Name,
			Description: pgtype.Text{String: req.Description, Valid: true},
			Username:    username,
		})
		if err != nil {
			return fmt.Errorf("Cannot create project")
		}
		return err
	})
	if err != nil {
		return &createProjectResponse{}, fmt.Errorf("failed to verify email: %w", err)
	}

	return &createProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description.String,
		CreateAt:    project.CreatedAt.Time.Format(time.RFC3339),
		UpdateAt:    project.UpdatedAt.Time.Format(time.RFC3339),
	}, nil

}

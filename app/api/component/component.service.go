package component

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ComponetServiceInterface interface {
	createComponent(ctx *gin.Context, req createComponentsRequest) (*createComponentsResonse, error)
	getComponentByID(ctx *gin.Context, id int64, projectId int32) (*createComponentsResonse, error)
}

func (s *ComponentService) createComponent(ctx *gin.Context, req createComponentsRequest) (*createComponentsResonse, error) {
	contentJson := util.ParseMapInterfaceToString(req.Content)
	slug := util.Slugify(req.Name)
	// Insert the page into the database
	page, err := s.store.CreateComponents(ctx, db.CreateComponentsParams{
		Name:          req.Name,
		ProjectID:     req.ProjectID,
		Description:   pgtype.Text{String: req.Description, Valid: true},
		Content:       pgtype.Text{String: contentJson, Valid: true},
		ComponentCode: pgtype.Text{String: slug, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating page: %v", err) // Handle the error from the database
	}

	// Create response
	return &createComponentsResonse{
		ID:            page.ID,
		ProjectID:     req.ProjectID,
		Name:          page.Name,
		ComponentCode: page.ComponentCode.String,
		CreatedAt:     page.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     page.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *ComponentService) getComponentByID(ctx *gin.Context, id int64, projectId int32) (*createComponentsResonse, error) {
	com, err := s.store.GetComponentsById(ctx, db.GetComponentsByIdParams{
		ID:        id,
		ProjectID: projectId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("Cannot component [%s]", id, err)
		}
	}
	return &createComponentsResonse{
		ID:            com.ID,
		ProjectID:     com.ProjectID,
		Description:   com.Description.String,
		Name:          com.Name,
		Content:       com.Content.String,
		ComponentCode: com.ComponentCode.String,
		CreatedAt:     com.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     com.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}

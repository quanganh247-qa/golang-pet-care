package page

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PageServiceInterface interface {
	createPage(ctx *gin.Context, req createPageRequest) (*createPageResponse, error)
}

func (s *PageService) createPage(ctx *gin.Context, req createPageRequest) (*createPageResponse, error) {
	contentJson, err := json.Marshal(req.Content)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert to json")
	}
	fmt.Println(string(contentJson))
	// Insert the page into the database
	page, err := s.store.CreatePage(ctx, db.CreatePageParams{
		ProjectID: pgtype.Int4{Int32: int32(req.ProjectID), Valid: true},
		Name:      req.Name,
		Slug:      req.Slug,
		Content:   pgtype.Text{String: string(contentJson), Valid: true}, // Pass the valid JSON content to the database
	})
	if err != nil {
		return nil, fmt.Errorf("error creating page: %v", err) // Handle the error from the database
	}

	// Create response
	return &createPageResponse{
		ID:        page.ID,
		ProjectID: req.ProjectID,
		Name:      page.Name,
		Slug:      page.Slug,
		CreatedAt: page.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: page.UpdatedAt.Time.Format(time.RFC3339),
	}, nil

}

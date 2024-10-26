package pet

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PetServiceInterface interface {
	CreatePet(ctx *gin.Context, username string, req createPetRequest) (*createPetResponse, error)
	GetPet(ctx *gin.Context, petid int64) (*createPetResponse, error)
}

func (s *PetService) CreatePet(ctx *gin.Context, username string, req createPetRequest) (*createPetResponse, error) {
	var pet createPetResponse
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.CreatePet(ctx, db.CreatePetParams{
			Username:     username,
			Name:         req.Name,
			Type:         req.Type,
			Breed:        pgtype.Text{String: req.Breed, Valid: true},
			Age:          pgtype.Int4{Int32: int32(req.Age), Valid: true},
			Weight:       pgtype.Float8{Float64: req.Weight, Valid: true},
			Gender:       pgtype.Text{String: req.Gender, Valid: true},
			Healthnotes:  pgtype.Text{String: req.Healthnotes, Valid: true},
			Profileimage: pgtype.Text{String: req.Profileimage, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create pet: %w", err)
		}
		pet = createPetResponse{
			Petid:    res.Petid,
			Username: res.Username,
			Name:     res.Name,
			Type:     res.Type,
			Breed:    res.Breed.String,
			Age:      int16(res.Age.Int32),
			Weight:   res.Weight.Float64,
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &pet, nil
}

func (s *PetService) GetPet(ctx *gin.Context, petid int64) (*createPetResponse, error) {
	var pet createPetResponse
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.GetPetByID(ctx, petid)
		if err != nil {
			return fmt.Errorf("failed to get pet: %w", err)
		}
		pet = createPetResponse{
			Petid:    res.Petid,
			Username: res.Username,
			Name:     res.Name,
			Type:     res.Type,
			Breed:    res.Breed.String,
			Age:      int16(res.Age.Int32),
			Weight:   res.Weight.Float64,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &pet, nil
}

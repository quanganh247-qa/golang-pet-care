package vaccination

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type VaccinationServiceInterface interface {
	CreateVaccination(ctx *gin.Context, req createVaccinationRequest) (*VaccinationResponse, error)
	GetVaccinationByID(ctx *gin.Context, vaccinationID int64) (*VaccinationResponse, error)
	ListVaccinationsByPetID(ctx *gin.Context, petID int64) ([]VaccinationResponse, error)
	UpdateVaccination(ctx *gin.Context, req updateVaccinationRequest) error
	DeleteVaccination(ctx context.Context, vaccinationID int64) error
}

type VaccinationService struct {
	storeDB db.Store
}

func (s *VaccinationService) CreateVaccination(ctx *gin.Context, req createVaccinationRequest) (*VaccinationResponse, error) {
	res, err := s.storeDB.CreateVaccination(ctx, db.CreateVaccinationParams{
		Petid:            pgtype.Int8{Int64: req.PetID, Valid: true},
		VaccineName:      req.VaccineName,
		DateAdministered: pgtype.Timestamp{Time: req.DateAdministered, Valid: true},
		NextDueDate:      pgtype.Timestamp{Time: req.NextDueDate, Valid: !req.NextDueDate.IsZero()},
		VaccineProvider:  pgtype.Text{String: req.VaccineProvider, Valid: req.VaccineProvider != ""},
		BatchNumber:      pgtype.Text{String: req.BatchNumber, Valid: req.BatchNumber != ""},
		Notes:            pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create vaccination: %w", err)
	}

	return &VaccinationResponse{
		VaccinationID:    res.Vaccinationid,
		PetID:            res.Petid.Int64,
		VaccineName:      res.VaccineName,
		DateAdministered: res.DateAdministered.Time,
		NextDueDate:      res.NextDueDate.Time,
		VaccineProvider:  res.VaccineProvider.String,
		BatchNumber:      res.BatchNumber.String,
		Notes:            res.Notes.String,
	}, nil
}

func (s *VaccinationService) GetVaccinationByID(ctx *gin.Context, vaccinationID int64) (*VaccinationResponse, error) {
	res, err := s.storeDB.GetVaccinationByID(ctx, vaccinationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vaccination: %w", err)
	}

	return &VaccinationResponse{
		VaccinationID:    res.Vaccinationid,
		PetID:            res.Petid.Int64,
		VaccineName:      res.VaccineName,
		DateAdministered: res.DateAdministered.Time,
		NextDueDate:      res.NextDueDate.Time,
		VaccineProvider:  res.VaccineProvider.String,
		BatchNumber:      res.BatchNumber.String,
		Notes:            res.Notes.String,
	}, nil
}

func (s *VaccinationService) ListVaccinationsByPetID(ctx *gin.Context, petID int64) ([]VaccinationResponse, error) {
	res, err := s.storeDB.ListVaccinationsByPetID(ctx, pgtype.Int8{Int64: petID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list vaccinations for pet: %w", err)
	}

	var vaccinations []VaccinationResponse
	for _, r := range res {
		vaccinations = append(vaccinations, VaccinationResponse{
			VaccinationID:    r.Vaccinationid,
			PetID:            r.Petid.Int64,
			VaccineName:      r.VaccineName,
			DateAdministered: r.DateAdministered.Time,
			NextDueDate:      r.NextDueDate.Time,
			VaccineProvider:  r.VaccineProvider.String,
			BatchNumber:      r.BatchNumber.String,
			Notes:            r.Notes.String,
		})
	}

	return vaccinations, nil
}

func (s *VaccinationService) UpdateVaccination(ctx *gin.Context, req updateVaccinationRequest) error {
	err := s.storeDB.UpdateVaccination(ctx, db.UpdateVaccinationParams{
		Vaccinationid:    req.VaccinationID,
		VaccineName:      req.VaccineName,
		DateAdministered: pgtype.Timestamp{Time: req.DateAdministered, Valid: true},
		NextDueDate:      pgtype.Timestamp{Time: req.NextDueDate, Valid: !req.NextDueDate.IsZero()},
		VaccineProvider:  pgtype.Text{String: req.VaccineProvider, Valid: req.VaccineProvider != ""},
		BatchNumber:      pgtype.Text{String: req.BatchNumber, Valid: req.BatchNumber != ""},
		Notes:            pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to update vaccination: %w", err)
	}
	return nil
}

func (s *VaccinationService) DeleteVaccination(ctx context.Context, vaccinationID int64) error {
	err := s.storeDB.DeleteVaccination(ctx, vaccinationID)
	if err != nil {
		return fmt.Errorf("failed to delete vaccination: %w", err)
	}
	return nil
}

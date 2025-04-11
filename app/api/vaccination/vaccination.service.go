package vaccination

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type VaccinationServiceInterface interface {
	CreateVaccination(ctx *gin.Context, req createVaccinationRequest) (*VaccinationResponse, error)
	GetVaccinationByID(ctx *gin.Context, vaccinationID int64) (*VaccinationResponse, error)
	ListVaccinationsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]VaccinationResponse, error)
	UpdateVaccination(ctx *gin.Context, req updateVaccinationRequest) error
	DeleteVaccination(ctx context.Context, vaccinationID int64) error
}

type VaccinationService struct {
	storeDB db.Store
}

func (s *VaccinationService) CreateVaccination(ctx *gin.Context, req createVaccinationRequest) (*VaccinationResponse, error) {

	res, err := s.storeDB.CreateVaccination(ctx, db.CreateVaccinationParams{
		Petid:            pgtype.Int8{Int64: req.PetID, Valid: true},
		Vaccinename:      req.VaccineName,
		Dateadministered: pgtype.Timestamp{Time: req.DateAdministered, Valid: true},
		Nextduedate:      pgtype.Timestamp{Time: req.NextDueDate, Valid: !req.NextDueDate.IsZero()},
		Vaccineprovider:  pgtype.Text{String: req.VaccineProvider, Valid: req.VaccineProvider != ""},
		Batchnumber:      pgtype.Text{String: req.BatchNumber, Valid: req.BatchNumber != ""},
		Notes:            pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create vaccination: %w", err)
	}

	return &VaccinationResponse{
		VaccinationID:    res.Vaccinationid,
		PetID:            res.Petid.Int64,
		VaccineName:      res.Vaccinename,
		DateAdministered: res.Dateadministered.Time,
		NextDueDate:      res.Nextduedate.Time,
		VaccineProvider:  res.Vaccineprovider.String,
		BatchNumber:      res.Batchnumber.String,
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
		VaccineName:      res.Vaccinename,
		DateAdministered: res.Dateadministered.Time,
		NextDueDate:      res.Nextduedate.Time,
		VaccineProvider:  res.Vaccineprovider.String,
		BatchNumber:      res.Batchnumber.String,
		Notes:            res.Notes.String,
	}, nil
}

func (s *VaccinationService) ListVaccinationsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]VaccinationResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	res, err := s.storeDB.ListVaccinationsByPetID(ctx, db.ListVaccinationsByPetIDParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
		Petid:  pgtype.Int8{Int64: petID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list vaccinations for pet: %w", err)
	}

	var vaccinations []VaccinationResponse
	for _, r := range res {
		vaccinations = append(vaccinations, VaccinationResponse{
			VaccinationID:    r.Vaccinationid,
			PetID:            r.Petid.Int64,
			VaccineName:      r.Vaccinename,
			DateAdministered: r.Dateadministered.Time,
			NextDueDate:      r.Nextduedate.Time,
			VaccineProvider:  r.Vaccineprovider.String,
			BatchNumber:      r.Batchnumber.String,
			Notes:            r.Notes.String,
		})
	}

	return vaccinations, nil
}

func (s *VaccinationService) UpdateVaccination(ctx *gin.Context, req updateVaccinationRequest) error {
	err := s.storeDB.UpdateVaccination(ctx, db.UpdateVaccinationParams{
		Vaccinename:      req.VaccineName,
		Dateadministered: pgtype.Timestamp{Time: req.DateAdministered, Valid: true},
		Nextduedate:      pgtype.Timestamp{Time: req.NextDueDate, Valid: !req.NextDueDate.IsZero()},
		Vaccineprovider:  pgtype.Text{String: req.VaccineProvider, Valid: req.VaccineProvider != ""},
		Batchnumber:      pgtype.Text{String: req.BatchNumber, Valid: req.BatchNumber != ""},
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

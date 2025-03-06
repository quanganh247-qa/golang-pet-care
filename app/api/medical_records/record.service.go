package medical_records

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/util"
=======
>>>>>>> 3bf345d (happy new year)
=======
	"github.com/quanganh247-qa/go-blog-be/app/util"
>>>>>>> e859654 (Elastic search)
)

type MedicalRecordServiceInterface interface {
	CreateMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	CreateMedicalHistory(ctx *gin.Context, req *MedicalHistoryRequest, recordID int64) (*MedicalHistoryResponse, error)
<<<<<<< HEAD
<<<<<<< HEAD
	GetMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	ListMedicalHistory(ctx *gin.Context, recordID int64, pagination *util.Pagination) ([]MedicalHistoryResponse, error)
	GetMedicalHistoryByID(ctx *gin.Context, medicalHistoryID int64) (*MedicalHistoryResponse, error)
=======
>>>>>>> 3bf345d (happy new year)
=======
	CreateAllergy(ctx *gin.Context, req AllergyRequest, recordID int64) (*Allergy, error)
	GetMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	ListMedicalHistory(ctx *gin.Context, recordID int64, pagination *util.Pagination) ([]MedicalHistoryResponse, error)
	GetMedicalHistoryByID(ctx *gin.Context, medicalHistoryID int64) (*MedicalHistoryResponse, error)
>>>>>>> e859654 (Elastic search)
}

// Quản lý Bệnh án
func (s *MedicalRecordService) CreateMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error) {

	var rec db.MedicalRecord
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		rec, err = q.CreateMedicalRecord(ctx, pgtype.Int8{Int64: petID, Valid: true})
		if err != nil {
			log.Printf("failed to create medical record: %v", err)
			return fmt.Errorf("failed to create medical record: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to create medical record: %v", err)
		return nil, fmt.Errorf("failed to create medical record: %w", err)
	}

	return &MedicalRecordResponse{
		ID:        rec.ID,
		PetID:     rec.PetID.Int64,
		CreatedAt: rec.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: rec.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

// medical history
func (s *MedicalRecordService) CreateMedicalHistory(ctx *gin.Context, req *MedicalHistoryRequest, recordID int64) (*MedicalHistoryResponse, error) {
	var rec db.MedicalHistory
	var err error

	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, req.DiagnosisDate)

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		rec, err = q.CreateMedicalHistory(ctx, db.CreateMedicalHistoryParams{
			MedicalRecordID: pgtype.Int8{Int64: recordID, Valid: true},
			Condition:       pgtype.Text{String: req.Condition, Valid: true},
			DiagnosisDate:   pgtype.Timestamp{Time: t, Valid: true},
			Treatment:       pgtype.Int8{Int64: req.Treatment, Valid: true},
			Notes:           pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			log.Printf("failed to create medical history: %v", err)
			return fmt.Errorf("failed to create medical history: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to create medical history: %v", err)
		return nil, fmt.Errorf("failed to create medical history: %w", err)
	}

	return &MedicalHistoryResponse{
		ID:              rec.ID,
		MedicalRecordID: rec.MedicalRecordID.Int64,
		Condition:       rec.Condition.String,
		DiagnosisDate:   rec.DiagnosisDate.Time.Format("2006-01-02 15:04:05"),
		Treatment:       rec.Treatment.Int64,
		Notes:           rec.Notes.String,
		CreatedAt:       rec.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       rec.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}
<<<<<<< HEAD
<<<<<<< HEAD
=======

func (s *MedicalRecordService) CreateAllergy(ctx *gin.Context, req AllergyRequest, recordID int64) (*Allergy, error) {
	var allergy db.Allergy
	var err error

	allergen, err := json.Marshal(req.Allergen)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal allergen: %w", err)
	}
	reaction, err := json.Marshal(req.Reaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reaction: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		allergy, err = q.CreateAllergy(ctx, db.CreateAllergyParams{
			MedicalRecordID: pgtype.Int8{Int64: recordID, Valid: true},
			Allergen:        allergen,
			Severity:        pgtype.Text{String: req.Severity, Valid: true},
			Reaction:        reaction,
			Notes:           pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create allergy: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create allergy: %w", err)
	}

	return &Allergy{
		ID:              allergy.ID,
		MedicalRecordID: allergy.MedicalRecordID.Int64,
		Allergen:        string(allergen),
		Severity:        allergy.Severity.String,
		Reaction:        string(reaction),
		Notes:           allergy.Notes.String,
	}, nil
}
>>>>>>> e859654 (Elastic search)

func (s *MedicalRecordService) GetMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error) {
	// First verify pet exists and is active
	pet, err := s.storeDB.GetPetByID(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pet: %w", err)
	}
	if !pet.IsActive.Bool {
		return nil, fmt.Errorf("pet is not active")
	}

	record, err := s.storeDB.GetMedicalRecord(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical record: %w", err)
	}

	return &MedicalRecordResponse{
		ID:        record.ID,
		PetID:     record.PetID.Int64,
		CreatedAt: record.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: record.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *MedicalRecordService) ListMedicalHistory(ctx *gin.Context, recordID int64, pagination *util.Pagination) ([]MedicalHistoryResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize
	medicalHistories, err := s.storeDB.GetMedicalHistory(ctx, db.GetMedicalHistoryParams{
		MedicalRecordID: pgtype.Int8{Int64: recordID, Valid: true},
		Limit:           int32(pagination.PageSize),
		Offset:          int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	var medicalHistoryResponse []MedicalHistoryResponse
	for _, medicalHistory := range medicalHistories {
		medicalHistoryResponse = append(medicalHistoryResponse, MedicalHistoryResponse{
			ID:              medicalHistory.ID,
			MedicalRecordID: medicalHistory.MedicalRecordID.Int64,
		})
	}

	return medicalHistoryResponse, nil
}

func (s *MedicalRecordService) GetMedicalHistoryByID(ctx *gin.Context, medicalHistoryID int64) (*MedicalHistoryResponse, error) {
	medicalHistory, err := s.storeDB.GetMedicalHistoryByID(ctx, medicalHistoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	return &MedicalHistoryResponse{
		ID:              medicalHistory.ID,
		MedicalRecordID: medicalHistory.MedicalRecordID.Int64,
	}, nil
}
<<<<<<< HEAD
=======
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)

package medical_records

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type MedicalRecordServiceInterface interface {
	CreateMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	CreateMedicalHistory(ctx *gin.Context, req *MedicalHistoryRequest, recordID int64) (*MedicalHistoryResponse, error)
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

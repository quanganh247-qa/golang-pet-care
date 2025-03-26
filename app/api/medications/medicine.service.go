package medications

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicineServiceInterface interface {
	CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error)
	GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error)
	ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]createMedicineResponse, error)
	UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error
}

func (s *MedicineService) CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error) {

	var medicine db.Medicine
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		var expirationDate time.Time
		if req.ExpirationDate != "" {
			expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
			if err != nil {
				return fmt.Errorf("failed to parse expiration date: %w", err)
			}
		}
		medicine, err = q.CreateMedicine(ctx, db.CreateMedicineParams{
			Name:           req.MedicineName,
			Description:    pgtype.Text{String: req.Description, Valid: true},
			Usage:          pgtype.Text{String: req.Usage, Valid: true},
			Dosage:         pgtype.Text{String: req.Dosage, Valid: true},
			Frequency:      pgtype.Text{String: req.Frequency, Valid: true},
			Duration:       pgtype.Text{String: req.Duration, Valid: true},
			SideEffects:    pgtype.Text{String: req.SideEffects, Valid: true},
			Quantity:       pgtype.Int8{Int64: req.Quantity, Valid: true},
			ExpirationDate: pgtype.Date{Time: expirationDate, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create medicine: %w", err)
		}

		err = s.es.IndexMedicine(&medicine)
		if err != nil {
			return fmt.Errorf("failed to index medicine: %w", err)
		}
		return nil
	})

	return &createMedicineResponse{
		MedicineName:   medicine.Name,
		Description:    medicine.Description.String,
		Usage:          medicine.Usage.String,
		Dosage:         medicine.Dosage.String,
		Frequency:      medicine.Frequency.String,
		Duration:       medicine.Duration.String,
		SideEffects:    medicine.SideEffects.String,
		Quantity:       medicine.Quantity.Int64,
		ExpirationDate: medicine.ExpirationDate.Time.Format("2006-01-02"),
	}, nil
}

func (s *MedicineService) GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error) {
	medicine, err := s.storeDB.GetMedicineByID(ctx, medicineid)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine: %w", err)
	}

	return &createMedicineResponse{
		MedicineName: medicine.Name,
	}, nil
}

func (s *MedicineService) ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]createMedicineResponse, error) {
	// var medicines []createMedicineResponse
	// offset := (pagination.Page - 1) * pagination.PageSize

	// res, err := s.storeDB.ListMedicinesByPet(ctx, db.ListMedicinesByPetParams{
	// 	Limit:  int32(pagination.PageSize),
	// 	Offset: int32(offset),
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list Medicines: %w", err)
	// }

	// for _, r := range res {
	// 	medicines = append(medicines, createMedicineResponse{
	// 		MedicineName: r.Name,
	// 		Dosage:       r.Dosage.String,
	// 		Frequency:    r.Frequency.String,
	// 		Duration:     r.Duration.String,
	// 		SideEffects:  r.SideEffects.String,
	// 		Description:  r.Description.String,
	// 		Usage:        r.Usage.String,
	// 	})
	// }

	// return medicines, nil
	return nil, nil
}

func (s *MedicineService) UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error {

	return nil
}

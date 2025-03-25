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
	// DeleteMedicine(ctx context.Context, Medicineid int64) error
	// ListMedicinesByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]createMedicineResponse, error)
	// SetMedicineInactive(ctx context.Context, Medicineid int64) error
}

func (s *MedicineService) CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error) {
	var medicine createMedicineResponse
	// // Parse the StartDate string to time.Time
	// startDate, err := time.Parse("2006-01-02 15:04:05", req.StartDate)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid start date format: %w", err)
	// }

	// // Parse the EndDate string to time.Time (if provided)
	// var endDate pgtype.Timestamp
	// if req.EndDate != "" {
	// 	parsedEndDate, err := time.Parse("2006-01-02 15:04:05", req.EndDate)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid end date format: %w", err)
	// 	}
	// 	endDate = pgtype.Timestamp{Time: parsedEndDate, Valid: true}
	// } else {
	// 	endDate = pgtype.Timestamp{Valid: false}
	// }
	startDate, endDate, err := util.ParseStringToTime(req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.InsertMedicine(ctx, db.InsertMedicineParams{
			PetID:          req.PetID,
			MedicationName: req.MedicineName,
			Dosage:         req.Dosage,
			Frequency:      req.Frequency,
			StartDate:      pgtype.Timestamp{Time: startDate, Valid: true},
			EndDate:        pgtype.Timestamp{Time: endDate, Valid: true},
			Notes:          pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create Medicine: %w", err)
		}
		medicine = createMedicineResponse{
			MedicineName: res.MedicationName,
			Dosage:       res.Dosage,
			Frequency:    res.Frequency,
			StartDate:    res.StartDate.Time.Format("2006-01-02 15:04:05"),
			EndDate:      res.EndDate.Time.Format("2006-01-02 15:04:05"),
			Notes:        res.Notes.String,
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &medicine, nil
}

func (s *MedicineService) GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error) {
	var medicine createMedicineResponse
	res, err := s.storeDB.GetMedicinesByID(ctx, medicineid)
	if err != nil {
		return nil, fmt.Errorf("failed to get Medicine: %w", err)
	}
	medicine = createMedicineResponse{
		MedicineName: res.MedicationName,
		Dosage:       res.Dosage,
		Frequency:    res.Frequency,
		StartDate:    res.StartDate.Time.Format("2006-01-02 15:04:05"),
		EndDate:      res.EndDate.Time.Format("2006-01-02 15:04:05"),
		Notes:        res.Notes.String,
	}

	return &medicine, nil
}

func (s *MedicineService) ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]createMedicineResponse, error) {
	var medicines []createMedicineResponse
	offset := (pagination.Page - 1) * pagination.PageSize

	res, err := s.storeDB.GetAllMedicinesByPet(ctx, db.GetAllMedicinesByPetParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
		PetID:  petID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list Medicines: %w", err)
	}

	for _, r := range res {
		medicines = append(medicines, createMedicineResponse{
			MedicineName: r.MedicationName,
			Dosage:       r.Dosage,
			Frequency:    r.Frequency,
			StartDate:    r.StartDate.Time.Format("2006-01-02 15:04:05"),
			EndDate:      r.EndDate.Time.Format("2006-01-02 15:04:05"),
			Notes:        r.Notes.String,
		})
	}

	return medicines, nil
}

func (s *MedicineService) UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error {
	fmt.Println(req)
	var err error
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		start, end, err = util.ParseStringToTime(req.StartDate, req.EndDate)
		if err != nil {
			return fmt.Errorf("failed to parse time: %w", err)
		}
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.UpdateMedicine(ctx, db.UpdateMedicineParams{
			MedicationID:   medicineid,
			MedicationName: pgtype.Text{String: req.MedicineName, Valid: true},
			Dosage:         pgtype.Text{String: req.Dosage, Valid: true},
			Frequency:      pgtype.Text{String: req.Frequency, Valid: true},
			StartDate:      pgtype.Timestamp{Time: start, Valid: false},
			EndDate:        pgtype.Timestamp{Time: end, Valid: false},
			Notes:          pgtype.Text{String: req.Notes, Valid: false},
		})
		if err != nil {
			return fmt.Errorf("failed to update Medicine: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// func (s *MedicineService) ListMedicinesByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]createMedicineResponse, error) {
// 	var Medicines []createMedicineResponse
// 	offset := (pagination.Page - 1) * pagination.PageSize

// 	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		listParams := db.ListMedicinesByUsernameParams{
// 			Username: username,
// 			Limit:    int32(pagination.PageSize),
// 			Offset:   int32(offset),
// 		}

// 		res, err := q.ListMedicinesByUsername(ctx, listParams)
// 		if err != nil {
// 			return fmt.Errorf("failed to list Medicines for user %s: %w", username, err)
// 		}

// 		for _, r := range res {
// 			Medicines = append(Medicines, createMedicineResponse{
// 				Medicineid: r.Medicineid,
// 				Username:   r.Username,
// 				Name:       r.Name,
// 				Type:       r.Type,
// 				Breed:      r.Breed.String,
// 				Age:        int16(r.Age.Int32),
// 				Weight:     r.Weight.Float64,
// 			})
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("transaction failed: %w", err)
// 	}

// 	return Medicines, nil
// }

// func (s *MedicineService) SetMedicineInactive(ctx context.Context, Medicineid int64) error {
// 	params := db.SetMedicineInactiveParams{
// 		Medicineid: Medicineid,
// 		IsActive:   pgtype.Bool{Bool: false, Valid: true},
// 	}
// 	return s.storeDB.SetMedicineInactive(ctx, params)
// }

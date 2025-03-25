package medications

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicineServiceInterface interface {
	CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error)
	GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error)
	ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]createMedicineResponse, error)
	UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error
}

func (s *MedicineService) CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error) {

	return nil, nil
}

func (s *MedicineService) GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error) {
	return nil, nil
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

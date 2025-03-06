package disease

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiseaseServiceInterface interface {
	CreateTreatmentService(ctx *gin.Context, treatmentPhase CreateTreatmentRequest) (*db.PetTreatment, error)
	CreateTreatmentPhaseService(ctx *gin.Context, treatmentPhases []CreateTreatmentPhaseRequest, id int64) (*[]db.TreatmentPhase, error)
	AssignMedicinesToTreatmentPhase(ctx *gin.Context, treatmentPhases []AssignMedicineRequest, phaseID int64) (*[]db.PhaseMedicine, error)
	GetTreatmentsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]CreateTreatmentResponse, error)
	GetTreatmentPhasesByTreatmentID(ctx *gin.Context, treatmentID int64, pagination *util.Pagination) (*[]TreatmentPhase, error)
	GetMedicinesByPhase(ctx *gin.Context, phaseID int64, pagination *util.Pagination) ([]PhaseMedicine, error)
	UpdateTreatmentPhaseStatus(ctx *gin.Context, phaseID int64, req UpdateTreatmentPhaseStatusRequest) error
	GetActiveTreatments(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]Treatment, error)
	GetTreatmentProgress(ctx *gin.Context, id int64) ([]TreatmentProgressDetail, error)

	CreateDisease(ctx context.Context, arg CreateDiseaseRequest) (*db.Disease, error)
	// GetDiceaseAnhMedicinesInfoService(ctx *gin.Context, disease string) ([]DiseaseMedicineInfo, error)
	// GetTreatmentByDiseaseID(ctx *gin.Context, diseaseID int64, pagination *util.Pagination) ([]TreatmentPlan, error)
}

func (s *DiseaseService) CreateDisease(ctx context.Context, arg CreateDiseaseRequest) (*db.Disease, error) {
	var disease db.Disease
	var err error

	// Convert []string to JSON bytes
	symptomsJSON, err := json.Marshal(arg.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("error marshaling symptoms: %w", err)
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		disease, err = q.CreateDisease(ctx, db.CreateDiseaseParams{
			Name:        arg.Name,
			Description: pgtype.Text{String: arg.Description, Valid: true},
			Symptoms:    symptomsJSON,
		})
		if err != nil {
			return fmt.Errorf("error while creating disease: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Tự động index vào Elasticsearch
	if err := s.es.IndexDisease(&disease); err != nil {
		log.Printf("Warning: Failed to index disease: %v", err)
	}

	return &disease, nil
}

func (s *DiseaseService) CreateTreatmentService(ctx *gin.Context, treatmentPhase CreateTreatmentRequest) (*db.PetTreatment, error) {
	var PetTreatment db.PetTreatment
	var err error

	startDate, err := time.Parse("2006-01-02", treatmentPhase.StartDate)
	if err != nil {
		log.Println("error while parsing start date: ", err)
		return nil, fmt.Errorf("error while parsing start date: %w", err)
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		PetTreatment, err = q.CreateTreatment(ctx, db.CreateTreatmentParams{
			PetID:     pgtype.Int8{Int64: treatmentPhase.PetID, Valid: true},
			DiseaseID: pgtype.Int8{Int64: treatmentPhase.DiseaseID, Valid: true},
			StartDate: pgtype.Date{Time: startDate, Valid: true},
			Status:    pgtype.Text{String: "ongoing", Valid: true},
			Notes:     pgtype.Text{String: treatmentPhase.Notes, Valid: true},
		})
		if err != nil {
			log.Println("error while creating pet treatment: ", err)
			return fmt.Errorf("error while creating pet treatment: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Println("error while creating pet treatment: ", err)
		return nil, err
	}
	return &PetTreatment, nil
}

// Insert Treatment Phases
func (s *DiseaseService) CreateTreatmentPhaseService(ctx *gin.Context, treatmentPhases []CreateTreatmentPhaseRequest, id int64) (*[]db.TreatmentPhase, error) {
	var insertedPhases []db.TreatmentPhase
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		for _, phase := range treatmentPhases {
			// Parse start date
			startDate, err := time.Parse("2006-01-02", phase.StartDate)
			if err != nil {
				log.Println("error while parsing start date: ", err)
				return fmt.Errorf("error while parsing start date: %w", err)
			}

			// Insert each phase
			treatmentPhase, err := q.CreateTreatmentPhase(ctx, db.CreateTreatmentPhaseParams{
				TreatmentID: pgtype.Int8{Int64: id, Valid: true},
				StartDate:   pgtype.Date{Time: startDate, Valid: true},
				PhaseName:   pgtype.Text{String: phase.PhaseName, Valid: true},
				Description: pgtype.Text{String: phase.Description, Valid: true},
				Status:      pgtype.Text{String: phase.Status, Valid: true},
			})
			if err != nil {
				log.Println("error while creating treatment phase: ", err)
				return fmt.Errorf("error while creating treatment phase: %w", err)
			}

			// Append the inserted phase to the list
			insertedPhases = append(insertedPhases, treatmentPhase)
		}
		return nil
	})
	if err != nil {
		log.Println("error while creating treatment phases: ", err)
		return nil, err
	}
	return &insertedPhases, nil
}

// Assign Medicines to Treatment Phases
func (s *DiseaseService) AssignMedicinesToTreatmentPhase(ctx *gin.Context, treatmentPhases []AssignMedicineRequest, phaseID int64) (*[]db.PhaseMedicine, error) {
	var medicine db.PhaseMedicine
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		for _, phase := range treatmentPhases {
			// Insert each medicine
			medicine, err = q.AssignMedicationToTreatmentPhase(ctx, db.AssignMedicationToTreatmentPhaseParams{
				PhaseID:    phaseID,
				MedicineID: phase.MedicineID,
				Dosage:     pgtype.Text{String: phase.Dosage, Valid: true},
				Frequency:  pgtype.Text{String: phase.Frequency, Valid: true},
				Duration:   pgtype.Text{String: phase.Duration, Valid: true},
				Notes:      pgtype.Text{String: phase.Notes, Valid: true},
			})
			if err != nil {
				log.Println("error while creating treatment medicines: ", err)
				return fmt.Errorf("error while creating treatment medicines: %w", err)
			}
		}
		return nil

	})
	if err != nil {
		log.Println("error while creating treatment medicines: ", err)
		return nil, fmt.Errorf("error while creating treatment medicines: %w", err)
	}
	return &[]db.PhaseMedicine{medicine}, nil

}

// Get All Treatments for a Pet
func (s *DiseaseService) GetTreatmentsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]CreateTreatmentResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize
	treatments, err := s.storeDB.GetTreatmentsByPet(ctx, db.GetTreatmentsByPetParams{
		PetID:  pgtype.Int8{Int64: petID, Valid: true},
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting treatment by pet id: %w", err)
	}

	var result []CreateTreatmentResponse
	for _, treatment := range treatments {
		result = append(result, CreateTreatmentResponse{
			TreatmentID: treatment.TreatmentID,
			Disease:     treatment.Disease,
			StartDate:   treatment.StartDate.Time.Format("2006-01-02"),
			EndDate:     treatment.EndDate.Time.Format("2006-01-02"),
			Status:      treatment.Status.String,
		})
	}
	return &result, nil
}

// Get Treatment Phases for a Treatment
func (s *DiseaseService) GetTreatmentPhasesByTreatmentID(ctx *gin.Context, treatmentID int64, pagination *util.Pagination) (*[]TreatmentPhase, error) {
	offset := (pagination.Page - 1) * pagination.PageSize
	phases, err := s.storeDB.GetTreatmentPhasesByTreatment(ctx, db.GetTreatmentPhasesByTreatmentParams{
		ID:     treatmentID,
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting treatment phases by treatment id: %w", err)
	}
	var result []TreatmentPhase
	for _, phase := range phases {
		result = append(result, TreatmentPhase{
			ID:          phase.ID,
			TreatmentID: treatmentID,
			PhaseName:   phase.PhaseName.String,
			Description: phase.Description.String,
			Status:      phase.Status.String,
			StartDate:   phase.StartDate.Time.Format("2006-01-02"),
			CreatedAt:   phase.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return &result, nil
}

// Get Medicines for a Treatment Phase
func (s *DiseaseService) GetMedicinesByPhase(ctx *gin.Context, phaseID int64, pagination *util.Pagination) ([]PhaseMedicine, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	medicines, err := s.storeDB.GetMedicationsByPhase(ctx, db.GetMedicationsByPhaseParams{
		PhaseID: phaseID,
		Limit:   int32(pagination.PageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		log.Println("error while getting medications for phase: ", err)
		return nil, fmt.Errorf("failed to get medications for phase %d: %w", phaseID, err)
	}

	result := make([]PhaseMedicine, 0, len(medicines))
	for _, medicine := range medicines {
		result = append(result, PhaseMedicine{
			PhaseID:      phaseID,
			MedicineName: medicine.Name,
			Dosage:       medicine.Dosage.String,
			Frequency:    medicine.Frequency.String,
			Duration:     medicine.Duration.String,
			Notes:        medicine.Notes.String,
			CreatedAt:    medicine.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

// Update Treatment Phase Status
func (s *DiseaseService) UpdateTreatmentPhaseStatus(ctx *gin.Context, phaseID int64, req UpdateTreatmentPhaseStatusRequest) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdateTreatmentPhaseStatus(ctx, db.UpdateTreatmentPhaseStatusParams{
			ID:     phaseID,
			Status: pgtype.Text{String: req.Status, Valid: true},
		})
		if err != nil {
			log.Println("error while updating treatment phase status: ", err)
			return fmt.Errorf("error while updating treatment phase status: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Println("error while updating treatment phase status: ", err)
		return fmt.Errorf("error while updating treatment phase status: %w", err)
	}
	return nil
}

// Get All Active Treatments
func (s *DiseaseService) GetActiveTreatments(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]Treatment, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	treatments, err := s.storeDB.GetActiveTreatments(ctx, db.GetActiveTreatmentsParams{
		Petid:  petID,
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		log.Println("error while getting active treatments: ", err)
		return nil, fmt.Errorf("error while getting active treatments: %w", err)
	}

	var result []Treatment
	for _, treatment := range treatments {
		result = append(result, Treatment{
			ID:        treatment.ID,
			PetName:   treatment.PetName,
			Disease:   treatment.Disease,
			StartDate: treatment.StartDate.Time.Format("2006-01-02"),
			EndDate:   treatment.EndDate.Time.Format("2006-01-02"),
			Status:    treatment.Status.String,
		})
	}
	return result, nil
}

// Get Treatment Progress
func (s *DiseaseService) GetTreatmentProgress(ctx *gin.Context, id int64) ([]TreatmentProgressDetail, error) {
	progress, err := s.storeDB.GetTreatmentProgress(ctx, id)
	if err != nil {
		log.Println("error while getting treatment progress: ", err)
		return nil, fmt.Errorf("error while getting treatment progress: %w", err)
	}

	var result []TreatmentProgressDetail
	for _, p := range progress {
		result = append(result, TreatmentProgressDetail{
			PhaseName:    p.PhaseName.String,
			Status:       p.Status.String,
			StartDate:    p.StartDate.Time.Format("2006-01-02"),
			NumMedicines: int32(p.NumMedicines),
		})
	}
	return result, nil
}

// // -- Query lấy phác đồ điều trị đầy đủ
// func (s *DiceaseService) GetDiceaseAnhMedicinesInfoService(ctx *gin.Context, disease string) ([]DiseaseMedicineInfo, error) {

// 	// 	searchPattern := "%" + disease + "%"

// 	// 	rows, err := s.storeDB.GetDiceaseAndMedicinesInfo(ctx, searchPattern)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}

// 	// 	// Group medicines by disease
// 	// 	diseaseMap := make(map[int]*DiseaseMedicineInfo)

// 	// 	for _, row := range rows {

// 	// 		var symptoms []string
// 	// 		if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
// 	// 			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
// 	// 		}

// 	// 		disease, exists := diseaseMap[int(row.DiseaseID)]
// 	// 		if !exists {
// 	// 			disease = &DiseaseMedicineInfo{
// 	// 				DiseaseID:          int(row.DiseaseID),
// 	// 				DiseaseName:        row.DiseaseName,
// 	// 				DiseaseDescription: &row.DiseaseDescription.String,
// 	// 				Symptoms:           symptoms,
// 	// 				Medicines:          []MedicineInfo{},
// 	// 			}
// 	// 			diseaseMap[int(row.DiseaseID)] = disease
// 	// 		}

// 	// 		disease.Medicines = append(disease.Medicines, MedicineInfo{
// 	// 			MedicineID:   int(row.MedicineID.Int64),
// 	// 			MedicineName: row.MedicineName.String,
// 	// 			Usage:        &row.MedicineUsage.String,
// 	// 			Dosage:       &row.Dosage.String,
// 	// 			Frequency:    &row.Frequency.String,
// 	// 			Duration:     &row.Duration.String,
// 	// 			SideEffects:  &row.SideEffects.String,
// 	// 		})
// 	// 	}
// 	// 	// Convert map to slice
// 	// 	result := make([]DiseaseMedicineInfo, 0, len(diseaseMap))
// 	// 	for _, disease := range diseaseMap {
// 	// 		result = append(result, *disease)
// 	// 	}
// 	// 	return result, nil
// 	// }

// 	// func (s *DiceaseService) GetDiseaseTreatmentPlanWithPhasesService(ctx *gin.Context, disease string) ([]TreatmentPlan, error) {

// 	// 	searchPattern := "%" + disease + "%"

// 	// 	rows, err := s.storeDB.GetDiseaseTreatmentPlanWithPhases(ctx, searchPattern)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to get diagnostic")
// 	// 	}

// 	// 	// Group medicines by disease
// 	// 	diseaseMap := make(map[string]*TreatmentPlan)
// 	// 	phaseMap := make(map[string]map[int]*PhaseDetail) // diseaseID_phaseNumber -> PhaseDetail

// 	// 	for _, row := range rows {
// 	// 		var symptoms []string
// 	// 		if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
// 	// 			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
// 	// 		}

// 	// 		// Create or get disease entry
// 	// 		if _, exists := diseaseMap[row.DiseaseName]; !exists {
// 	// 			diseaseMap[row.DiseaseName] = &TreatmentPlan{
// 	// 				DiseaseID:       int(row.DiseaseID),
// 	// 				DiseaseName:     row.DiseaseName,
// 	// 				Description:     row.DiseaseDescription.String,
// 	// 				Symptoms:        symptoms,
// 	// 				TreatmentPhases: []PhaseDetail{},
// 	// 			}
// 	// 		}

// 	// 		// Create phase map for disease if it doesn't exist
// 	// 		diseaseKey := fmt.Sprintf("%s", row.DiseaseName)
// 	// 		if _, exists := phaseMap[diseaseKey]; !exists {
// 	// 			phaseMap[diseaseKey] = make(map[int]*PhaseDetail)
// 	// 		}

// 	// 		// Create or get phase entry
// 	// 		if _, exists := phaseMap[diseaseKey][int(row.PhaseNumber.Int32)]; !exists {
// 	// 			phaseMap[diseaseKey][int(row.PhaseNumber.Int32)] = &PhaseDetail{
// 	// 				PhaseID:     int(row.PhaseID),
// 	// 				PhaseNumber: int(row.PhaseNumber.Int32),
// 	// 				PhaseName:   row.PhaseName.String,
// 	// 				Duration:    row.PhaseDuration.String,
// 	// 				Medicines:   []MedicineInfo{},
// 	// 			}
// 	// 		}

// 	// 		// Add medicine to phase
// 	// 		medicine := MedicineInfo{
// 	// 			MedicineID:   int(row.MedicineID),
// 	// 			MedicineName: row.MedicineName,
// 	// 			Dosage:       &row.Dosage.String,
// 	// 			Usage:        &row.MedicineUsage.String,
// 	// 			Frequency:    &row.Frequency.String,
// 	// 			Duration:     &row.Duration.String,
// 	// 			SideEffects:  &row.SideEffects.String,
// 	// 		}
// 	// 		phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines = append(
// 	// 			phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines,
// 	// 			medicine,
// 	// 		)
// 	// 	}

// 	// 	// Convert maps to slice and organize phases
// 	// 	result := make([]TreatmentPlan, 0, len(diseaseMap))
// 	// 	for diseaseName, plan := range diseaseMap {
// 	// 		diseaseKey := fmt.Sprintf("%s", diseaseName)

// 	// 		// Get all phases for this disease
// 	// 		phases := make([]PhaseDetail, 0, len(phaseMap[diseaseKey]))
// 	// 		for _, phase := range phaseMap[diseaseKey] {
// 	// 			phases = append(phases, *phase)
// 	// 		}

// 	// 		// Sort phases by phase number
// 	// 		sort.Slice(phases, func(i, j int) bool {
// 	// 			return phases[i].PhaseNumber < phases[j].PhaseNumber
// 	// 		})

// 	// 		plan.TreatmentPhases = phases
// 	// 		result = append(result, *plan)
// 	// 	}

// 	// 	// Sort results by disease name for consistency
// 	// 	sort.Slice(result, func(i, j int) bool {
// 	// 		return result[i].DiseaseName < result[j].DiseaseName
// 	// 	})

// 	// return result, nil
// 	return nil, nil
// }

// // get treatment by disease id
// func (s *DiceaseService) GetTreatmentByDiseaseID(ctx *gin.Context, diseaseID int64, pagination *util.Pagination) ([]TreatmentPlan, error) {

// 	// offset := (pagination.Page - 1) * pagination.PageSize

// 	// treatments, err := s.storeDB.GetTreatmentByDiseaseId(ctx, db.GetTreatmentByDiseaseIdParams{
// 	// 	ID:     diseaseID,
// 	// 	Limit:  int32(pagination.PageSize),
// 	// 	Offset: int32(offset),
// 	// })
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("error while getting treatment by disease id: %w", err)
// 	// }
// 	// // / Group medicines by disease
// 	// diseaseMap := make(map[string]*TreatmentPlan)
// 	// phaseMap := make(map[string]map[int]*PhaseDetail) // diseaseID_phaseNumber -> PhaseDetail

// 	// for _, row := range treatments {
// 	// 	var symptoms []string
// 	// 	if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
// 	// 		return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
// 	// 	}

// 	// 	// Create or get disease entry
// 	// 	if _, exists := diseaseMap[row.DiseaseName]; !exists {
// 	// 		diseaseMap[row.DiseaseName] = &TreatmentPlan{
// 	// 			DiseaseID:       int(row.DiseaseID),
// 	// 			DiseaseName:     row.DiseaseName,
// 	// 			Description:     row.DiseaseDescription.String,
// 	// 			Symptoms:        symptoms,
// 	// 			TreatmentPhases: []PhaseDetail{},
// 	// 		}
// 	// 	}

// 	// 	// Create phase map for disease if it doesn't exist
// 	// 	diseaseKey := fmt.Sprintf("%s", row.DiseaseName)
// 	// 	if _, exists := phaseMap[diseaseKey]; !exists {
// 	// 		phaseMap[diseaseKey] = make(map[int]*PhaseDetail)
// 	// 	}

// 	// 	// Create or get phase entry
// 	// 	if _, exists := phaseMap[diseaseKey][int(row.PhaseNumber.Int32)]; !exists {
// 	// 		phaseMap[diseaseKey][int(row.PhaseNumber.Int32)] = &PhaseDetail{
// 	// 			PhaseID:          int(row.PhaseID),
// 	// 			PhaseNumber:      int(row.PhaseNumber.Int32),
// 	// 			PhaseName:        row.PhaseName.String,
// 	// 			PhaseDescription: row.PhaseDescription.String,
// 	// 			Duration:         row.PhaseDuration.String,
// 	// 			PhaseNotes:       row.PhaseNotes.String,
// 	// 			Medicines:        []MedicineInfo{},
// 	// 		}
// 	// 	}

// 	// 	// Add medicine to phase
// 	// 	medicine := MedicineInfo{
// 	// 		MedicineID:   int(row.MedicineID),
// 	// 		MedicineName: row.MedicineName,
// 	// 		Dosage:       &row.Dosage.String,
// 	// 		Usage:        &row.MedicineUsage.String,
// 	// 		Frequency:    &row.Frequency.String,
// 	// 		Duration:     &row.Duration.String,
// 	// 		SideEffects:  &row.SideEffects.String,
// 	// 	}
// 	// 	phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines = append(
// 	// 		phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines,
// 	// 		medicine,
// 	// 	)
// 	// }

// 	// // Convert maps to slice and organize phases
// 	// result := make([]TreatmentPlan, 0, len(diseaseMap))
// 	// for diseaseName, plan := range diseaseMap {
// 	// 	diseaseKey := fmt.Sprintf("%s", diseaseName)

// 	// 	// Get all phases for this disease
// 	// 	phases := make([]PhaseDetail, 0, len(phaseMap[diseaseKey]))
// 	// 	for _, phase := range phaseMap[diseaseKey] {
// 	// 		phases = append(phases, *phase)
// 	// 	}

// 	// 	// Sort phases by phase number
// 	// 	sort.Slice(phases, func(i, j int) bool {
// 	// 		return phases[i].PhaseNumber < phases[j].PhaseNumber
// 	// 	})

// 	// 	plan.TreatmentPhases = phases
// 	// 	result = append(result, *plan)
// 	// }

// 	// // Sort results by disease name for consistency
// 	// sort.Slice(result, func(i, j int) bool {
// 	// 	return result[i].DiseaseName < result[j].DiseaseName
// 	// })

// 	// return result, nil
// 	return nil, nil
// }

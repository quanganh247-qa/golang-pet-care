package disease

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiceaseServiceInterface interface {
	GetDiceaseAnhMedicinesInfoService(ctx *gin.Context, disease string) ([]DiseaseMedicineInfo, error)
	GetDiseaseTreatmentPlanWithPhasesService(ctx *gin.Context, disease string) ([]TreatmentPlan, error)
	GetTreatmentByDiseaseID(ctx *gin.Context, diseaseID int64, pagination *util.Pagination) ([]TreatmentPlan, error)
}

// -- Query lấy phác đồ điều trị đầy đủ
func (s *DiceaseService) GetDiceaseAnhMedicinesInfoService(ctx *gin.Context, disease string) ([]DiseaseMedicineInfo, error) {

	searchPattern := "%" + disease + "%"

	rows, err := s.storeDB.GetDiceaseAndMedicinesInfo(ctx, searchPattern)
	if err != nil {
		return nil, err
	}

	// Group medicines by disease
	diseaseMap := make(map[int]*DiseaseMedicineInfo)

	for _, row := range rows {

		var symptoms []string
		if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
		}

		disease, exists := diseaseMap[int(row.DiseaseID)]
		if !exists {
			disease = &DiseaseMedicineInfo{
				DiseaseID:          int(row.DiseaseID),
				DiseaseName:        row.DiseaseName,
				DiseaseDescription: &row.DiseaseDescription.String,
				Symptoms:           symptoms,
				Medicines:          []MedicineInfo{},
			}
			diseaseMap[int(row.DiseaseID)] = disease
		}

		disease.Medicines = append(disease.Medicines, MedicineInfo{
			MedicineID:   int(row.MedicineID.Int64),
			MedicineName: row.MedicineName.String,
			Usage:        &row.MedicineUsage.String,
			Dosage:       &row.Dosage.String,
			Frequency:    &row.Frequency.String,
			Duration:     &row.Duration.String,
			SideEffects:  &row.SideEffects.String,
		})
	}
	// Convert map to slice
	result := make([]DiseaseMedicineInfo, 0, len(diseaseMap))
	for _, disease := range diseaseMap {
		result = append(result, *disease)
	}
	return result, nil
}

func (s *DiceaseService) GetDiseaseTreatmentPlanWithPhasesService(ctx *gin.Context, disease string) ([]TreatmentPlan, error) {

	searchPattern := "%" + disease + "%"

	rows, err := s.storeDB.GetDiseaseTreatmentPlanWithPhases(ctx, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnostic")
	}

	// Group medicines by disease
	diseaseMap := make(map[string]*TreatmentPlan)
	phaseMap := make(map[string]map[int]*PhaseDetail) // diseaseID_phaseNumber -> PhaseDetail

	for _, row := range rows {
		var symptoms []string
		if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
		}

		// Create or get disease entry
		if _, exists := diseaseMap[row.DiseaseName]; !exists {
			diseaseMap[row.DiseaseName] = &TreatmentPlan{
				DiseaseName:     row.DiseaseName,
				Description:     row.DiseaseDescription.String,
				Symptoms:        symptoms,
				TreatmentPhases: []PhaseDetail{},
			}
		}

		// Create phase map for disease if it doesn't exist
		diseaseKey := fmt.Sprintf("%s", row.DiseaseName)
		if _, exists := phaseMap[diseaseKey]; !exists {
			phaseMap[diseaseKey] = make(map[int]*PhaseDetail)
		}

		// Create or get phase entry
		if _, exists := phaseMap[diseaseKey][int(row.PhaseNumber.Int32)]; !exists {
			phaseMap[diseaseKey][int(row.PhaseNumber.Int32)] = &PhaseDetail{
				PhaseNumber: int(row.PhaseNumber.Int32),
				PhaseName:   row.PhaseName.String,
				Duration:    row.PhaseDuration.String,
				Medicines:   []MedicineInfo{},
			}
		}

		// Add medicine to phase
		medicine := MedicineInfo{
			MedicineID:   int(row.MedicineID),
			MedicineName: row.MedicineName,
			Dosage:       &row.Dosage.String,
			Usage:        &row.MedicineUsage.String,
			Frequency:    &row.Frequency.String,
			Duration:     &row.Duration.String,
			SideEffects:  &row.SideEffects.String,
		}
		phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines = append(
			phaseMap[diseaseKey][int(row.PhaseNumber.Int32)].Medicines,
			medicine,
		)
	}

	// Convert maps to slice and organize phases
	result := make([]TreatmentPlan, 0, len(diseaseMap))
	for diseaseName, plan := range diseaseMap {
		diseaseKey := fmt.Sprintf("%s", diseaseName)

		// Get all phases for this disease
		phases := make([]PhaseDetail, 0, len(phaseMap[diseaseKey]))
		for _, phase := range phaseMap[diseaseKey] {
			phases = append(phases, *phase)
		}

		// Sort phases by phase number
		sort.Slice(phases, func(i, j int) bool {
			return phases[i].PhaseNumber < phases[j].PhaseNumber
		})

		plan.TreatmentPhases = phases
		result = append(result, *plan)
	}

	// Sort results by disease name for consistency
	sort.Slice(result, func(i, j int) bool {
		return result[i].DiseaseName < result[j].DiseaseName
	})

	return result, nil
}

// get treatment by disease id
func (s *DiceaseService) GetTreatmentByDiseaseID(ctx *gin.Context, diseaseID int64, pagination *util.Pagination) ([]TreatmentPlan, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	treatments, err := s.storeDB.GetTreatmentByDiseaseId(ctx, db.GetTreatmentByDiseaseIdParams{
		ID:     diseaseID,
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting treatment by disease id: %w", err)
	}
	diseaseMap := make(map[string]*TreatmentPlan)
	phaseMap := make(map[string]map[int]*PhaseDetail) // diseaseID_phaseNumber -> PhaseDetail

	for _, row := range treatments {
		var symptoms []string
		if err := json.Unmarshal(row.Symptoms, &symptoms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
		}

		// Create or get disease entry
		if _, exists := diseaseMap[row.DiseaseName]; !exists {
			diseaseMap[row.DiseaseName] = &TreatmentPlan{
				DiseaseID:       int(row.DiseaseID),
				DiseaseName:     row.DiseaseName,
				Description:     row.DiseaseDescription.String,
				Symptoms:        symptoms,
				TreatmentPhases: []PhaseDetail{},
			}
		}

		// Create phase map for disease if it doesn't exist
		diseaseKey := fmt.Sprintf("%s", row.DiseaseName)
		if _, exists := phaseMap[diseaseKey]; !exists {
			phaseMap[diseaseKey] = make(map[int]*PhaseDetail)
		}

		// Create or get phase entry
		if _, exists := phaseMap[diseaseKey][int(row.PhaseNumber.Int32)]; !exists {
			phaseMap[diseaseKey][int(row.PhaseNumber.Int32)] = &PhaseDetail{
				PhaseID:          int(row.PhaseID),
				PhaseNumber:      int(row.PhaseNumber.Int32),
				PhaseName:        row.PhaseName.String,
				Duration:         row.PhaseDuration.String,
				PhaseDescription: row.PhaseDescription.String,
				PhaseNotes:       row.PhaseNotes.String,
			}
		}
	}

	// Convert maps to slice and organize phases
	result := make([]TreatmentPlan, 0, len(diseaseMap))
	for diseaseName, plan := range diseaseMap {
		diseaseKey := fmt.Sprintf("%s", diseaseName)

		// Get all phases for this disease
		phases := make([]PhaseDetail, 0, len(phaseMap[diseaseKey]))
		for _, phase := range phaseMap[diseaseKey] {
			phases = append(phases, *phase)
		}

		// Sort phases by phase number
		sort.Slice(phases, func(i, j int) bool {
			return phases[i].PhaseNumber < phases[j].PhaseNumber
		})

		plan.TreatmentPhases = phases
		result = append(result, *plan)
	}

	// Sort results by disease name for consistency
	sort.Slice(result, func(i, j int) bool {
		return result[i].DiseaseName < result[j].DiseaseName
	})

	return result, nil
}

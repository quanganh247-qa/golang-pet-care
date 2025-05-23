package disease

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/minio"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiseaseServiceInterface interface {
	CreateTreatmentService(ctx *gin.Context, treatmentPhase CreateTreatmentRequest) (*db.PetTreatment, error)
	CreateTreatmentPhaseService(ctx *gin.Context, treatmentPhases []CreateTreatmentPhaseRequest, id int64) (*[]db.TreatmentPhase, error)
	AssignMedicinesToTreatmentPhase(ctx *gin.Context, treatmentPhases []AssignMedicineRequest, phaseID int64) (*[]db.PhaseMedicine, error)
	GetTreatmentsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]Treatment, error)
	GetTreatmentPhasesByTreatmentID(ctx *gin.Context, treatmentID int64, pagination *util.Pagination) (*[]TreatmentPhase, error)
	GetMedicinesByPhase(ctx *gin.Context, phaseID int64, pagination *util.Pagination) ([]PhaseMedicine, error)
	UpdateTreatmentPhaseStatus(ctx *gin.Context, phaseID int64, req UpdateTreatmentPhaseStatusRequest) error
	GetActiveTreatments(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]Treatment, error)
	GetTreatmentProgress(ctx *gin.Context, id int64) ([]TreatmentProgressDetail, error)
	UpdateTreatmentStatus(ctx *gin.Context, treatmentID int64, status string) error

	CreateDisease(ctx context.Context, arg CreateDiseaseRequest) (*db.Disease, error)
	GenerateMedicineOnlyPrescriptionPDF(ctx context.Context, treatmentID int64, outputFile string) (*PrescriptionResponse, error)

	// CreateAllergyService(ctx *gin.Context, petID int64, req CreateAllergyRequest) (*PetAllergy, error)
	// GetAllergiesByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]PetAllergy, error)
}

func (s *DiseaseService) CreateDisease(ctx context.Context, arg CreateDiseaseRequest) (*db.Disease, error) {
	var disease db.Disease
	var err error

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
			PetID:       pgtype.Int8{Int64: treatmentPhase.PetID, Valid: true},
			Name:        pgtype.Text{String: treatmentPhase.Name, Valid: true},
			Diseases:    pgtype.Text{String: treatmentPhase.Diseases, Valid: true},
			DoctorID:    pgtype.Int4{Int32: int32(treatmentPhase.DoctorID), Valid: true},
			Type:        pgtype.Text{String: treatmentPhase.Type, Valid: true},
			Status:      pgtype.Text{String: treatmentPhase.Status, Valid: true},
			StartDate:   pgtype.Date{Time: startDate, Valid: true},
			Description: pgtype.Text{String: treatmentPhase.Notes, Valid: true},
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
				Quantity:   pgtype.Int4{Int32: int32(phase.Quantity), Valid: true},
				Dosage:     pgtype.Text{String: phase.Dosage, Valid: true},
				Frequency:  pgtype.Text{String: phase.Frequency, Valid: true},
				Duration:   pgtype.Text{String: phase.Duration, Valid: true},
				Notes:      pgtype.Text{String: phase.Notes, Valid: true},
			})
			if err != nil {
				log.Println("error while creating treatment medicines: ", err)
				return fmt.Errorf("error while creating treatment medicines: %w", err)
			}

			medicineDetail, err := q.GetMedicineByID(ctx, phase.MedicineID)
			if err != nil {
				return fmt.Errorf("error while getting medicine: %w", err)
			}

			totalAmount := medicineDetail.UnitPrice.Float64 * float64(phase.Quantity)

			// Create export medicine transaction to track inventory
			_, err = q.CreateMedicineTransaction(ctx, db.CreateMedicineTransactionParams{
				MedicineID:      phase.MedicineID,
				Quantity:        phase.Quantity,
				TransactionType: "export",
				Notes:           pgtype.Text{String: fmt.Sprintf("Assigned to treatment phase ID: %d", phaseID), Valid: true},
				UnitPrice:       pgtype.Float8{Float64: medicineDetail.UnitPrice.Float64, Valid: true},
				TotalAmount:     pgtype.Float8{Float64: totalAmount, Valid: true},
				SupplierID:      pgtype.Int8{Int64: medicineDetail.SupplierID.Int64, Valid: true},
				ExpirationDate:  pgtype.Date{Time: medicineDetail.ExpirationDate.Time, Valid: true},
				// PrescriptionID:  pgtype.Int8{Int64: phase.AppointmentID, Valid: true},
				AppointmentID: pgtype.Int8{Int64: phase.AppointmentID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("error while creating export medicine transaction: %w", err)
			}

			// Update medicine quantity in inventory
			err = q.UpdateMedicineQuantity(ctx, db.UpdateMedicineQuantityParams{
				ID:       phase.MedicineID,
				Quantity: pgtype.Int8{Int64: medicineDetail.Quantity.Int64 - phase.Quantity, Valid: true},
			})
			if err != nil {
				log.Println("error while updating medicine quantity: ", err)
				return fmt.Errorf("error while updating medicine quantity: %w", err)
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
func (s *DiseaseService) GetTreatmentsByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]Treatment, error) {
	// Calculate pagination offset
	offset := (pagination.Page - 1) * pagination.PageSize

	// Fetch treatments for the pet
	treatments, err := s.storeDB.GetTreatmentsByPet(ctx, db.GetTreatmentsByPetParams{
		PetID:  pgtype.Int8{Int64: petID, Valid: true},
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting treatment by pet id: %w", err)
	}

	// Initialize the result slice
	var result []Treatment

	// Iterate over each treatment
	for _, treatment := range treatments {
		// Fetch the doctor for the treatment
		doctor, err := s.storeDB.GetDoctor(ctx, int64(treatment.DoctorID.Int32))
		if err != nil {
			log.Println("error while getting doctor: ", err)
			return nil, fmt.Errorf("error while getting doctor: %w", err)
		}

		// Fetch all phases for the treatment
		phases, err := s.storeDB.GetAllTreatmentPhasesByTreatmentID(ctx, pgtype.Int8{Int64: treatment.ID, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("error while getting treatment phases by treatment id: %w", err)
		}

		// Initialize slice for treatment phases
		var phaseResult []TreatmentPhase

		// Iterate over each phase
		for _, phase := range phases {
			// Fetch medications for the phase
			medicines, err := s.storeDB.GetMedicationsByPhase(ctx, phase.ID)
			if err != nil {
				log.Println("error while getting medications for phase: ", err)
				return nil, fmt.Errorf("failed to get medications for phase %d: %w", phase.ID, err)
			}

			// Initialize slice for phase medications
			var meds []PhaseMedicine

			// Construct PhaseMedicine structs for each medication
			for _, medicine := range medicines {
				meds = append(meds, PhaseMedicine{
					PhaseID:      phase.ID,
					MedicineID:   medicine.ID,
					MedicineName: medicine.Name,
					Dosage:       medicine.Dosage.String,
					Frequency:    medicine.Frequency.String,
					Duration:     medicine.Duration.String,
					Notes:        medicine.Notes.String,
					CreatedAt:    medicine.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				})
			}

			// Construct TreatmentPhase struct with medications
			phaseResult = append(phaseResult, TreatmentPhase{
				ID:          phase.ID,
				TreatmentID: treatment.ID,
				PhaseName:   phase.PhaseName.String,
				Description: phase.Description.String,
				Status:      phase.Status.String,
				StartDate:   phase.StartDate.Time.Format("2006-01-02"),
				CreatedAt:   phase.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				Medications: meds,
			})
		}

		// Construct Treatment struct with phases
		result = append(result, Treatment{
			ID:          treatment.ID,
			Type:        treatment.Type.String,
			Disease:     treatment.Diseases.String,
			StartDate:   treatment.StartDate.Time.Format("2006-01-02"),
			EndDate:     treatment.EndDate.Time.Format("2006-01-02"),
			Status:      treatment.Status.String,
			DoctorName:  doctor.Name,
			Name:        treatment.Name.String,
			Phases:      phaseResult,
			Description: treatment.Description.String,
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
		medicines, err := s.storeDB.GetMedicationsByPhase(ctx, phase.ID)
		if err != nil {
			log.Println("error while getting medications for phase: ", err)
			return nil, fmt.Errorf("failed to get medications for phase %d: %w", phase.ID, err)
		}

		meds := make([]PhaseMedicine, 0, len(medicines))
		for _, medicine := range medicines {
			meds = append(meds, PhaseMedicine{
				PhaseID:      phase.ID,
				MedicineName: medicine.Name,
				Dosage:       medicine.Dosage.String,
				Quantity:     int64(medicine.Quantity.Int32),
				MedicineID:   medicine.ID,
				Frequency:    medicine.Frequency.String,
				Duration:     medicine.Duration.String,
				Notes:        medicine.Notes.String,
				CreatedAt:    medicine.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			})
		}
		result = append(result, TreatmentPhase{
			ID:          phase.ID,
			TreatmentID: treatmentID,
			PhaseName:   phase.PhaseName.String,
			Description: phase.Description.String,
			Status:      phase.Status.String,
			StartDate:   phase.StartDate.Time.Format("2006-01-02"),
			CreatedAt:   phase.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			Medications: meds,
		})
	}

	return &result, nil
}

// Get Medicines for a Treatment Phase
func (s *DiseaseService) GetMedicinesByPhase(ctx *gin.Context, phaseID int64, pagination *util.Pagination) ([]PhaseMedicine, error) {

	// offset := (pagination.Page - 1) * pagination.PageSize

	// medicines, err := s.storeDB.GetMedicationsByPhase(ctx, db.GetMedicationsByPhaseParams{
	// 	PhaseID: phaseID,
	// 	Limit:   int32(pagination.PageSize),
	// 	Offset:  int32(offset),
	// })
	medicines, err := s.storeDB.GetMedicationsByPhase(ctx, phaseID)
	if err != nil {
		log.Println("error while getting medications for phase: ", err)
		return nil, fmt.Errorf("failed to get medications for phase %d: %w", phaseID, err)
	}

	result := make([]PhaseMedicine, 0, len(medicines))
	for _, medicine := range medicines {
		result = append(result, PhaseMedicine{
			PhaseID:      phaseID,
			MedicineID:   medicine.ID,
			MedicineName: medicine.Name,
			Dosage:       medicine.Dosage.String,
			Frequency:    medicine.Frequency.String,
			Duration:     medicine.Duration.String,
			Quantity:     int64(medicine.Quantity.Int32),
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

func (s *DiseaseService) GenerateMedicineOnlyPrescriptionPDF(ctx context.Context, treatmentID int64, outputFile string) (*PrescriptionResponse, error) {
	// Fetch all required data in a single transaction
	var prescription Prescription
	var pdfBytes bytes.Buffer
	var pet db.Pet
	var err error
	var wg sync.WaitGroup
	var treatment db.PetTreatment
	var errs []error

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		wg.Add(2)
		go func() {
			defer wg.Done()
			// 1. Fetch treatment with related data
			treatment, err = q.GetTreatment(ctx, treatmentID)
			if err != nil {
				errs = append(errs, err)
				return
			}

		}()
		go func() {
			defer wg.Done()
			pet, err = q.GetPetByID(ctx, treatment.PetID.Int64)
			if err != nil {
				errs = append(errs, err)
				return
			}
		}()
		wg.Wait()

		if len(errs) > 0 {
			return errs[0]
		}

		// 3. Fetch doctor details
		doctor, err := q.GetDoctor(ctx, int64(treatment.DoctorID.Int32))
		if err != nil {
			return fmt.Errorf("failed to get doctor: %w", err)
		}

		// 4. Fetch medicines in one call
		medicines, err := q.GetMedicineByTreatmentID(ctx, pgtype.Int8{Int64: treatmentID, Valid: true})
		if err != nil {
			return fmt.Errorf("failed to get medicines: %w", err)
		}

		// Build prescription object
		prescription = Prescription{
			ID:              fmt.Sprintf("DT-%d", treatmentID),
			HospitalLogo:    "",
			HospitalName:    hospitalName,
			HospitalAddress: hospitalAddress,
			HospitalPhone:   hospitalPhone,
			PatientName:     pet.Name,
			PatientGender:   pet.Gender.String,
			PatientAge:      int(pet.Age.Int32),
			Diagnosis:       fmt.Sprintf("Disease ID: %d", treatment.Diseases),
			Notes:           treatment.Description.String,
			PrescribedDate:  treatment.StartDate.Time,
			DoctorName:      doctor.Name,
			DoctorTitle:     doctor.Specialization.String,
		}

		// Process medicines (fetch full details for each)
		for _, med := range medicines {
			medicine, err := q.GetMedicineByID(ctx, med.MedicineID)
			if err != nil {
				return fmt.Errorf("failed to get medicine details: %w", err)
			}

			prescription.Medicines = append(prescription.Medicines, MedicineInfo{
				MedicineName: medicine.Name,
				Dosage:       med.Dosage.String,
				Frequency:    med.Frequency.String,
				Duration:     med.Duration.String,
				Notes:        med.Notes.String,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Generate PDF
	if err := GeneratePrescriptionPDF(&prescription, &pdfBytes); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload to MinIO
	filePath := fmt.Sprintf("%s/files/treatment-%d.pdf", pet.Name, treatmentID)
	fileName := fmt.Sprintf("treatment-%d.pdf", treatmentID)
	fileData := pdfBytes.Bytes()

	mcclient, err := minio.GetMinIOClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get MinIO client: %w", err)
	}

	// Upload file
	if err := mcclient.UploadFile(ctx, "prescriptions", filePath, fileData); err != nil {
		return nil, fmt.Errorf("failed to upload PDF to MinIO: %w", err)
	}

	// Get URL
	fileURL, err := mcclient.GetPresignedURL(ctx, "prescriptions", filePath, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for file: %w", err)
	}

	// Create file record
	file, err := s.storeDB.CreateFile(ctx, db.CreateFileParams{
		FileName: fileName,
		FilePath: filePath,
		FileSize: int64(len(fileData)),
		FileType: "pdf",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return &PrescriptionResponse{
		PrescriptionID:  file.ID,
		PrescriptionURL: fileURL,
	}, nil
}

// Allergy
// func (s *DiseaseService) CreateAllergyService(ctx *gin.Context, petID int64, req CreateAllergyRequest) (*PetAllergy, error) {
// 	var res PetAllergy
// 	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		// 2. Create allergy
// 		allergy, err := q.CreatePetAllergy(ctx, db.CreatePetAllergyParams{
// 			PetID:  pgtype.Int8{Int64: petID, Valid: true},
// 			Type:   pgtype.Text{String: req.Type, Valid: true},
// 			Detail: pgtype.Text{String: req.Detail, Valid: true},
// 		})
// 		if err != nil {
// 			return fmt.Errorf("error while creating allergy: %w", err)
// 		}
// 		res = PetAllergy{
// 			ID:     allergy.ID,
// 			PetID:  petID,
// 			Type:   allergy.Type.String,
// 			Detail: allergy.Detail.String,
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error while creating allergy: %w", err)
// 	}
// 	return &res, nil
// }

// func (s *DiseaseService) GetAllergiesByPetID(ctx *gin.Context, petID int64, pagination *util.Pagination) (*[]PetAllergy, error) {

// 	offset := (pagination.Page - 1) * pagination.PageSize
// 	allergies, err := s.storeDB.ListPetAllergies(ctx, db.ListPetAllergiesParams{
// 		PetID:  pgtype.Int8{Int64: petID, Valid: true},
// 		Limit:  int32(pagination.PageSize),
// 		Offset: int32(offset),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error while getting allergies by pet id: %w", err)
// 	}
// 	var res []PetAllergy
// 	for _, allergy := range allergies {
// 		res = append(res, PetAllergy{
// 			ID:     allergy.ID,
// 			PetID:  petID,
// 			Type:   allergy.Type.String,
// 			Detail: allergy.Detail.String,
// 		})
// 	}
// 	return &res, nil
// }

func (s *DiseaseService) UpdateTreatmentStatus(ctx *gin.Context, treatmentID int64, status string) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdateTreatmentStatus(ctx, db.UpdateTreatmentStatusParams{
			ID:     treatmentID,
			Status: pgtype.Text{String: status, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error while updating treatment status: %w", err)
		}
		return nil
	})
}

package medical_records

import (
	"errors"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

var (
	ErrMedicalRecordNotFound = errors.New("medical record not found")
	ErrInvalidDiagnosisDate  = errors.New("invalid diagnosis date format")
	ErrInvalidPetID          = errors.New("invalid pet ID")
	ErrPrescriptionNotFound  = errors.New("prescription not found")
	ErrTestResultNotFound    = errors.New("test result not found")
)

type MedicalRecordApi struct {
	controller MedicalRecordControllerInterface
}

type MedicalRecordController struct {
	service MedicalRecordServiceInterface
}

type MedicalRecordService struct {
	storeDB db.Store
}
type MedicalRecord struct {
	ID        string    `json:"id"`
	PetID     string    `json:"pet_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MedicalRecordResponse struct {
	ID        int64  `json:"id"`
	PetID     int64  `json:"pet_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MedicalHistoryRequest struct {
	Condition     string `json:"condition" binding:"required"`
	DiagnosisDate string `json:"diagnosis_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Notes         string `json:"notes"`
}

type MedicalHistoryResponse struct {
	ID              int64  `json:"id"`
	MedicalRecordID int64  `json:"medical_record_id"`
	Condition       string `json:"condition"`
	DiagnosisDate   string `json:"diagnosis_date"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type Allergy struct {
	ID              int64  `json:"id"`
	MedicalRecordID int64  `json:"medical_record_id"`
	Allergen        string `json:"allergen"`
	Severity        string `json:"severity"`
	Reaction        string `json:"reaction"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type AllergyRequest struct {
	Allergen string `json:"allergen" binding:"required"`
	Severity string `json:"severity" binding:"required,oneof=mild moderate severe"`
	Reaction string `json:"reaction" binding:"required"`
	Notes    string `json:"notes"`
}

// New structures for complete medical history

// Examination represents a medical examination for a pet
type ExaminationRequest struct {
	MedicalHistoryID int64  `json:"medical_history_id" binding:"required"`
	ExamDate         string `json:"exam_date" binding:"required,datetime=2006-01-02 15:04:05"`
	ExamType         string `json:"exam_type" binding:"required"`
	Findings         string `json:"findings" binding:"required"`
	VetNotes         string `json:"vet_notes"`
	DoctorID         int64  `json:"doctor_id" binding:"required"`
}

type ExaminationResponse struct {
	ID               int64  `json:"id"`
	MedicalHistoryID int64  `json:"medical_history_id"`
	ExamDate         string `json:"exam_date"`
	ExamType         string `json:"exam_type"`
	Findings         string `json:"findings"`
	VetNotes         string `json:"vet_notes"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Prescription represents a medication prescription for a pet
type PrescriptionRequest struct {
	MedicalHistoryID int64                  `json:"medical_history_id" binding:"required"`
	ExaminationID    int64                  `json:"examination_id" binding:"required"`
	PrescriptionDate string                 `json:"prescription_date" binding:"required,datetime=2006-01-02 15:04:05"`
	DoctorID         int64                  `json:"doctor_id" binding:"required"`
	Notes            string                 `json:"notes"`
	Medications      []PrescribedMedication `json:"medications" binding:"required"`
}

type PrescribedMedication struct {
	MedicineID   int64  `json:"medicine_id" binding:"required"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	Duration     string `json:"duration" binding:"required"`
	Instructions string `json:"instructions"`
}

type PrescriptionResponse struct {
	ID               int64                            `json:"id"`
	MedicalHistoryID int64                            `json:"medical_history_id"`
	ExaminationID    int64                            `json:"examination_id"`
	PrescriptionDate string                           `json:"prescription_date"`
	DoctorID         int64                            `json:"doctor_id"`
	DoctorName       string                           `json:"doctor_name"`
	Notes            string                           `json:"notes"`
	Medications      []PrescriptionMedicationResponse `json:"medications"`
	CreatedAt        string                           `json:"created_at"`
	UpdatedAt        string                           `json:"updated_at"`
}

type PrescriptionMedicationResponse struct {
	ID             int64  `json:"id"`
	PrescriptionID int64  `json:"prescription_id"`
	MedicineID     int64  `json:"medicine_id"`
	MedicineName   string `json:"medicine_name"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	Duration       string `json:"duration"`
	Instructions   string `json:"instructions"`
}

// TestResult represents test results for a pet
type TestResultRequest struct {
	MedicalHistoryID int64  `json:"medical_history_id" binding:"required"`
	ExaminationID    int64  `json:"examination_id" binding:"required"`
	TestType         string `json:"test_type" binding:"required"`
	TestDate         string `json:"test_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Results          string `json:"results" binding:"required"`
	Interpretation   string `json:"interpretation"`
	FileURL          string `json:"file_url"`
	DoctorID         int64  `json:"doctor_id" binding:"required"`
}

type TestResultResponse struct {
	ID               int64  `json:"id"`
	MedicalHistoryID int64  `json:"medical_history_id"`
	ExaminationID    int64  `json:"examination_id"`
	TestType         string `json:"test_type"`
	TestDate         string `json:"test_date"`
	Results          string `json:"results"`
	Interpretation   string `json:"interpretation"`
	FileURL          string `json:"file_url"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// MedicalHistorySummary represents a comprehensive view of medical history
type MedicalHistorySummary struct {
	MedicalRecord MedicalRecordResponse    `json:"medical_record"`
	Examinations  []ExaminationResponse    `json:"examinations"`
	Prescriptions []PrescriptionResponse   `json:"prescriptions"`
	TestResults   []TestResultResponse     `json:"test_results"`
	Conditions    []MedicalHistoryResponse `json:"conditions"`
	Allergies     []Allergy                `json:"allergies"`
}

type SoapNoteResponse struct {
	ID               int64  `json:"id"`
	PetID            int64  `json:"pet_id"`
	Subjective       string `json:"subjective"`
	Objective        string `json:"objective"`
	Assessment       string `json:"assessment"`
	Plan             string `json:"plan"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	ConsultationDate string `json:"consultation_date"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// AppointmentVisitSummaryResponse là cấu trúc trả về kết quả tổng hợp của một lần khám bệnh
type AppointmentVisitSummaryResponse struct {
	AppointmentID     int64                   `json:"appointment_id"`
	VisitDate         time.Time               `json:"visit_date"`
	PetInfo           PetBasicInfo            `json:"pet_info"`
	OwnerInfo         OwnerBasicInfo          `json:"owner_info"`
	DoctorInfo        DoctorBasicInfo         `json:"doctor_info"`
	ServiceInfo       ServiceBasicInfo        `json:"service_info"`
	AppointmentStatus string                  `json:"appointment_status"`
	Reason            string                  `json:"appointment_reason"`
	SOAPNote          *SOAPNoteInfo           `json:"soap_note"`
	Treatments        []TreatmentBasicInfo    `json:"treatments"`
	Prescriptions     []PrescriptionBasicInfo `json:"prescriptions"`
	TestResults       []TestResultBasicInfo   `json:"test_results"`
	NextAppointment   *NextAppointmentInfo    `json:"next_appointment"`
}

// PetBasicInfo chứa thông tin cơ bản về thú cưng
type PetBasicInfo struct {
	PetID  int64   `json:"pet_id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Breed  string  `json:"breed"`
	Age    int32   `json:"age,omitempty"`
	Gender string  `json:"gender,omitempty"`
	Weight float64 `json:"weight,omitempty"`
}

// OwnerBasicInfo chứa thông tin cơ bản về chủ sở hữu
type OwnerBasicInfo struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
}

// DoctorBasicInfo chứa thông tin cơ bản về bác sĩ
type DoctorBasicInfo struct {
	DoctorID int64  `json:"doctor_id"`
	Name     string `json:"name"`
}

// ServiceBasicInfo chứa thông tin cơ bản về dịch vụ
type ServiceBasicInfo struct {
	ServiceID int64   `json:"service_id"`
	Name      string  `json:"name"`
	Duration  int32   `json:"duration,omitempty"`
	Price     float64 `json:"price,omitempty"`
}

// SOAPNoteInfo chứa thông tin SOAP note
type SOAPNoteInfo struct {
	ID         int32     `json:"id,omitempty"`
	Subjective string    `json:"subjective,omitempty"`
	Objective  string    `json:"objective,omitempty"`
	Assessment string    `json:"assessment,omitempty"`
	Plan       string    `json:"plan,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// VitalSignsInfo chứa thông tin dấu hiệu sinh tồn
type VitalSignsInfo struct {
	Temperature     float64 `json:"temperature,omitempty"`
	HeartRate       int32   `json:"heart_rate,omitempty"`
	RespirationRate int32   `json:"respiration_rate,omitempty"`
	Weight          float64 `json:"weight,omitempty"`
}

// TreatmentBasicInfo chứa thông tin cơ bản về điều trị
type TreatmentBasicInfo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type,omitempty"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty"`
	Description string    `json:"description,omitempty"`
}

// PrescriptionBasicInfo chứa thông tin cơ bản về đơn thuốc
type PrescriptionBasicInfo struct {
	ID           int64     `json:"id"`
	MedicineName string    `json:"medicine_name"`
	Dosage       string    `json:"dosage,omitempty"`
	Frequency    string    `json:"frequency,omitempty"`
	Duration     string    `json:"duration,omitempty"`
	Quantity     int32     `json:"quantity,omitempty"`
	Instructions string    `json:"instructions,omitempty"`
	IssuedDate   time.Time `json:"issued_date,omitempty"`
}

// TestResultBasicInfo chứa thông tin cơ bản về kết quả xét nghiệm
type TestResultBasicInfo struct {
	ID          int64     `json:"id"`
	TestName    string    `json:"test_name"`
	Result      string    `json:"result,omitempty"`
	NormalRange string    `json:"normal_range,omitempty"`
	Status      string    `json:"status,omitempty"`
	TestDate    time.Time `json:"test_date,omitempty"`
}

// NextAppointmentInfo chứa thông tin về lịch hẹn tiếp theo
type NextAppointmentInfo struct {
	AppointmentID int64     `json:"appointment_id"`
	Date          time.Time `json:"date"`
	ServiceName   string    `json:"service_name,omitempty"`
	DoctorName    string    `json:"doctor_name,omitempty"`
}

package doctor

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type DoctorController struct {
	service DoctorServiceInterface
}

type DoctorService struct {
	storeDB db.Store
}

// API structure
type DoctorApi struct {
	controller DoctorControllerInterface
}

// Request and response structures

type EditDoctorProfileRequest struct {
	Specialization  string      `json:"specialization"`
	YearsOfExp      int32       `json:"years_of_experience"`
	Education       string      `json:"education"`
	Certificate     string      `json:"certificate_number"`
	Bio             string      `json:"bio"`
	Qualifications  []string    `json:"qualifications"`
	TrainingRecords []string    `json:"training_records"`
	AvailableHours  []WorkHours `json:"available_hours"`
	Department      string      `json:"department"`
	Position        string      `json:"position"`
}

type WorkHours struct {
	DayOfWeek   int       `json:"day_of_week"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsAvailable bool      `json:"is_available"`
}

type DoctorPerformance struct {
	DoctorID              int64   `json:"doctor_id"`
	ServiceCompletionRate float64 `json:"service_completion_rate"`
	CustomerSatisfaction  float64 `json:"customer_satisfaction"`
	RevenueGenerated      float64 `json:"revenue_generated"`
	AttendanceRate        float64 `json:"attendance_rate"`
	CompletedAppointments int     `json:"completed_appointments"`
	TotalAppointments     int     `json:"total_appointments"`
}

type PerformanceReview struct {
	ReviewID     int64     `json:"review_id"`
	DoctorID     int64     `json:"doctor_id"`
	ReviewerID   int64     `json:"reviewer_id"`
	ReviewDate   time.Time `json:"review_date"`
	Rating       float64   `json:"rating"`
	Comments     string    `json:"comments"`
	Goals        []string  `json:"goals"`
	Achievements []string  `json:"achievements"`
}

type loginDoctorRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=25"`
}

type loginDoctorResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	Doctor                DoctorDetail `json:"doctor"`
}

type DoctorDetail struct {
	DoctorID       int64  `json:"doctor_id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	DoctorName     string `json:"doctor_name"`
	PhoneNumber    string `json:"phone_number"`
	Address        string `json:"address"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Specialization string `json:"specialization"`
	YearsOfExp     int32  `json:"years_of_experience"`
	Education      string `json:"education"`
	Certificate    string `json:"certificate_number"`
	Bio            string `json:"bio"`
	DataImage      []byte `json:"data_image,omitempty"`
}

type Shift struct {
	ID               int64     `json:"id"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	AssignedPatients int32     `json:"assigned_patients"`
	CreatedAt        time.Time `json:"created_at"`
	DoctorID         int64     `json:"doctor_id"`
	DoctorName       string    `json:"doctor_name"`
}

type ShiftResponse struct {
	ID               int64  `json:"id"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	AssignedPatients int32  `json:"assigned_patients"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
}

type CreateShiftRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	DoctorID  int64     `json:"doctor_id"`
}

type LeaveRequest struct {
	ID          int64     `json:"id"`
	DoctorID    int64     `json:"doctor_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	LeaveType   string    `json:"leave_type"` // vacation, sick, personal, other
	Reason      string    `json:"reason"`
	Status      string    `json:"status"` // pending, approved, rejected
	ReviewedBy  int64     `json:"reviewed_by,omitempty"`
	ReviewNotes string    `json:"review_notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateLeaveRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	LeaveType string    `json:"leave_type" binding:"required"`
	Reason    string    `json:"reason"`
}

type UpdateLeaveRequest struct {
	Status      string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewNotes string `json:"review_notes"`
}

type AttendanceRecord struct {
	WorkDate              time.Time `json:"work_date"`
	TotalAppointments     int       `json:"total_appointments"`
	CompletedAppointments int       `json:"completed_appointments"`
	WorkHours             float64   `json:"work_hours"`
}

type WorkloadMetrics struct {
	TotalAppointments     int     `json:"total_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	AvgWorkHoursPerDay    float64 `json:"avg_work_hours_per_day"`
	TotalWorkDays         int     `json:"total_work_days"`
}

type TrainingRecord struct {
	ID                  int64     `json:"id"`
	DoctorID            int64     `json:"doctor_id"`
	TrainingName        string    `json:"training_name"`
	Provider            string    `json:"provider"`
	CompletionDate      time.Time `json:"completion_date"`
	ExpiryDate          time.Time `json:"expiry_date,omitempty"`
	CertificationNumber string    `json:"certification_number"`
	Description         string    `json:"description"`
	Status              string    `json:"status"` // active, expired, revoked
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateTrainingRequest struct {
	TrainingName        string    `json:"training_name" binding:"required"`
	Provider            string    `json:"provider" binding:"required"`
	CompletionDate      time.Time `json:"completion_date" binding:"required"`
	ExpiryDate          time.Time `json:"expiry_date"`
	CertificationNumber string    `json:"certification_number"`
	Description         string    `json:"description"`
}

type Qualification struct {
	ID                 int64     `json:"id"`
	DoctorID           int64     `json:"doctor_id"`
	QualificationType  string    `json:"qualification_type"` // degree, certification, license
	QualificationName  string    `json:"qualification_name"`
	Institution        string    `json:"institution"`
	IssueDate          time.Time `json:"issue_date"`
	RegistrationNumber string    `json:"registration_number"`
	ValidUntil         time.Time `json:"valid_until,omitempty"`
	DocumentURL        string    `json:"document_url,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateQualificationRequest struct {
	QualificationType  string    `json:"qualification_type" binding:"required"`
	QualificationName  string    `json:"qualification_name" binding:"required"`
	Institution        string    `json:"institution" binding:"required"`
	IssueDate          time.Time `json:"issue_date" binding:"required"`
	RegistrationNumber string    `json:"registration_number"`
	ValidUntil         time.Time `json:"valid_until"`
	DocumentURL        string    `json:"document_url"`
}

type ExpirationAlert struct {
	DoctorID            int64     `json:"doctor_id"`
	DoctorName          string    `json:"doctor_name"`
	QualificationType   string    `json:"qualification_type,omitempty"`
	QualificationName   string    `json:"qualification_name,omitempty"`
	QualificationExpiry time.Time `json:"qualification_expiry,omitempty"`
	TrainingName        string    `json:"training_name,omitempty"`
	TrainingExpiry      time.Time `json:"training_expiry,omitempty"`
}

type DoctorKPIs struct {
	DoctorID              int64   `json:"doctor_id"`
	DoctorName            string  `json:"doctor_name"`
	TotalAppointments     int     `json:"total_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	AverageSatisfaction   float64 `json:"avg_satisfaction"`
	AttendanceRate        float64 `json:"attendance_rate"`
	RevenueGenerated      float64 `json:"revenue_generated"`
}

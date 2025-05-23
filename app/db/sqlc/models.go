// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type MovementTypeEnum string

const (
	MovementTypeEnumImport MovementTypeEnum = "import"
	MovementTypeEnumExport MovementTypeEnum = "export"
)

func (e *MovementTypeEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MovementTypeEnum(s)
	case string:
		*e = MovementTypeEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for MovementTypeEnum: %T", src)
	}
	return nil
}

type NullMovementTypeEnum struct {
	MovementTypeEnum MovementTypeEnum `json:"movement_type_enum"`
	Valid            bool             `json:"valid"` // Valid is true if MovementTypeEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMovementTypeEnum) Scan(value interface{}) error {
	if value == nil {
		ns.MovementTypeEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MovementTypeEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMovementTypeEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MovementTypeEnum), nil
}

type Appointment struct {
	AppointmentID     int64            `json:"appointment_id"`
	Petid             pgtype.Int8      `json:"petid"`
	Username          pgtype.Text      `json:"username"`
	DoctorID          pgtype.Int8      `json:"doctor_id"`
	ServiceID         pgtype.Int8      `json:"service_id"`
	Date              pgtype.Timestamp `json:"date"`
	ReminderSend      pgtype.Bool      `json:"reminder_send"`
	TimeSlotID        pgtype.Int8      `json:"time_slot_id"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	StateID           pgtype.Int4      `json:"state_id"`
	AppointmentReason pgtype.Text      `json:"appointment_reason"`
	Priority          pgtype.Text      `json:"priority"`
	ArrivalTime       pgtype.Timestamp `json:"arrival_time"`
	RoomID            pgtype.Int8      `json:"room_id"`
	ConfirmationSent  pgtype.Bool      `json:"confirmation_sent"`
	Notes             pgtype.Text      `json:"notes"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type Cart struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type CartItem struct {
	ID         int64         `json:"id"`
	CartID     int64         `json:"cart_id"`
	ProductID  int64         `json:"product_id"`
	Quantity   pgtype.Int4   `json:"quantity"`
	UnitPrice  float64       `json:"unit_price"`
	TotalPrice pgtype.Float8 `json:"total_price"`
}

type Consultation struct {
	ID            int32            `json:"id"`
	AppointmentID pgtype.Int8      `json:"appointment_id"`
	Subjective    []byte           `json:"subjective"`
	Objective     []byte           `json:"objective"`
	Assessment    []byte           `json:"assessment"`
	Plan          pgtype.Int8      `json:"plan"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}

type DeviceToken struct {
	ID         int64            `json:"id"`
	Username   string           `json:"username"`
	Token      string           `json:"token"`
	DeviceType pgtype.Text      `json:"device_type"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	LastUsedAt pgtype.Timestamp `json:"last_used_at"`
	ExpiredAt  pgtype.Timestamp `json:"expired_at"`
}

type Disease struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description pgtype.Text        `json:"description"`
	Symptoms    []byte             `json:"symptoms"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type Doctor struct {
	ID                int64         `json:"id"`
	UserID            int64         `json:"user_id"`
	Specialization    pgtype.Text   `json:"specialization"`
	YearsOfExperience pgtype.Int4   `json:"years_of_experience"`
	Education         pgtype.Text   `json:"education"`
	CertificateNumber pgtype.Text   `json:"certificate_number"`
	Bio               pgtype.Text   `json:"bio"`
	ConsultationFee   pgtype.Float8 `json:"consultation_fee"`
}

type Examination struct {
	ID               int64            `json:"id"`
	MedicalHistoryID int64            `json:"medical_history_id"`
	ExamDate         pgtype.Timestamp `json:"exam_date"`
	ExamType         string           `json:"exam_type"`
	Findings         string           `json:"findings"`
	VetNotes         pgtype.Text      `json:"vet_notes"`
	DoctorID         int64            `json:"doctor_id"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	UpdatedAt        pgtype.Timestamp `json:"updated_at"`
}

type File struct {
	ID         int64            `json:"id"`
	FileName   string           `json:"file_name"`
	FilePath   string           `json:"file_path"`
	FileSize   int64            `json:"file_size"`
	FileType   string           `json:"file_type"`
	UploadedAt pgtype.Timestamp `json:"uploaded_at"`
	PetID      pgtype.Int8      `json:"pet_id"`
}

type Invoice struct {
	ID            int32            `json:"id"`
	InvoiceNumber string           `json:"invoice_number"`
	Amount        float64          `json:"amount"`
	Date          pgtype.Date      `json:"date"`
	DueDate       pgtype.Date      `json:"due_date"`
	Status        string           `json:"status"`
	Description   pgtype.Text      `json:"description"`
	CustomerName  pgtype.Text      `json:"customer_name"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	Type          pgtype.Text      `json:"type"`
	AppointmentID pgtype.Int8      `json:"appointment_id"`
	TestOrderID   pgtype.Int8      `json:"test_order_id"`
	OrderID       pgtype.Int8      `json:"order_id"`
}

type InvoiceItem struct {
	ID        int32   `json:"id"`
	InvoiceID int32   `json:"invoice_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int32   `json:"quantity"`
}

type MedicalHistory struct {
	ID              int64            `json:"id"`
	MedicalRecordID pgtype.Int8      `json:"medical_record_id"`
	Condition       pgtype.Text      `json:"condition"`
	DiagnosisDate   pgtype.Timestamp `json:"diagnosis_date"`
	Notes           pgtype.Text      `json:"notes"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}

type MedicalRecord struct {
	ID        int64            `json:"id"`
	PetID     pgtype.Int8      `json:"pet_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Medicine struct {
	ID             int64              `json:"id"`
	Name           string             `json:"name"`
	Description    pgtype.Text        `json:"description"`
	Usage          pgtype.Text        `json:"usage"`
	Dosage         pgtype.Text        `json:"dosage"`
	Frequency      pgtype.Text        `json:"frequency"`
	Duration       pgtype.Text        `json:"duration"`
	SideEffects    pgtype.Text        `json:"side_effects"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
	ExpirationDate pgtype.Date        `json:"expiration_date"`
	Quantity       pgtype.Int8        `json:"quantity"`
	UnitPrice      pgtype.Float8      `json:"unit_price"`
	ReorderLevel   pgtype.Int8        `json:"reorder_level"`
	SupplierID     pgtype.Int8        `json:"supplier_id"`
}

type MedicineSupplier struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Email       pgtype.Text        `json:"email"`
	Phone       pgtype.Text        `json:"phone"`
	Address     pgtype.Text        `json:"address"`
	ContactName pgtype.Text        `json:"contact_name"`
	Notes       pgtype.Text        `json:"notes"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type MedicineTransaction struct {
	ID              int64              `json:"id"`
	MedicineID      int64              `json:"medicine_id"`
	Quantity        int64              `json:"quantity"`
	TransactionType string             `json:"transaction_type"`
	UnitPrice       pgtype.Float8      `json:"unit_price"`
	TotalAmount     pgtype.Float8      `json:"total_amount"`
	TransactionDate pgtype.Timestamptz `json:"transaction_date"`
	SupplierID      pgtype.Int8        `json:"supplier_id"`
	ExpirationDate  pgtype.Date        `json:"expiration_date"`
	Notes           pgtype.Text        `json:"notes"`
	PrescriptionID  pgtype.Int8        `json:"prescription_id"`
	AppointmentID   pgtype.Int8        `json:"appointment_id"`
	CreatedBy       pgtype.Text        `json:"created_by"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
}

type Notification struct {
	ID          int64            `json:"id"`
	Username    string           `json:"username"`
	Title       string           `json:"title"`
	Content     pgtype.Text      `json:"content"`
	IsRead      pgtype.Bool      `json:"is_read"`
	RelatedID   pgtype.Int4      `json:"related_id"`
	RelatedType pgtype.Text      `json:"related_type"`
	Datetime    pgtype.Timestamp `json:"datetime"`
	NotifyType  pgtype.Text      `json:"notify_type"`
}

type Order struct {
	ID              int64            `json:"id"`
	UserID          int64            `json:"user_id"`
	OrderDate       pgtype.Timestamp `json:"order_date"`
	TotalAmount     float64          `json:"total_amount"`
	PaymentStatus   pgtype.Text      `json:"payment_status"`
	CartItems       []byte           `json:"cart_items"`
	ShippingAddress pgtype.Text      `json:"shipping_address"`
	Notes           pgtype.Text      `json:"notes"`
}

type OrderedTest struct {
	ID              int32              `json:"id"`
	OrderID         pgtype.Int4        `json:"order_id"`
	TestID          pgtype.Int4        `json:"test_id"`
	PriceAtOrder    float64            `json:"price_at_order"`
	Status          pgtype.Text        `json:"status"`
	Results         pgtype.Text        `json:"results"`
	ResultsDate     pgtype.Timestamptz `json:"results_date"`
	TechnicianNotes pgtype.Text        `json:"technician_notes"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	DoctorID        pgtype.Int4        `json:"doctor_id"`
	NextDueDate     pgtype.Date        `json:"next_due_date"`
}

type Payment struct {
	ID             int32              `json:"id"`
	Amount         float64            `json:"amount"`
	PaymentMethod  string             `json:"payment_method"`
	PaymentStatus  string             `json:"payment_status"`
	OrderID        pgtype.Int4        `json:"order_id"`
	TestOrderID    pgtype.Int4        `json:"test_order_id"`
	TransactionID  pgtype.Text        `json:"transaction_id"`
	PaymentDetails []byte             `json:"payment_details"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
	AppointmentID  pgtype.Int8        `json:"appointment_id"`
}

type Pet struct {
	Petid           int64         `json:"petid"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	Weight          pgtype.Float8 `json:"weight"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	Username        string        `json:"username"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
	LastCheckupDate pgtype.Date   `json:"last_checkup_date"`
	IsActive        pgtype.Bool   `json:"is_active"`
	DataImage       []byte        `json:"data_image"`
	OriginalImage   pgtype.Text   `json:"original_image"`
}

type PetLog struct {
	LogID    int64            `json:"log_id"`
	Petid    int64            `json:"petid"`
	Datetime pgtype.Timestamp `json:"datetime"`
	Title    pgtype.Text      `json:"title"`
	Notes    pgtype.Text      `json:"notes"`
}

type PetSchedule struct {
	ID               int64            `json:"id"`
	PetID            pgtype.Int8      `json:"pet_id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Bool      `json:"end_type"`
	EndDate          pgtype.Date      `json:"end_date"`
	Notes            pgtype.Text      `json:"notes"`
	IsActive         pgtype.Bool      `json:"is_active"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	Removedat        pgtype.Timestamp `json:"removedat"`
}

type PetTreatment struct {
	ID          int64              `json:"id"`
	PetID       pgtype.Int8        `json:"pet_id"`
	StartDate   pgtype.Date        `json:"start_date"`
	EndDate     pgtype.Date        `json:"end_date"`
	Status      pgtype.Text        `json:"status"`
	Description pgtype.Text        `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	DoctorID    pgtype.Int4        `json:"doctor_id"`
	Name        pgtype.Text        `json:"name"`
	Type        pgtype.Text        `json:"type"`
	Diseases    pgtype.Text        `json:"diseases"`
}

type PetWeightHistory struct {
	ID         int64            `json:"id"`
	PetID      int64            `json:"pet_id"`
	WeightKg   float64          `json:"weight_kg"`
	RecordedAt pgtype.Timestamp `json:"recorded_at"`
	Notes      pgtype.Text      `json:"notes"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
}

type PhaseMedicine struct {
	PhaseID    int64              `json:"phase_id"`
	MedicineID int64              `json:"medicine_id"`
	Dosage     pgtype.Text        `json:"dosage"`
	Frequency  pgtype.Text        `json:"frequency"`
	Duration   pgtype.Text        `json:"duration"`
	Notes      pgtype.Text        `json:"notes"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	Quantity   pgtype.Int4        `json:"quantity"`
	IsReceived pgtype.Bool        `json:"is_received"`
}

type Prescription struct {
	ID               int64            `json:"id"`
	MedicalHistoryID int64            `json:"medical_history_id"`
	ExaminationID    int64            `json:"examination_id"`
	PrescriptionDate pgtype.Timestamp `json:"prescription_date"`
	DoctorID         int64            `json:"doctor_id"`
	Notes            pgtype.Text      `json:"notes"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	UpdatedAt        pgtype.Timestamp `json:"updated_at"`
}

type Product struct {
	ProductID      int64            `json:"product_id"`
	Name           string           `json:"name"`
	Description    pgtype.Text      `json:"description"`
	Price          float64          `json:"price"`
	StockQuantity  pgtype.Int4      `json:"stock_quantity"`
	Category       pgtype.Text      `json:"category"`
	DataImage      []byte           `json:"data_image"`
	OriginalImage  pgtype.Text      `json:"original_image"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	IsAvailable    pgtype.Bool      `json:"is_available"`
	RemovedAt      pgtype.Timestamp `json:"removed_at"`
	ExpirationDate pgtype.Date      `json:"expiration_date"`
	ReorderLevel   pgtype.Int4      `json:"reorder_level"`
	SupplierID     pgtype.Int8      `json:"supplier_id"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
}

type ProductStockMovement struct {
	MovementID   int64            `json:"movement_id"`
	ProductID    int64            `json:"product_id"`
	MovementType MovementTypeEnum `json:"movement_type"`
	Quantity     int32            `json:"quantity"`
	Reason       pgtype.Text      `json:"reason"`
	MovementDate time.Time        `json:"movement_date"`
	Price        pgtype.Numeric   `json:"price"`
}

type Room struct {
	ID                   int64            `json:"id"`
	Name                 string           `json:"name"`
	Type                 string           `json:"type"`
	Status               pgtype.Text      `json:"status"`
	CurrentAppointmentID pgtype.Int8      `json:"current_appointment_id"`
	AvailableAt          pgtype.Timestamp `json:"available_at"`
}

type Service struct {
	ID          int64            `json:"id"`
	Name        pgtype.Text      `json:"name"`
	Description pgtype.Text      `json:"description"`
	Duration    pgtype.Int2      `json:"duration"`
	Cost        pgtype.Float8    `json:"cost"`
	Category    pgtype.Text      `json:"category"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	Priority    pgtype.Int2      `json:"priority"`
	RemovedAt   pgtype.Timestamp `json:"removed_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type Shift struct {
	ID          int64            `json:"id"`
	DoctorID    int64            `json:"doctor_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	MaxPatients pgtype.Int4      `json:"max_patients"`
	Date        pgtype.Date      `json:"date"`
}

type SmtpConfig struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	SmtpHost  string    `json:"smtp_host"`
	SmtpPort  string    `json:"smtp_port"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type State struct {
	ID        int64            `json:"id"`
	State     string           `json:"state"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Test struct {
	ID             int32              `json:"id"`
	TestID         string             `json:"test_id"`
	CategoryID     pgtype.Text        `json:"category_id"`
	Name           string             `json:"name"`
	Description    pgtype.Text        `json:"description"`
	Price          float64            `json:"price"`
	TurnaroundTime string             `json:"turnaround_time"`
	IsActive       pgtype.Bool        `json:"is_active"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
	Quantity       pgtype.Int4        `json:"quantity"`
	ExpirationDate pgtype.Date        `json:"expiration_date"`
	BatchNumber    pgtype.Text        `json:"batch_number"`
	SupplierID     pgtype.Int4        `json:"supplier_id"`
	Type           pgtype.Text        `json:"type"`
}

type TestCategory struct {
	ID          int32              `json:"id"`
	CategoryID  string             `json:"category_id"`
	Name        string             `json:"name"`
	Description pgtype.Text        `json:"description"`
	IconName    pgtype.Text        `json:"icon_name"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type TestOrder struct {
	OrderID       int32              `json:"order_id"`
	AppointmentID pgtype.Int4        `json:"appointment_id"`
	OrderDate     pgtype.Timestamptz `json:"order_date"`
	Status        pgtype.Text        `json:"status"`
	TotalAmount   pgtype.Float8      `json:"total_amount"`
	Notes         pgtype.Text        `json:"notes"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
}

type TestResult struct {
	ID               int64            `json:"id"`
	MedicalHistoryID int64            `json:"medical_history_id"`
	ExaminationID    int64            `json:"examination_id"`
	TestType         string           `json:"test_type"`
	TestDate         pgtype.Timestamp `json:"test_date"`
	Results          string           `json:"results"`
	Interpretation   pgtype.Text      `json:"interpretation"`
	FileUrl          pgtype.Text      `json:"file_url"`
	DoctorID         int64            `json:"doctor_id"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	UpdatedAt        pgtype.Timestamp `json:"updated_at"`
}

type TimeSlot struct {
	ID             int64            `json:"id"`
	DoctorID       int32            `json:"doctor_id"`
	Date           pgtype.Date      `json:"date"`
	StartTime      pgtype.Time      `json:"start_time"`
	EndTime        pgtype.Time      `json:"end_time"`
	MaxPatients    pgtype.Int4      `json:"max_patients"`
	BookedPatients pgtype.Int4      `json:"booked_patients"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
	ShiftID        int64            `json:"shift_id"`
}

type TreatmentPhase struct {
	ID          int64              `json:"id"`
	TreatmentID pgtype.Int8        `json:"treatment_id"`
	PhaseName   pgtype.Text        `json:"phase_name"`
	Description pgtype.Text        `json:"description"`
	Status      pgtype.Text        `json:"status"`
	StartDate   pgtype.Date        `json:"start_date"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	IsLocked    pgtype.Bool        `json:"is_locked"`
}

type User struct {
	ID              int64            `json:"id"`
	Username        string           `json:"username"`
	HashedPassword  string           `json:"hashed_password"`
	FullName        string           `json:"full_name"`
	Email           string           `json:"email"`
	PhoneNumber     pgtype.Text      `json:"phone_number"`
	Address         pgtype.Text      `json:"address"`
	DataImage       []byte           `json:"data_image"`
	OriginalImage   pgtype.Text      `json:"original_image"`
	Role            pgtype.Text      `json:"role"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	IsVerifiedEmail pgtype.Bool      `json:"is_verified_email"`
	RemovedAt       pgtype.Timestamp `json:"removed_at"`
	Status          pgtype.Text      `json:"status"`
}

type Vaccination struct {
	Vaccinationid    int64            `json:"vaccinationid"`
	Petid            pgtype.Int8      `json:"petid"`
	Vaccinename      string           `json:"vaccinename"`
	Dateadministered pgtype.Timestamp `json:"dateadministered"`
	Nextduedate      pgtype.Timestamp `json:"nextduedate"`
	Vaccineprovider  pgtype.Text      `json:"vaccineprovider"`
	Batchnumber      pgtype.Text      `json:"batchnumber"`
	Notes            pgtype.Text      `json:"notes"`
}

type VerifyEmail struct {
	ID         int64            `json:"id"`
	Username   string           `json:"username"`
	Email      string           `json:"email"`
	SecretCode int64            `json:"secret_code"`
	IsUsed     bool             `json:"is_used"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	ExpiredAt  pgtype.Timestamp `json:"expired_at"`
}

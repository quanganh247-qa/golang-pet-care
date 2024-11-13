// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Activitylog struct {
	Logid        int64            `json:"logid"`
	Petid        pgtype.Int8      `json:"petid"`
	Activitytype string           `json:"activitytype"`
	Starttime    pgtype.Timestamp `json:"starttime"`
	Duration     pgtype.Interval  `json:"duration"`
	Notes        pgtype.Text      `json:"notes"`
}

type Appointment struct {
	AppointmentID int64            `json:"appointment_id"`
	Petid         pgtype.Int8      `json:"petid"`
	DoctorID      pgtype.Int8      `json:"doctor_id"`
	ServiceID     pgtype.Int8      `json:"service_id"`
	Date          pgtype.Timestamp `json:"date"`
	Status        pgtype.Text      `json:"status"`
	Notes         pgtype.Text      `json:"notes"`
	ReminderSend  pgtype.Bool      `json:"reminder_send"`
	TimeSlotID    pgtype.Int8      `json:"time_slot_id"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}

type Checkout struct {
	CheckoutID    int64            `json:"checkout_id"`
	Petid         pgtype.Int8      `json:"petid"`
	DoctorID      pgtype.Int8      `json:"doctor_id"`
	Date          pgtype.Timestamp `json:"date"`
	TotalTmount   float64          `json:"total_tmount"`
	PaymentStatus pgtype.Text      `json:"payment_status"`
	PaymentMethod pgtype.Text      `json:"payment_method"`
	Notes         pgtype.Text      `json:"notes"`
}

type Checkoutservice struct {
	CheckoutserviceID int64         `json:"checkoutservice_id"`
	Checkoutid        pgtype.Int8   `json:"checkoutid"`
	Serviceid         pgtype.Int8   `json:"serviceid"`
	Quantity          pgtype.Int4   `json:"quantity"`
	Unitprice         pgtype.Float8 `json:"unitprice"`
	Subtotal          pgtype.Float8 `json:"subtotal"`
}

type Devicetoken struct {
	ID         int64            `json:"id"`
	Username   string           `json:"username"`
	Token      string           `json:"token"`
	DeviceType pgtype.Text      `json:"device_type"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	LastUsedAt pgtype.Timestamp `json:"last_used_at"`
	ExpiredAt  pgtype.Timestamp `json:"expired_at"`
}

type Doctor struct {
	ID                int64          `json:"id"`
	UserID            int64          `json:"user_id"`
	Specialization    pgtype.Text    `json:"specialization"`
	YearsOfExperience pgtype.Int4    `json:"years_of_experience"`
	Education         pgtype.Text    `json:"education"`
	CertificateNumber pgtype.Text    `json:"certificate_number"`
	Bio               pgtype.Text    `json:"bio"`
	ConsultationFee   pgtype.Numeric `json:"consultation_fee"`
}

type Doctorschedule struct {
	ID              int64            `json:"id"`
	DoctorID        int64            `json:"doctor_id"`
	DayOfWeek       pgtype.Int4      `json:"day_of_week"`
	StartTime       pgtype.Timestamp `json:"start_time"`
	EndTime         pgtype.Timestamp `json:"end_time"`
	IsActive        pgtype.Bool      `json:"is_active"`
	MaxAppointments pgtype.Int4      `json:"max_appointments"`
}

type Feedingschedule struct {
	Feedingscheduleid int64            `json:"feedingscheduleid"`
	Petid             pgtype.Int8      `json:"petid"`
	Mealtime          pgtype.Time      `json:"mealtime"`
	Foodtype          string           `json:"foodtype"`
	Quantity          float64          `json:"quantity"`
	Frequency         string           `json:"frequency"`
	Lastfed           pgtype.Timestamp `json:"lastfed"`
	Notes             pgtype.Text      `json:"notes"`
	Isactive          pgtype.Bool      `json:"isactive"`
}

type Medication struct {
	MedicationID   int64            `json:"medication_id"`
	PetID          int64            `json:"pet_id"`
	MedicationName string           `json:"medication_name"`
	Dosage         string           `json:"dosage"`
	Frequency      string           `json:"frequency"`
	StartDate      pgtype.Timestamp `json:"start_date"`
	EndDate        pgtype.Timestamp `json:"end_date"`
	Notes          pgtype.Text      `json:"notes"`
}

type Notification struct {
	Notificationid   int64            `json:"notificationid"`
	Petid            pgtype.Int8      `json:"petid"`
	Title            string           `json:"title"`
	Body             pgtype.Text      `json:"body"`
	Duedate          pgtype.Timestamp `json:"duedate"`
	Repeatinterval   pgtype.Text      `json:"repeatinterval"`
	Iscompleted      pgtype.Bool      `json:"iscompleted"`
	Notificationsent pgtype.Bool      `json:"notificationsent"`
}

type Notificationhistory struct {
	ID             int64              `json:"id"`
	NotificationID int64              `json:"notification_id"`
	UserID         int64              `json:"user_id"`
	DeviceTokenID  pgtype.Int8        `json:"device_token_id"`
	Title          string             `json:"title"`
	Body           pgtype.Text        `json:"body"`
	Data           []byte             `json:"data"`
	SentAt         pgtype.Timestamptz `json:"sent_at"`
	DeliveredAt    pgtype.Timestamptz `json:"delivered_at"`
	OpenedAt       pgtype.Timestamptz `json:"opened_at"`
	ErrorMessage   pgtype.Text        `json:"error_message"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
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
	OriginalImage   string        `json:"original_image"`
}

type Petservicelocation struct {
	Locationid    int64          `json:"locationid"`
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	Latitude      pgtype.Numeric `json:"latitude"`
	Longitude     pgtype.Numeric `json:"longitude"`
	Contactnumber pgtype.Text    `json:"contactnumber"`
	Businesshours pgtype.Text    `json:"businesshours"`
	Servicetypes  []string       `json:"servicetypes"`
	Rating        pgtype.Numeric `json:"rating"`
	Isverified    pgtype.Bool    `json:"isverified"`
}

type Service struct {
	Serviceid   int64            `json:"serviceid"`
	Typeid      pgtype.Int8      `json:"typeid"`
	Name        string           `json:"name"`
	Price       pgtype.Float8    `json:"price"`
	Duration    pgtype.Interval  `json:"duration"`
	Description pgtype.Text      `json:"description"`
	Isavailable pgtype.Bool      `json:"isavailable"`
	RemovedAt   pgtype.Timestamp `json:"removed_at"`
}

type Servicetype struct {
	Typeid          int64       `json:"typeid"`
	Servicetypename string      `json:"servicetypename"`
	Description     pgtype.Text `json:"description"`
	Iconurl         pgtype.Text `json:"iconurl"`
}

type Timeslot struct {
	ID        int64            `json:"id"`
	DoctorID  int64            `json:"doctor_id"`
	StartTime pgtype.Timestamp `json:"start_time"`
	EndTime   pgtype.Timestamp `json:"end_time"`
	IsActive  pgtype.Bool      `json:"is_active"`
	Day       pgtype.Date      `json:"day"`
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
	OriginalImage   string           `json:"original_image"`
	Role            pgtype.Text      `json:"role"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	IsVerifiedEmail pgtype.Bool      `json:"is_verified_email"`
	RemovedAt       pgtype.Timestamp `json:"removed_at"`
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
	SecretCode string           `json:"secret_code"`
	IsUsed     bool             `json:"is_used"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	ExpiredAt  pgtype.Timestamp `json:"expired_at"`
}

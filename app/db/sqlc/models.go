// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"time"

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
	Appointmentid int64            `json:"appointmentid"`
	Petid         pgtype.Int8      `json:"petid"`
	Doctorid      pgtype.Int8      `json:"doctorid"`
	Serviceid     pgtype.Int8      `json:"serviceid"`
	Date          pgtype.Timestamp `json:"date"`
	Status        pgtype.Text      `json:"status"`
	Notes         pgtype.Text      `json:"notes"`
	Remindersent  pgtype.Bool      `json:"remindersent"`
	TimeSlotID    pgtype.Int8      `json:"time_slot_id"`
}

type Checkout struct {
	Checkoutid    int64            `json:"checkoutid"`
	Petid         pgtype.Int8      `json:"petid"`
	Doctorid      pgtype.Int8      `json:"doctorid"`
	Date          pgtype.Timestamp `json:"date"`
	TotalAmount   float64          `json:"total_amount"`
	Paymentstatus pgtype.Text      `json:"paymentstatus"`
	Paymentmethod pgtype.Text      `json:"paymentmethod"`
	Note          pgtype.Text      `json:"note"`
}

type Checkoutservice struct {
	CheckoutserviceID int64         `json:"checkoutservice_id"`
	Checkoutid        pgtype.Int8   `json:"checkoutid"`
	Serviceid         pgtype.Int8   `json:"serviceid"`
	Quantity          pgtype.Int4   `json:"quantity"`
	Unitprice         pgtype.Float8 `json:"unitprice"`
	Subtotal          pgtype.Float8 `json:"subtotal"`
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
	ID              int64       `json:"id"`
	DoctorID        int64       `json:"doctor_id"`
	DayOfWeek       pgtype.Int4 `json:"day_of_week"`
	StartTime       pgtype.Time `json:"start_time"`
	EndTime         pgtype.Time `json:"end_time"`
	IsActive        pgtype.Bool `json:"is_active"`
	MaxAppointments pgtype.Int4 `json:"max_appointments"`
}

type Doctortimeoff struct {
	ID            int64            `json:"id"`
	DoctorID      int64            `json:"doctor_id"`
	StartDatetime pgtype.Timestamp `json:"start_datetime"`
	EndDatetime   pgtype.Timestamp `json:"end_datetime"`
	Reason        pgtype.Text      `json:"reason"`
}

type Feedingschedule struct {
	Feedingscheduleid int64            `json:"feedingscheduleid"`
	Petid             pgtype.Int8      `json:"petid"`
	Mealtime          pgtype.Time      `json:"mealtime"`
	Foodtype          string           `json:"foodtype"`
	Quantity          pgtype.Numeric   `json:"quantity"`
	Frequency         string           `json:"frequency"`
	Lastfed           pgtype.Timestamp `json:"lastfed"`
	Notes             pgtype.Text      `json:"notes"`
	Isactive          pgtype.Bool      `json:"isactive"`
}

type Pet struct {
	Petid           int64         `json:"petid"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	Profileimage    pgtype.Text   `json:"profileimage"`
	Weight          pgtype.Float8 `json:"weight"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	Username        string        `json:"username"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
	LastCheckupDate pgtype.Date   `json:"last_checkup_date"`
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

type Reminder struct {
	Reminderid       int64            `json:"reminderid"`
	Petid            pgtype.Int8      `json:"petid"`
	Title            string           `json:"title"`
	Description      pgtype.Text      `json:"description"`
	Duedate          pgtype.Timestamp `json:"duedate"`
	Repeatinterval   pgtype.Text      `json:"repeatinterval"`
	Iscompleted      pgtype.Bool      `json:"iscompleted"`
	Notificationsent pgtype.Bool      `json:"notificationsent"`
}

type Service struct {
	Serviceid   int64           `json:"serviceid"`
	Typeid      pgtype.Int8     `json:"typeid"`
	Name        string          `json:"name"`
	Price       interface{}     `json:"price"`
	Duration    pgtype.Interval `json:"duration"`
	Description pgtype.Text     `json:"description"`
	Isavailable pgtype.Bool     `json:"isavailable"`
}

type Servicetype struct {
	Typeid          int64       `json:"typeid"`
	Servicetypename string      `json:"servicetypename"`
	Description     pgtype.Text `json:"description"`
	Iconurl         pgtype.Text `json:"iconurl"`
}

type Timeslot struct {
	ID              int64       `json:"id"`
	StartTime       pgtype.Time `json:"start_time"`
	EndTime         pgtype.Time `json:"end_time"`
	DurationMinutes int32       `json:"duration_minutes"`
}

type User struct {
	ID              int64              `json:"id"`
	Username        string             `json:"username"`
	HashedPassword  string             `json:"hashed_password"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	PhoneNumber     pgtype.Text        `json:"phone_number"`
	Address         pgtype.Text        `json:"address"`
	Avatar          pgtype.Text        `json:"avatar"`
	Role            pgtype.Text        `json:"role"`
	CreatedAt       time.Time          `json:"created_at"`
	IsVerifiedEmail pgtype.Bool        `json:"is_verified_email"`
	RemovedAt       pgtype.Timestamptz `json:"removed_at"`
}

type Vaccination struct {
	Vaccinationid    int64       `json:"vaccinationid"`
	Petid            pgtype.Int8 `json:"petid"`
	Vaccinename      string      `json:"vaccinename"`
	Dateadministered pgtype.Date `json:"dateadministered"`
	Nextduedate      pgtype.Date `json:"nextduedate"`
	Vaccineprovider  pgtype.Text `json:"vaccineprovider"`
	Batchnumber      pgtype.Text `json:"batchnumber"`
	Notes            pgtype.Text `json:"notes"`
}

type VerifyEmail struct {
	ID         int64              `json:"id"`
	Username   string             `json:"username"`
	Email      string             `json:"email"`
	SecretCode string             `json:"secret_code"`
	IsUsed     bool               `json:"is_used"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	ExpiredAt  pgtype.Timestamptz `json:"expired_at"`
}

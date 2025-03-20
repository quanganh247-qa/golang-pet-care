<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	DoctorName     string `json:"doctor_name"`
=======
	FullName       string `json:"full_name"`
>>>>>>> 4ccd381 (Update appointment flow)
=======
	DoctorName     string `json:"doctor_name"`
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
	FullName       string `json:"full_name"`
>>>>>>> 4ccd381 (Update appointment flow)
=======
	DoctorName     string `json:"doctor_name"`
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
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
}

type CreateShiftRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	DoctorID  int64     `json:"doctor_id"`
}
<<<<<<< HEAD
<<<<<<< HEAD
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)

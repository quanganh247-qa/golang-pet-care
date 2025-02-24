package clinic_reporting

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ClinicReportingService struct {
	storeDB db.Store
}

// ClinicAnalyticsResponse holds the aggregated analytics for the entire clinic.
type ClinicAnalyticsResponse struct {
	TotalAppointments    int `json:"total_appointments"`
	UpcomingAppointments int `json:"upcoming_appointments"`
	PastAppointments     int `json:"past_appointments"`
	TotalPets            int `json:"total_pets"`
	TotalDoctors         int `json:"total_doctors"`
}

type GenerateSOAPNoteRequest struct {
	Text string `json:"text"`
}

type Note struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	PatientID  uint           `json:"patient_id"`
	Subjective SOAPSubjective `json:"subjective"`
	Objective  SOAPObjective  `json:"objective"`
	Assessment SOAPAssessment `json:"assessment"`
	Plan       SOAPPlan       `json:"plan"`
	VoiceData  []byte         `json:"voice_data"`
	CreatedAt  time.Time      `json:"created_at"`
}

type SOAPSubjective struct {
	PatientInfo    string   `json:"patient_info"`
	ChiefComplaint string   `json:"chief_complaint"`
	Duration       string   `json:"duration"`
	History        string   `json:"history"`
	Symptoms       []string `json:"symptoms"`
	Behavior       struct {
		ActivityLevel string `json:"activity_level"`
		Appetite      string `json:"appetite"`
	} `json:"behavior"`
}

type SOAPObjective struct {
	VitalSigns struct {
		Temperature      string `json:"temperature"`
		GeneralCondition string `json:"general_condition"`
		HeartRate        string `json:"heart_rate"`
		Weight           string `json:"weight"`
		RespiratoryRate  string `json:"respiratory_rate"`
	} `json:"vital_signs"`
	ExaminationFindings []string `json:"examination_findings"`
}

type SOAPAssessment struct {
	PrimaryDiagnosis string   `json:"primary_diagnosis"`
	Differentials    []string `json:"differentials"`
}

type SOAPPlan struct {
	ImmediateTreatment []string `json:"immediate_treatment"`
	Medications        []string `json:"medications"`
	FollowUp           string   `json:"follow_up"`
	ClientEducation    []string `json:"client_education"`
}

type CreateSoapNoteRequest struct {
	AppointmentID uint `json:"appointment_id"`
	Note          string `json:"note"`
}

type UpdateSoapNoteRequest struct {
	ID            int    `json:"id"`
	AppointmentID int    `json:"appointment_id"`
	Note          string `json:"note"`
}

type SoapNoteResponse struct {
	ID            int    `json:"id"`
	AppointmentID int    `json:"appointment_id"`
	Note          string `json:"note"`
}

package clinic_reporting

// ClinicAnalyticsResponse holds the aggregated analytics for the entire clinic.
type ClinicAnalyticsResponse struct {
	TotalAppointments    int `json:"total_appointments"`
	UpcomingAppointments int `json:"upcoming_appointments"`
	PastAppointments     int `json:"past_appointments"`
	TotalPets            int `json:"total_pets"`
	TotalDoctors         int `json:"total_doctors"`
}

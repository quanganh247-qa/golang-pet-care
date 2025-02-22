package clinic_reporting

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// ClinicReportingServiceInterface defines the interface for the analytics service.
type ClinicReportingServiceInterface interface {
	GetClinicAnalytics(ctx *gin.Context) (*ClinicAnalyticsResponse, error)
}

// ClinicReportingService implements the analytics collection functionality.
type ClinicReportingService struct {
	storeDB db.Store
}

func (s *ClinicReportingService) GetClinicAnalytics(ctx *gin.Context) (*ClinicAnalyticsResponse, error) {
	// Retrieve all appointments.
	appointments, err := s.storeDB.ListAllAppointments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve appointments: %w", err)
	}

	// Retrieve all pets.
	pets, err := s.storeDB.ListPets(ctx, db.ListPetsParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pets: %w", err)
	}

	// Retrieve all doctors.
	doctors, err := s.storeDB.GetAllDoctors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve doctors: %w", err)
	}

	totalAppointments := len(appointments)
	upcomingAppointments := 0
	pastAppointments := 0
	now := time.Now()

	// Assume each appointment has an AppointmentTime field of type time.Time.
	for _, app := range appointments {
		if app.Date.Time.After(now) {
			upcomingAppointments++
		} else {
			pastAppointments++
		}
	}

	return &ClinicAnalyticsResponse{
		TotalAppointments:    totalAppointments,
		UpcomingAppointments: upcomingAppointments,
		PastAppointments:     pastAppointments,
		TotalPets:            len(pets),
		TotalDoctors:         len(doctors),
	}, nil
}

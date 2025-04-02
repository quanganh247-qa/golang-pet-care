package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

// Routes sets up all appointment related routes
// @Summary Appointment routes setup
// @Description Initializes all appointment related API endpoints
// @Tags appointments
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
	appointment := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(appointment)

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
			},
		},
	}

	{
		// @Summary Create a new appointment
		// @Description Create a new appointment for a pet
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param appointment body CreateAppointmentRequest true "Appointment details"
		// @Success 201 {object} AppointmentResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment [post]
		// @Security BearerAuth
		authRoute.POST("appointment/", appointmentApi.controller.createAppointment)

		// @Summary Confirm an appointment
		// @Description Confirm an existing appointment by ID
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Success 200 {object} SuccessResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 404 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/confirm/{id} [post]
		// @Security BearerAuth
		authRoute.POST("appointment/confirm/:id", appointmentApi.controller.confirmAppointment)

		// @Summary Check-in an appointment
		// @Description Check-in a patient for an existing appointment
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Success 200 {object} SuccessResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 404 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/check-in/{id} [post]
		// @Security BearerAuth
		authRoute.POST("appointment/check-in/:id", appointmentApi.controller.checkinAppointment)

		// @Summary Get user appointments
		// @Description Get all appointments for the authenticated user
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Success 200 {array} AppointmentResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/user [get]
		// @Security BearerAuth
		authRoute.GET("appointment/user", appointmentApi.controller.getAppointmentsByUser)

		// @Summary Get appointment by ID
		// @Description Get appointment details by ID
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Success 200 {object} AppointmentResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 404 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/{id} [get]
		// @Security BearerAuth
		authRoute.GET("appointment/:id", appointmentApi.controller.getAppointmentByID)

		// @Summary Get appointments by doctor
		// @Description Get all appointments for a specific doctor
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param doctor_id path string true "Doctor ID"
		// @Success 200 {array} createAppointmentResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/doctor/{doctor_id} [get]
		// @Security BearerAuth
		authRoute.GET("appointment/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)

		// @Summary Get all appointments
		// @Description Get all appointments with pagination, filtering by date and option
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param date query string true "Date in format YYYY-MM-DD"
		// @Param option query string false "Filter option"
		// @Param page query int false "Page number"
		// @Param page_size query int false "Page size"
		// @Success 200 {object} util.PaginationResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointments [get]
		// @Security BearerAuth
		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)

		// @Summary Update an appointment
		// @Description Update an existing appointment by ID
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Param appointment body updateAppointmentRequest true "Updated appointment details"
		// @Success 200 {object} SuccessResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 404 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/{id} [put]
		// @Security BearerAuth
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)

		// @Summary Get available time slots
		// @Description Get available time slots for a specific doctor on a given date
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param doctor_id path string true "Doctor ID"
		// @Param date query string true "Date in format YYYY-MM-DD"
		// @Success 200 {array} timeSlotResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /doctor/{doctor_id}/time-slot [get]
		// @Security BearerAuth
		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)

		// @Summary Get history appointments by pet ID
		// @Description Get all history appointments for a specific pet
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param pet_id path string true "Pet ID"
		// @Success 200 {array} historyAppointmentResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointments/pet/{pet_id}/history [get]
		// @Security BearerAuth
		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)

		// @Summary Get queue
		// @Description Get all queue items for the authenticated user
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Success 200 {array} QueueItem
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointments/queue [get]
		// @Security BearerAuth
		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)

		// @Summary Update queue item status
		// @Description Update the status of a queue item
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Queue Item ID"
		// @Param status body struct{Status string} true "New status"
		// @Success 200 {object} gin.H
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointments/queue/{id}/status [put]
		// @Security BearerAuth
		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)

		// SOAP endpoints
		// @Summary Create SOAP note
		// @Description Create a SOAP note for an appointment
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Param soap body CreateSOAPRequest true "SOAP note details"
		// @Success 200 {object} SOAPResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/{id}/soap [post]
		// @Security BearerAuth
		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)

		// @Summary Update SOAP note
		// @Description Update a SOAP note for an appointment
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Param soap body UpdateSOAPRequest true "Updated SOAP note details"
		// @Success 200 {object} SOAPResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/{id}/soap [put]
		// @Security BearerAuth
		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)

		// @Summary Get SOAP note
		// @Description Get a SOAP note for an appointment
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param id path string true "Appointment ID"
		// @Success 200 {object} SOAPResponse
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointment/{id}/soap [get]
		// @Security BearerAuth
		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)

		// Statistics
		// @Summary Get appointment distribution
		// @Description Get appointment distribution statistics by service
		// @Tags appointments
		// @Accept json
		// @Produce json
		// @Param start_date query string true "Start date in format YYYY-MM-DD"
		// @Param end_date query string true "End date in format YYYY-MM-DD"
		// @Success 200 {object} gin.H{success=bool,data=[]AppointmentDistribution}
		// @Failure 400 {object} ErrorResponse
		// @Failure 401 {object} ErrorResponse
		// @Failure 500 {object} ErrorResponse
		// @Router /appointments/statistic [get]
		// @Security BearerAuth
		authRoute.GET("appointments/statistic", appointmentApi.controller.GetAppointmentDistribution)
	}
}

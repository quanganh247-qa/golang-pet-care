// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error)
	CreatePet(ctx context.Context, arg CreatePetParams) (Pet, error)
	CreateReminder(ctx context.Context, arg CreateReminderParams) (Reminder, error)
	CreateReminder(ctx context.Context, arg CreateReminderParams) (Reminder, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceType(ctx context.Context, arg CreateServiceTypeParams) (Servicetype, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVaccination(ctx context.Context, arg CreateVaccinationParams) (Vaccination, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	DeleteFeedingSchedule(ctx context.Context, feedingScheduleID int64) error
	DeletePet(ctx context.Context, petid int64) error
	DeleteReminder(ctx context.Context, reminderID int64) error
	DeleteService(ctx context.Context, serviceid int64) error
	DeleteServiceType(ctx context.Context, typeid int64) error
	DeleteVaccination(ctx context.Context, vaccinationid int64) error
	GetActiveDoctors(ctx context.Context, arg GetActiveDoctorsParams) ([]GetActiveDoctorsRow, error)
	GetAllServices(ctx context.Context, arg GetAllServicesParams) ([]Service, error)
<<<<<<< HEAD
	GetAllTimeSlots(ctx context.Context, arg GetAllTimeSlotsParams) ([]GetAllTimeSlotsRow, error)
=======
>>>>>>> 1c933f4 (update service api)
	GetAllTimeSlots(ctx context.Context, arg GetAllTimeSlotsParams) ([]GetAllTimeSlotsRow, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetAppointmentsOfDoctor(ctx context.Context, id int64) ([]GetAppointmentsOfDoctorRow, error)
	GetDoctor(ctx context.Context, id int64) (GetDoctorRow, error)
	GetFeedingScheduleByPetID(ctx context.Context, petid pgtype.Int8) ([]Feedingschedule, error)
	GetPetByID(ctx context.Context, petid int64) (Pet, error)
	GetServiceByID(ctx context.Context, serviceid int64) (Service, error)
	GetServiceType(ctx context.Context, typeid int64) (Servicetype, error)
	GetTimeSlotByID(ctx context.Context, id int64) (GetTimeSlotByIDRow, error)
	GetTimeslotsAvailable(ctx context.Context, arg GetTimeslotsAvailableParams) ([]GetTimeslotsAvailableRow, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetVaccinationByID(ctx context.Context, vaccinationid int64) (Vaccination, error)
	InsertDoctor(ctx context.Context, arg InsertDoctorParams) (Doctor, error)
	InsertDoctorSchedule(ctx context.Context, arg InsertDoctorScheduleParams) (Doctorschedule, error)
	InsertTimeslot(ctx context.Context, arg InsertTimeslotParams) (Timeslot, error)
	ListPets(ctx context.Context, arg ListPetsParams) ([]Pet, error)
	UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error
	// Replace $2 with the specific date (YYYY-MM-DD)
	UpdateDoctorAvailable(ctx context.Context, arg UpdateDoctorAvailableParams) error
	UpdateNotification(ctx context.Context, appointmentID int64) error
	UpdatePet(ctx context.Context, arg UpdatePetParams) error
	UpdateService(ctx context.Context, arg UpdateServiceParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVaccination(ctx context.Context, arg UpdateVaccinationParams) error
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)

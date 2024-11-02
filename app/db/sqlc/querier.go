// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error)
	CreatePet(ctx context.Context, arg CreatePetParams) (Pet, error)
	CreateReminder(ctx context.Context, arg CreateReminderParams) (Reminder, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceType(ctx context.Context, arg CreateServiceTypeParams) (Servicetype, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	DeletePet(ctx context.Context, petid int64) error
	DeleteService(ctx context.Context, serviceid int64) error
	DeleteServiceType(ctx context.Context, typeid int64) error
	GetActiveDoctors(ctx context.Context, arg GetActiveDoctorsParams) ([]GetActiveDoctorsRow, error)
	GetAllServices(ctx context.Context, arg GetAllServicesParams) ([]Service, error)
	GetAllTimeSlots(ctx context.Context, arg GetAllTimeSlotsParams) ([]GetAllTimeSlotsRow, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetAppointmentsOfDoctorWithDetails(ctx context.Context, id int64) ([]GetAppointmentsOfDoctorWithDetailsRow, error)
	GetDoctor(ctx context.Context, id int64) (GetDoctorRow, error)
	GetPetByID(ctx context.Context, petid int64) (Pet, error)
	GetServiceByID(ctx context.Context, serviceid int64) (Service, error)
	GetServiceType(ctx context.Context, typeid int64) (Servicetype, error)
	GetTimeSlotByID(ctx context.Context, id int64) (GetTimeSlotByIDRow, error)
	GetTimeslotsAvailable(ctx context.Context, arg GetTimeslotsAvailableParams) ([]GetTimeslotsAvailableRow, error)
	GetUser(ctx context.Context, username string) (User, error)
	InsertDoctor(ctx context.Context, arg InsertDoctorParams) (Doctor, error)
	InsertDoctorSchedule(ctx context.Context, arg InsertDoctorScheduleParams) (Doctorschedule, error)
	InsertTimeslot(ctx context.Context, arg InsertTimeslotParams) (Timeslot, error)
	InsertTokenInfo(ctx context.Context, arg InsertTokenInfoParams) (TokenInfo, error)
	ListPets(ctx context.Context, arg ListPetsParams) ([]Pet, error)
	UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error
	// Replace $2 with the specific date (YYYY-MM-DD)
	UpdateDoctorAvailable(ctx context.Context, arg UpdateDoctorAvailableParams) error
	UpdateNotification(ctx context.Context, appointmentID int64) error
	UpdatePet(ctx context.Context, arg UpdatePetParams) error
	UpdateService(ctx context.Context, arg UpdateServiceParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)

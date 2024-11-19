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
	CreatePetSchedule(ctx context.Context, arg CreatePetScheduleParams) error
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateServiceType(ctx context.Context, arg CreateServiceTypeParams) (Servicetype, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (int64, error)
	CreateVaccination(ctx context.Context, arg CreateVaccinationParams) (Vaccination, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	DeleteDeviceToken(ctx context.Context, arg DeleteDeviceTokenParams) error
	DeletePet(ctx context.Context, petid int64) error
	DeletePetLog(ctx context.Context, logID int64) error
	DeleteService(ctx context.Context, serviceid int64) error
	DeleteServiceType(ctx context.Context, typeid int64) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteVaccination(ctx context.Context, vaccinationid int64) error
	GetActiveDoctors(ctx context.Context, arg GetActiveDoctorsParams) ([]GetActiveDoctorsRow, error)
	GetAllSchedulesByPet(ctx context.Context, arg GetAllSchedulesByPetParams) ([]PetSchedule, error)
	GetAllServices(ctx context.Context, arg GetAllServicesParams) ([]Service, error)
	GetAllTimeSlots(ctx context.Context, arg GetAllTimeSlotsParams) ([]GetAllTimeSlotsRow, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetAppointmentsOfDoctorWithDetails(ctx context.Context, id int64) ([]GetAppointmentsOfDoctorWithDetailsRow, error)
	GetDeviceTokenByUsername(ctx context.Context, username string) ([]Devicetoken, error)
	// 1. Query cơ bản để lấy thông tin bệnh và thuốc điều trị
	GetDiceaseAndMedicinesInfo(ctx context.Context, lower string) ([]GetDiceaseAndMedicinesInfoRow, error)
	GetDiseaseTreatmentPlanWithPhases(ctx context.Context, lower string) ([]GetDiseaseTreatmentPlanWithPhasesRow, error)
	GetDoctor(ctx context.Context, id int64) (GetDoctorRow, error)
	GetPetByID(ctx context.Context, petid int64) (Pet, error)
	GetPetLogByID(ctx context.Context, arg GetPetLogByIDParams) (GetPetLogByIDRow, error)
	GetPetLogsByPetID(ctx context.Context, arg GetPetLogsByPetIDParams) ([]GetPetLogsByPetIDRow, error)
	GetServiceByID(ctx context.Context, serviceid int64) (Service, error)
	GetServiceType(ctx context.Context, typeid int64) (Servicetype, error)
	GetTimeSlotByID(ctx context.Context, id int64) (GetTimeSlotByIDRow, error)
	GetTimeslotsAvailable(ctx context.Context, arg GetTimeslotsAvailableParams) ([]GetTimeslotsAvailableRow, error)
	GetUser(ctx context.Context, username string) (GetUserRow, error)
	GetVaccinationByID(ctx context.Context, vaccinationid int64) (Vaccination, error)
	InsertDeviceToken(ctx context.Context, arg InsertDeviceTokenParams) (Devicetoken, error)
	InsertDoctor(ctx context.Context, arg InsertDoctorParams) (Doctor, error)
	InsertDoctorSchedule(ctx context.Context, arg InsertDoctorScheduleParams) (Doctorschedule, error)
	InsertPetLog(ctx context.Context, arg InsertPetLogParams) (PetLog, error)
	InsertTimeslot(ctx context.Context, arg InsertTimeslotParams) (Timeslot, error)
	ListPetSchedulesByUsername(ctx context.Context, username string) ([]ListPetSchedulesByUsernameRow, error)
	ListPets(ctx context.Context, arg ListPetsParams) ([]Pet, error)
	ListPetsByUsername(ctx context.Context, arg ListPetsByUsernameParams) ([]Pet, error)
	ListVaccinationsByPetID(ctx context.Context, petid pgtype.Int8) ([]Vaccination, error)
	SetPetInactive(ctx context.Context, arg SetPetInactiveParams) error
	UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error
	// Replace $2 with the specific date (YYYY-MM-DD)
	UpdateDoctorAvailable(ctx context.Context, arg UpdateDoctorAvailableParams) error
	UpdateNotification(ctx context.Context, appointmentID int64) error
	UpdatePet(ctx context.Context, arg UpdatePetParams) error
	UpdatePetLog(ctx context.Context, arg UpdatePetLogParams) error
	UpdateService(ctx context.Context, arg UpdateServiceParams) error
	UpdateVaccination(ctx context.Context, arg UpdateVaccinationParams) error
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
	VerifiedUser(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)

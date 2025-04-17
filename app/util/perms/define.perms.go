package perms

type (
	Permission string
)

// Define roles
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleDoctor = "doctor"
)

// Define permissions
const (
	// Doctor
	ManageDoctor Permission = "MANAGE_DOCTOR"

	// Appointment
	CreateAppointment Permission = "CREATE_APPOINTMENT"
	ReadAppointment   Permission = "READ_APPOINTMENT"
	UpdateAppointment Permission = "UPDATE_APPOINTMENT"
	DeleteAppointment Permission = "DELETE_APPOINTMENT"
	ManageUsers       Permission = "MANAGE_USERS"
	ManageDoctors     Permission = "MANAGE_DOCTORS"

	// Service
	ManageServices Permission = "MANAGE_SERVICES"
	ReadService    Permission = "READ_SERVICE"
	CreateService  Permission = "CREATE_SERVICE"
	UpdateService  Permission = "UPDATE_SERVICE"
	DeleteService  Permission = "DELETE_SERVICE"

	// Medical Records
	ManageMedicalRecords Permission = "MANAGE_MEDICAL_RECORDS"
	ReadMedicalRecord    Permission = "READ_MEDICAL_RECORD"

	// Disease and Treatment
	ManageDisease   Permission = "MANAGE_DISEASE"
	ManageTreatment Permission = "MANAGE_TREATMENT"

	// Pet
	ManagePets Permission = "MANAGE_PETS"
	ReadPet    Permission = "READ_PET"
	CreatePet  Permission = "CREATE_PET"
	UpdatePet  Permission = "UPDATE_PET"
	DeletePet  Permission = "DELETE_PET"

	ViewReports Permission = "VIEW_REPORTS"

	// Test
	ManageTest Permission = "MANAGE_TEST"
	ReadTest   Permission = "READ_TEST"
	CreateTest Permission = "CREATE_TEST"

	// System
	ManageSystem Permission = "MANAGE_SYSTEM"
)

// RolePermissions maps roles to their respective permissions
var RolePermissions = map[string][]Permission{
	RoleAdmin: {CreateAppointment, ReadAppointment, UpdateAppointment,
		DeleteAppointment, ManageUsers, ManageDoctors,
		ViewReports, ManageServices, ManageMedicalRecords,
		ManageDoctor, ManageSystem},
	RoleDoctor: {CreateAppointment, ReadAppointment, UpdateAppointment,
		ViewReports, ManageServices, ManageMedicalRecords,
		ManageTreatment, ManageDisease, ManageDoctor,
		ManageTest,
	},
	RoleUser: {CreateAppointment, ReadAppointment, ReadMedicalRecord,
		ReadService, ReadPet, ReadTest},
}

// CheckPermission checks if the user's role has the required permissions
func CheckPermission(requiredPermissions []Permission, userRole string) bool {
	// Get permissions for the user's role
	userPermissions, exists := RolePermissions[userRole]
	if !exists {
		return false // Role not found
	}

	// Check if the user's permissions include all required permissions
	for _, requiredPerm := range requiredPermissions {
		hasPermission := false
		for _, userPerm := range userPermissions {
			if userPerm == requiredPerm {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return false
		}
	}

	return true
}

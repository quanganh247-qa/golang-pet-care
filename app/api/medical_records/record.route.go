package medical_records

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	MedicalRecord := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(MedicalRecord)
	// perRoute := routerGroup.RouterPermission(MedicalRecord)

	// Apply cache middleware to GET endpoints with appropriate time-to-live
	MedicalRecord.Use(middleware.CacheMiddleware(time.Minute*5, "medical_records", []string{"GET"}))

	// Khoi tao api
	MedicalRecordApi := &MedicalRecordApi{
		&MedicalRecordController{
			service: &MedicalRecordService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		// Existing routes for basic medical records
		authRoute.GET("/pets/:pet_id/medical-records", MedicalRecordApi.controller.GetMedicalRecord)
		authRoute.POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
		authRoute.GET("/pets/:pet_id/medical-records/:medical_record_id/medical-histories", MedicalRecordApi.controller.ListMedicalHistory)
		authRoute.POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
		authRoute.GET("/pets/:pet_id/medical-histories/:id", MedicalRecordApi.controller.GetMedicalHistoryByID)
		authRoute.GET("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.GetMedicalHistoryByPetID)

		// New routes for comprehensive medical history

		// Complete medical history endpoint
		authRoute.GET("/pets/:pet_id/complete-medical-history", MedicalRecordApi.controller.GetCompleteMedicalHistory)

		// Examination routes
		authRoute.POST("/medical-records/examinations", MedicalRecordApi.controller.CreateExamination)
		authRoute.GET("/medical-records/examinations/:examination_id", MedicalRecordApi.controller.GetExamination)
		authRoute.GET("/medical-records/examinations", MedicalRecordApi.controller.ListExaminations)

		// Prescription routes
		authRoute.POST("/medical-records/prescriptions", MedicalRecordApi.controller.CreatePrescription)
		authRoute.GET("/medical-records/prescriptions/:prescription_id", MedicalRecordApi.controller.GetPrescription)
		authRoute.GET("/medical-records/prescriptions", MedicalRecordApi.controller.ListPrescriptions)

		// Test result routes
		authRoute.POST("/medical-records/test-results", MedicalRecordApi.controller.CreateTestResult)
		authRoute.GET("/medical-records/test-results/:test_result_id", MedicalRecordApi.controller.GetTestResult)
		authRoute.GET("/medical-records/test-results", MedicalRecordApi.controller.ListTestResults)
	}
}

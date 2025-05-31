package reports

import (
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type Handler struct {
	store db.Store
}

func NewHandler(store db.Store) *Handler {
	return &Handler{store: store}
}

type FinancialReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type FinancialReportResponse struct {
	TotalRevenue     float64            `json:"total_revenue"`
	RevenueByService map[string]float64 `json:"revenue_by_service"`
	RevenueByProduct map[string]float64 `json:"revenue_by_product"`
	TotalExpenses    float64            `json:"total_expenses"`
	Profit           float64            `json:"profit"`
	MonthlyTrends    []MonthlyData      `json:"monthly_trends"`
}

type MonthlyData struct {
	Month    string  `json:"month"`
	Revenue  float64 `json:"revenue"`
	Expenses float64 `json:"expenses"`
	Profit   float64 `json:"profit"`
}

// Doctor statistics request and response structures
type DoctorStatsRequest struct {
	DoctorID  int64  `form:"doctor_id"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type DoctorStatsResponse struct {
	DoctorID              int64                    `json:"doctor_id"`
	DoctorName            string                   `json:"doctor_name"`
	TotalAppointments     int64                    `json:"total_appointments"`
	CompletedAppointments int64                    `json:"completed_appointments"`
	ScheduledAppointments int64                    `json:"scheduled_appointments"`
	CancelledAppointments int64                    `json:"cancelled_appointments"`
	UniquePatientsServed  int64                    `json:"unique_patients_served"`
	AppointmentsByMonth   []MonthlyAppointmentData `json:"appointments_by_month"`
	PatientsByMonth       []MonthlyPatientData     `json:"patients_by_month"`
	CompletionRate        float64                  `json:"completion_rate"`
	WorkingDays           int64                    `json:"working_days"`
	AvgAppointmentsPerDay float64                  `json:"avg_appointments_per_day"`
}

type MonthlyAppointmentData struct {
	Month     string `json:"month"`
	Year      int    `json:"year"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Scheduled int64  `json:"scheduled"`
	Cancelled int64  `json:"cancelled"`
}

type MonthlyPatientData struct {
	Month          string `json:"month"`
	Year           int    `json:"year"`
	UniquePatients int64  `json:"unique_patients"`
	NewPatients    int64  `json:"new_patients"`
	ReturnPatients int64  `json:"return_patients"`
}

// All doctors statistics response
type AllDoctorsStatsResponse struct {
	TotalDoctors  int64                 `json:"total_doctors"`
	PeriodStats   *DoctorStatsResponse  `json:"period_stats"`
	DoctorsList   []DoctorStatsResponse `json:"doctors_list"`
	TopPerformers []DoctorStatsResponse `json:"top_performers"`
}

func (h *Handler) GetFinancialReport(ctx *gin.Context) {
	var req FinancialReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Get appointments revenue
	appointmentsRevenue, err := h.getAppointmentsRevenue(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get appointments revenue"})
		return
	}

	// Get product sales revenue
	productRevenue, productRevenueByCategory, err := h.getProductRevenue(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product revenue"})
		return
	}

	// Get service revenue by category
	serviceRevenueByCategory, err := h.getServiceRevenueByCategory(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get service revenue by category"})
		return
	}

	// Get expenses (medicine purchases, etc.)
	expenses, err := h.getTotalExpenses(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get expenses"})
		return
	}

	// Calculate monthly trends
	monthlyTrends, err := h.getMonthlyFinancialTrends(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get monthly trends"})
		return
	}

	totalRevenue := appointmentsRevenue + productRevenue
	profit := totalRevenue - expenses

	response := FinancialReportResponse{
		TotalRevenue:     totalRevenue,
		RevenueByService: serviceRevenueByCategory,
		RevenueByProduct: productRevenueByCategory,
		TotalExpenses:    expenses,
		Profit:           profit,
		MonthlyTrends:    monthlyTrends,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) getAppointmentsRevenue(ctx *gin.Context, startDate, endDate time.Time) (float64, error) {
	// Query to get revenue from appointments
	var revenue float64

	statement := `
		SELECT COALESCE(SUM(p.amount), 0) as revenue
		FROM payments p
		JOIN appointments a ON p.appointment_id = a.appointment_id
		WHERE p.payment_status = 'completed' 
		AND p.created_at BETWEEN $1 AND $2
	`

	params := []interface{}{startDate, endDate}
	outputs := []string{"revenue"}

	result, err := h.store.ExecStatementOne(ctx, statement, params, outputs)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if value, ok := result["revenue"]; ok {
		revenue, _ = value.(float64)
	}

	return revenue, nil
}

func (h *Handler) getProductRevenue(ctx *gin.Context, startDate, endDate time.Time) (float64, map[string]float64, error) {
	// Query to get revenue from product sales
	var totalRevenue float64
	revenueByCategory := make(map[string]float64)

	statement := `
		SELECT 
			COALESCE(SUM(o.total_amount), 0) as total
		FROM 
			orders o
		WHERE 
			o.payment_status = 'completed' 
			AND o.order_date BETWEEN $1 AND $2

	`

	params := []interface{}{startDate, endDate}

	results, err := h.store.ExecStatementMany(ctx, statement, params)
	if err != nil {
		return 0, nil, err
	}

	for _, result := range results {
		if total, ok := result["total"].(float64); ok {
			totalRevenue = total
		}

		// category, categoryOk := result["category"].(string)
		// categoryRevenue, revenueOk := result["category_revenue"].(float64)

		// if categoryOk && revenueOk {
		// 	revenueByCategory[category] = categoryRevenue
		// }
	}

	return totalRevenue, revenueByCategory, nil
}

func (h *Handler) getServiceRevenueByCategory(ctx *gin.Context, startDate, endDate time.Time) (map[string]float64, error) {
	revenueByCategory := make(map[string]float64)

	results, err := h.store.ExecStatementMany(ctx, `
		SELECT 
			s.category,
			COALESCE(SUM(p.amount), 0) as category_revenue
		FROM 
			payments p
		JOIN 
			appointments a ON p.appointment_id = a.appointment_id 
		JOIN 
			services s ON a.service_id = s.id
		WHERE 
			p.payment_status = 'completed' 
			AND p.created_at BETWEEN $1 AND $2
		GROUP BY 
			s.category
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if category, ok := result["category"].(string); ok {
			categoryRevenue, revenueOk := result["category_revenue"].(float64)

			if revenueOk {
				revenueByCategory[category] = categoryRevenue
			}
		}

	}

	return revenueByCategory, nil
}

func (h *Handler) getTotalExpenses(ctx *gin.Context, startDate, endDate time.Time) (float64, error) {
	// Query to get expenses (primarily from medicine transactions)
	var expenses float64
	result, err := h.store.ExecStatementOne(ctx, `
		SELECT COALESCE(SUM(mt.quantity * mt.unit_price), 0)
		FROM medicine_transactions mt
		WHERE mt.transaction_type = 'import'
		AND mt.transaction_date BETWEEN $1 AND $2
	`, []interface{}{startDate, endDate}, []string{"expenses"})

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if value, ok := result["expenses"]; ok {
		expenses, _ = value.(float64)
	}

	return expenses, nil
}

func (h *Handler) getMonthlyFinancialTrends(ctx *gin.Context, startDate, endDate time.Time) ([]MonthlyData, error) {
	var trends []MonthlyData

	results, err := h.store.ExecStatementMany(ctx, `
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', $1::timestamp)::date,
				date_trunc('month', $2::timestamp)::date,
				'1 month'::interval
			) AS month_start
		),
		revenue_data AS (
			SELECT 
				date_trunc('month', p.created_at) as month,
				COALESCE(SUM(p.amount), 0) as revenue
			FROM payments p
			WHERE p.payment_status = 'completed'
			AND p.created_at BETWEEN $1 AND $2
			GROUP BY date_trunc('month', p.created_at)
		),
		expense_data AS (
			SELECT 
				date_trunc('month', mt.transaction_date) as month,
				COALESCE(SUM(mt.quantity * mt.unit_price), 0) as expenses
			FROM medicine_transactions mt
			WHERE mt.transaction_type = 'import'
			AND mt.transaction_date BETWEEN $1 AND $2
			GROUP BY date_trunc('month', mt.transaction_date)
		)
		SELECT 
			to_char(m.month_start, 'YYYY-MM') as month,
			COALESCE(r.revenue, 0) as revenue,
			COALESCE(e.expenses, 0) as expenses,
			COALESCE(r.revenue, 0) - COALESCE(e.expenses, 0) as profit
		FROM months m
		LEFT JOIN revenue_data r ON m.month_start = r.month
		LEFT JOIN expense_data e ON m.month_start = e.month
		ORDER BY m.month_start
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var data MonthlyData

		if month, ok := result["month"].(string); ok {
			data.Month = month
		}

		if revenue, ok := result["revenue"].(float64); ok {
			data.Revenue = revenue
		}
		if expenses, ok := result["expenses"].(float64); ok {
			data.Expenses = expenses
		}
		if profit, ok := result["profit"].(float64); ok {
			data.Profit = profit
		}

		trends = append(trends, data)
	}

	return trends, nil
}

// Medical records reports
type MedicalRecordsReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type MedicalRecordsReportResponse struct {
	TotalExaminations  int64              `json:"total_examinations"`
	CommonDiseases     []DiseaseStats     `json:"common_diseases"`
	ExaminationsByType map[string]int     `json:"examinations_by_type"`
	MonthlyTrends      []MonthlyExamStats `json:"monthly_trends"`
}

type DiseaseStats struct {
	Disease    string  `json:"disease"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type MonthlyExamStats struct {
	Month        string `json:"month"`
	Examinations int64  `json:"examinations"`
}

func (h *Handler) GetMedicalRecordsReport(ctx *gin.Context) {
	var req MedicalRecordsReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Get total number of examinations
	totalExams, err := h.getTotalExaminations(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total examinations"})
		return
	}

	// Get common diseases
	commonDiseases, err := h.getCommonDiseases(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get common diseases"})
		return
	}

	// Get examinations by type
	examsByType, err := h.getExaminationsByType(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get examinations by type"})
		return
	}

	// Get monthly trends
	monthlyTrends, err := h.getMonthlyExaminationTrends(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get monthly trends"})
		return
	}

	response := MedicalRecordsReportResponse{
		TotalExaminations:  totalExams,
		CommonDiseases:     commonDiseases,
		ExaminationsByType: examsByType,
		MonthlyTrends:      monthlyTrends,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) getTotalExaminations(ctx *gin.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	result, err := h.store.ExecStatementOne(ctx, `
		SELECT COUNT(*) 
		FROM consultations c
		JOIN appointments a ON c.appointment_id = a.appointment_id
		WHERE c.created_at BETWEEN $1 AND $2
		AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
	`, []interface{}{startDate, endDate}, []string{"count"})

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if value, ok := result["count"]; ok {
		count, _ = value.(int64)
	}

	return count, nil
}

func (h *Handler) getCommonDiseases(ctx *gin.Context, startDate, endDate time.Time) ([]DiseaseStats, error) {
	var diseases []DiseaseStats

	results, err := h.store.ExecStatementMany(ctx, `
		WITH disease_counts AS (
			SELECT 
				COALESCE(c.assessment->>'diagnosis', 'Unknown') as disease,
				COUNT(*) as count
			FROM 
				consultations c
			JOIN 
				appointments a ON c.appointment_id = a.appointment_id
			WHERE 
				c.created_at BETWEEN $1 AND $2
				AND c.assessment IS NOT NULL
				AND c.assessment->>'diagnosis' IS NOT NULL
				AND c.assessment->>'diagnosis' != ''
			GROUP BY 
				c.assessment->>'diagnosis'
			ORDER BY 
				count DESC
			LIMIT 10
		),
		total AS (
			SELECT COUNT(*) as total 
			FROM consultations c
			JOIN appointments a ON c.appointment_id = a.appointment_id
			WHERE c.created_at BETWEEN $1 AND $2
			AND c.assessment IS NOT NULL
			AND c.assessment->>'diagnosis' IS NOT NULL
			AND c.assessment->>'diagnosis' != ''
		)
		SELECT 
			dc.disease,
			dc.count,
			(dc.count * 100.0 / NULLIF(t.total, 0)) as percentage
		FROM 
			disease_counts dc
		CROSS JOIN 
			total t
		ORDER BY 
			dc.count DESC
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var ds DiseaseStats

		if disease, ok := result["disease"].(string); ok {
			ds.Disease = disease
		}
		if count, ok := result["count"].(int64); ok {
			ds.Count = count
		}
		if percentage, ok := result["percentage"].(pgtype.Numeric); ok {
			// Convert numeric value to float64, accounting for the exponent
			// Get the big.Int as a float64
			floatValue, _ := new(big.Float).SetInt(percentage.Int).Float64()

			// Apply the exponent by multiplying by 10^Exp
			ds.Percentage = floatValue * math.Pow10(int(percentage.Exp))
		}

		diseases = append(diseases, ds)
	}

	return diseases, nil
}

func (h *Handler) getExaminationsByType(ctx *gin.Context, startDate, endDate time.Time) (map[string]int, error) {
	examsByType := make(map[string]int)

	results, err := h.store.ExecStatementMany(ctx, `
		SELECT 
			s.category as service_type,
			COUNT(*) as count
		FROM 
			appointments a
		JOIN 
			services s ON a.service_id = s.id
		WHERE 
			a.date BETWEEN $1 AND $2
			AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
		GROUP BY 
			s.category
		ORDER BY 
			count DESC
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var serviceType string
		var count int

		if serviceType, ok := result["service_type"].(string); ok {
			examsByType[serviceType] = count
		}
		if count, ok := result["count"].(int); ok {
			examsByType[serviceType] = count
		}
	}

	return examsByType, nil
}

func (h *Handler) getMonthlyExaminationTrends(ctx *gin.Context, startDate, endDate time.Time) ([]MonthlyExamStats, error) {
	var trends []MonthlyExamStats

	results, err := h.store.ExecStatementMany(ctx, `
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', $1::timestamp)::date,
				date_trunc('month', $2::timestamp)::date,
				'1 month'::interval
			) AS month_start
		),
		exams_data AS (
			SELECT 
				date_trunc('month', c.created_at) as month,
				COUNT(*) as count
			FROM 
				consultations c
			JOIN 
				appointments a ON c.appointment_id = a.appointment_id
			WHERE 
				a.created_at BETWEEN $1 AND $2
				AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
			GROUP BY 
				date_trunc('month', c.created_at)
		)
		SELECT 
			to_char(m.month_start, 'YYYY-MM') as month,
			COALESCE(e.count, 0) as examinations
		FROM 
			months m
		LEFT JOIN 
			exams_data e ON m.month_start = e.month
		ORDER BY 
			m.month_start
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var data MonthlyExamStats

		if month, ok := result["month"].(string); ok {
			data.Month = month
		}
		if count, ok := result["examinations"].(int64); ok {
			data.Examinations = count
		}

		trends = append(trends, data)
	}

	return trends, nil
}

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type UpdateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status string `json:"status" binding:"required,oneof=available occupied maintenance"`
}

type RoomResponse struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Status               string     `json:"status"`
	Type                 string     `json:"type"`
	CurrentAppointmentID *int64     `json:"current_appointment_id,omitempty"`
	AvailableAt          *time.Time `json:"available_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type ListRoomsResponse struct {
	Rooms      []RoomResponse `json:"rooms"`
	TotalCount int32          `json:"total_count"`
}

func (h *Handler) ListRooms(ctx *gin.Context) ([]RoomResponse, error) {

	rooms, err := h.store.GetAvailableRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %v", err)
	}

	response := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		response[i] = RoomResponse{
			ID:   room.ID,
			Name: room.Name,
			Type: room.Type,
		}
	}

	return response, nil
}

func (h *Handler) ListRoomHandler(ctx *gin.Context) {
	rooms, err := h.ListRooms(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}

// GetDoctorStats lấy thống kê appointment và bệnh nhân của một bác sĩ cụ thể
func (h *Handler) GetDoctorStats(ctx *gin.Context) {
	var req DoctorStatsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Get doctor info
	doctor, err := h.store.GetDoctor(ctx, req.DoctorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Get doctor statistics
	stats, err := h.getDoctorStatistics(ctx, req.DoctorID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get doctor statistics"})
		return
	}

	stats.DoctorName = doctor.Name

	ctx.JSON(http.StatusOK, stats)
}

// GetAllDoctorsStats lấy thống kê của tất cả bác sĩ
func (h *Handler) GetAllDoctorsStats(ctx *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"required"`
		EndDate   string `form:"end_date" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Get all doctors
	doctors, err := h.store.GetDoctors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get doctors list"})
		return
	}

	var doctorsList []DoctorStatsResponse
	var totalAppointments, totalCompleted, totalScheduled, totalCancelled int64
	var totalUniquePatients int64

	for _, doctor := range doctors {
		stats, err := h.getDoctorStatistics(ctx, doctor.ID, startDate, endDate)
		if err != nil {
			continue // Skip this doctor if error occurs
		}

		// Get doctor name
		doctorInfo, err := h.store.GetDoctor(ctx, doctor.ID)
		if err != nil {
			continue
		}
		if conditions := doctorInfo.Role.String; conditions != "doctor" {
			continue // Skip if not a doctor
		}
		stats.DoctorName = doctorInfo.Name

		doctorsList = append(doctorsList, *stats)

		// Accumulate totals
		totalAppointments += stats.TotalAppointments
		totalCompleted += stats.CompletedAppointments
		totalScheduled += stats.ScheduledAppointments
		totalCancelled += stats.CancelledAppointments
		totalUniquePatients += stats.UniquePatientsServed
	}

	// Create overall stats
	overallStats := &DoctorStatsResponse{
		DoctorID:              0,
		DoctorName:            "All Doctors",
		TotalAppointments:     totalAppointments,
		CompletedAppointments: totalCompleted,
		ScheduledAppointments: totalScheduled,
		CancelledAppointments: totalCancelled,
		UniquePatientsServed:  totalUniquePatients,
		CompletionRate:        float64(totalCompleted) / float64(totalAppointments) * 100,
	}

	// Sort doctors by total appointments (top performers)
	topPerformers := make([]DoctorStatsResponse, len(doctorsList))
	copy(topPerformers, doctorsList)

	// Simple sort by total appointments (descending)
	for i := 0; i < len(topPerformers)-1; i++ {
		for j := i + 1; j < len(topPerformers); j++ {
			if topPerformers[i].TotalAppointments < topPerformers[j].TotalAppointments {
				topPerformers[i], topPerformers[j] = topPerformers[j], topPerformers[i]
			}
		}
	}

	// Take top 5
	if len(topPerformers) > 5 {
		topPerformers = topPerformers[:5]
	}

	response := AllDoctorsStatsResponse{
		TotalDoctors:  int64(len(doctors)),
		PeriodStats:   overallStats,
		DoctorsList:   doctorsList,
		TopPerformers: topPerformers,
	}

	ctx.JSON(http.StatusOK, response)
}

// getDoctorStatistics lấy thống kê chi tiết của một bác sĩ
func (h *Handler) getDoctorStatistics(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) (*DoctorStatsResponse, error) {

	var err error
	var totalAppointments, completedAppointments, scheduledAppointments, cancelledAppointments, uniquePatients, workingDays int64
	var monthlyAppointments []MonthlyAppointmentData
	var monthlyPatients []MonthlyPatientData
	var completionRate, avgAppointmentsPerDay float64

	// Get total appointments
	totalAppointments, err = h.getDoctorAppointmentCount(ctx, doctorID, startDate, endDate, "")
	if err != nil {
		return nil, err
	}

	// Get completed appointments
	completedAppointments, err = h.getDoctorAppointmentCount(ctx, doctorID, startDate, endDate, "Completed")
	if err != nil {
		return nil, err
	}

	// Get scheduled appointments
	scheduledAppointments, err = h.getDoctorAppointmentCount(ctx, doctorID, startDate, endDate, "Scheduled")
	if err != nil {
		return nil, err
	}

	// Get cancelled appointments
	cancelledAppointments, err = h.getDoctorAppointmentCount(ctx, doctorID, startDate, endDate, "Cancelled")
	if err != nil {
		return nil, err
	}

	// Get unique patients served
	uniquePatients, err = h.getDoctorUniquePatients(ctx, doctorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get monthly appointment data
	monthlyAppointments, err = h.getDoctorMonthlyAppointments(ctx, doctorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get monthly patient data
	monthlyPatients, err = h.getDoctorMonthlyPatients(ctx, doctorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get working days
	workingDays, err = h.getDoctorWorkingDays(ctx, doctorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate completion rate
	completionRate = float64(0)
	if totalAppointments > 0 {
		completionRate = float64(completedAppointments) / float64(totalAppointments) * 100
	}

	// Calculate average appointments per day
	avgAppointmentsPerDay = float64(0)
	if workingDays > 0 {
		avgAppointmentsPerDay = float64(totalAppointments) / float64(workingDays)
	}

	return &DoctorStatsResponse{
		DoctorID:              doctorID,
		TotalAppointments:     totalAppointments,
		CompletedAppointments: completedAppointments,
		ScheduledAppointments: scheduledAppointments,
		CancelledAppointments: cancelledAppointments,
		UniquePatientsServed:  uniquePatients,
		AppointmentsByMonth:   monthlyAppointments,
		PatientsByMonth:       monthlyPatients,
		CompletionRate:        completionRate,
		WorkingDays:           workingDays,
		AvgAppointmentsPerDay: avgAppointmentsPerDay,
	}, nil
}

// getDoctorAppointmentCount đếm số lượng appointment của bác sĩ theo trạng thái
func (h *Handler) getDoctorAppointmentCount(ctx *gin.Context, doctorID int64, startDate, endDate time.Time, status string) (int64, error) {
	var statement string
	var params []interface{}

	if status == "" {
		statement = `
			SELECT COUNT(*) as count
			FROM appointments a
			WHERE a.doctor_id = $1 
			AND a.date BETWEEN $2 AND $3
		`
		params = []interface{}{doctorID, startDate, endDate}
	} else {
		statement = `
			SELECT COUNT(*) as count
			FROM appointments a
			JOIN states s ON a.state_id = s.id
			WHERE a.doctor_id = $1 
			AND a.date BETWEEN $2 AND $3
			AND s.state = $4
		`
		params = []interface{}{doctorID, startDate, endDate, status}
	}

	result, err := h.store.ExecStatementOne(ctx, statement, params, []string{"count"})
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var count int64
	if value, ok := result["count"]; ok {
		count, _ = value.(int64)
	}

	return count, nil
}

// getDoctorUniquePatients đếm số lượng bệnh nhân duy nhất mà bác sĩ đã phục vụ
func (h *Handler) getDoctorUniquePatients(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) (int64, error) {
	statement := `
		SELECT COUNT(DISTINCT a.petid) as count
		FROM appointments a
		WHERE a.doctor_id = $1 
		AND a.date BETWEEN $2 AND $3
	`

	result, err := h.store.ExecStatementOne(ctx, statement, []interface{}{doctorID, startDate, endDate}, []string{"count"})
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var count int64
	if value, ok := result["count"]; ok {
		count, _ = value.(int64)
	}

	return count, nil
}

// getDoctorMonthlyAppointments lấy thống kê appointment theo tháng
func (h *Handler) getDoctorMonthlyAppointments(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) ([]MonthlyAppointmentData, error) {
	statement := `
        SELECT 
            EXTRACT(YEAR FROM a.date) as year,
            EXTRACT(MONTH FROM a.date) as month,
            COUNT(*) as total,
            COUNT(CASE WHEN s.state = 'Completed' THEN 1 END) as completed,
            COUNT(CASE WHEN s.state = 'Scheduled' THEN 1 END) as scheduled,
            COUNT(CASE WHEN s.state = 'Cancelled' THEN 1 END) as cancelled
        FROM appointments a
        LEFT JOIN states s ON a.state_id = s.id
        WHERE a.doctor_id = $1 
        AND a.date BETWEEN $2 AND $3
        GROUP BY EXTRACT(YEAR FROM a.date), EXTRACT(MONTH FROM a.date)
        ORDER BY year, month
    `

	results, err := h.store.ExecStatementMany(ctx.Request.Context(), statement, []interface{}{doctorID, startDate, endDate})
	if err != nil {
		return nil, err
	}

	var monthlyData []MonthlyAppointmentData
	for _, result := range results {
		// Handle year conversion
		var year int64
		if yearVal, ok := result["year"]; ok && yearVal != nil {
			switch v := yearVal.(type) {
			case int64:
				year = v
			case float64:
				year = int64(v)
			case string:
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					year = parsed
				}
			}
		}

		// Handle month conversion
		var month int64
		if monthVal, ok := result["month"]; ok && monthVal != nil {
			switch v := monthVal.(type) {
			case int64:
				month = v
			case float64:
				month = int64(v)
			case string:
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					month = parsed
				}
			}
		}

		// Validate month range
		if month < 1 || month > 12 {
			continue // Skip invalid months
		}

		// Handle count conversions
		var total, completed, scheduled, cancelled int64

		if totalVal, ok := result["total"]; ok && totalVal != nil {
			if v, ok := totalVal.(int64); ok {
				total = v
			} else if v, ok := totalVal.(float64); ok {
				total = int64(v)
			}
		}

		if completedVal, ok := result["completed"]; ok && completedVal != nil {
			if v, ok := completedVal.(int64); ok {
				completed = v
			} else if v, ok := completedVal.(float64); ok {
				completed = int64(v)
			}
		}

		if scheduledVal, ok := result["scheduled"]; ok && scheduledVal != nil {
			if v, ok := scheduledVal.(int64); ok {
				scheduled = v
			} else if v, ok := scheduledVal.(float64); ok {
				scheduled = int64(v)
			}
		}

		if cancelledVal, ok := result["cancelled"]; ok && cancelledVal != nil {
			if v, ok := cancelledVal.(int64); ok {
				cancelled = v
			} else if v, ok := cancelledVal.(float64); ok {
				cancelled = int64(v)
			}
		}

		monthName := time.Month(month).String()

		monthlyData = append(monthlyData, MonthlyAppointmentData{
			Month:     monthName,
			Year:      int(year),
			Total:     total,
			Completed: completed,
			Scheduled: scheduled,
			Cancelled: cancelled,
		})
	}

	return monthlyData, nil
}

// getDoctorMonthlyPatients lấy thống kê bệnh nhân theo tháng
func (h *Handler) getDoctorMonthlyPatients(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) ([]MonthlyPatientData, error) {
	statement := `
		SELECT 
			EXTRACT(YEAR FROM a.date) as year,
			EXTRACT(MONTH FROM a.date) as month,
			COUNT(DISTINCT a.petid) as unique_patients,
			COUNT(DISTINCT CASE 
				WHEN (SELECT COUNT(*) FROM appointments a2 WHERE a2.petid = a.petid AND a2.date < a.date) = 0 
				THEN a.petid 
			END) as new_patients
		FROM appointments a
		WHERE a.doctor_id = $1 
		AND a.date BETWEEN $2 AND $3
		GROUP BY EXTRACT(YEAR FROM a.date), EXTRACT(MONTH FROM a.date)
		ORDER BY year, month
	`

	results, err := h.store.ExecStatementMany(ctx.Request.Context(), statement, []interface{}{doctorID, startDate, endDate})
	if err != nil {
		return nil, err
	}

	var monthlyData []MonthlyPatientData
	for _, result := range results {
		year, _ := result["year"].(int64)
		month, _ := result["month"].(int64)
		uniquePatients, _ := result["unique_patients"].(int64)
		newPatients, _ := result["new_patients"].(int64)
		returnPatients := uniquePatients - newPatients

		monthName := time.Month(month).String()

		monthlyData = append(monthlyData, MonthlyPatientData{
			Month:          monthName,
			Year:           int(year),
			UniquePatients: uniquePatients,
			NewPatients:    newPatients,
			ReturnPatients: returnPatients,
		})
	}

	return monthlyData, nil
}

// getDoctorWorkingDays đếm số ngày làm việc của bác sĩ
func (h *Handler) getDoctorWorkingDays(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) (int64, error) {
	statement := `
		SELECT COUNT(DISTINCT DATE(s.date)) as working_days
		FROM shifts s
		WHERE s.doctor_id = $1 
		AND s.date BETWEEN $2 AND $3
	`

	result, err := h.store.ExecStatementOne(ctx, statement, []interface{}{doctorID, startDate, endDate}, []string{"working_days"})
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var count int64
	if value, ok := result["working_days"]; ok {
		count, _ = value.(int64)
	}

	return count, nil
}

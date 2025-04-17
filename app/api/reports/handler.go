package reports

import (
	"database/sql"
	"math"
	"math/big"
	"net/http"
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
		FROM appointments
		WHERE date BETWEEN $1 AND $2
		AND state_id = (SELECT id FROM states WHERE state = 'Completed')
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
				mh.condition as disease,
				COUNT(*) as count
			FROM 
				medical_history mh
			JOIN 
				medical_records mr ON mh.medical_record_id = mr.id
			WHERE 
				mh.diagnosis_date BETWEEN $1 AND $2
			GROUP BY 
				mh.condition
			ORDER BY 
				count DESC
			LIMIT 10
		),
		total AS (
			SELECT COUNT(*) as total 
			FROM medical_history mh
			WHERE mh.diagnosis_date BETWEEN $1 AND $2
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
				date_trunc('month', a.date) as month,
				COUNT(*) as count
			FROM 
				appointments a
			WHERE 
				a.date BETWEEN $1 AND $2
				AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
			GROUP BY 
				date_trunc('month', a.date)
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

// Performance reports
type PerformanceReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type PerformanceReportResponse struct {
	StaffPerformance     []StaffPerformanceStats  `json:"staff_performance"`
	AverageTimeStats     TimeStats                `json:"average_time_stats"`
	ServiceEfficiency    []ServiceEfficiencyStats `json:"service_efficiency"`
	CustomerSatisfaction float64                  `json:"customer_satisfaction"`
}

type StaffPerformanceStats struct {
	DoctorID     int64   `json:"doctor_id"`
	DoctorName   string  `json:"doctor_name"`
	Appointments int     `json:"appointments"`
	AvgDuration  float64 `json:"avg_duration_minutes"`
	Revenue      float64 `json:"revenue"`
}

type TimeStats struct {
	AvgWaitingTime     float64 `json:"avg_waiting_time_minutes"`
	AvgExaminationTime float64 `json:"avg_examination_time_minutes"`
	AvgTotalTime       float64 `json:"avg_total_time_minutes"`
}

type ServiceEfficiencyStats struct {
	ServiceID        int64   `json:"service_id"`
	ServiceName      string  `json:"service_name"`
	AvgDuration      float64 `json:"avg_duration_minutes"`
	StandardDuration int     `json:"standard_duration"`
	Efficiency       float64 `json:"efficiency_percentage"`
}

func (h *Handler) GetPerformanceReport(ctx *gin.Context) {
	var req PerformanceReportRequest
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

	// Get staff performance statistics
	staffPerformance, err := h.getStaffPerformance(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff performance"})
		return
	}

	// Get average time statistics
	timeStats, err := h.getAverageTimeStats(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get average time statistics"})
		return
	}

	// Get service efficiency
	serviceEfficiency, err := h.getServiceEfficiency(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get service efficiency"})
		return
	}

	// Get customer satisfaction rating
	satisfaction, err := h.getCustomerSatisfaction(ctx, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer satisfaction"})
		return
	}

	response := PerformanceReportResponse{
		StaffPerformance:     staffPerformance,
		AverageTimeStats:     timeStats,
		ServiceEfficiency:    serviceEfficiency,
		CustomerSatisfaction: satisfaction,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) getStaffPerformance(ctx *gin.Context, startDate, endDate time.Time) ([]StaffPerformanceStats, error) {
	var stats []StaffPerformanceStats

	results, err := h.store.ExecStatementMany(ctx, `
		SELECT 
			d.user_id as doctor_id,
			u.full_name as doctor_name,
			COUNT(a.id) as appointments,
			AVG(EXTRACT(EPOCH FROM (c.end_time - c.start_time)) / 60) as avg_duration_minutes,
			COALESCE(SUM(p.amount), 0) as revenue
		FROM 
			doctors d
		JOIN 
			users u ON d.user_id = u.id
		LEFT JOIN 
			appointments a ON a.doctor_id = d.user_id
		LEFT JOIN 
			consultations c ON c.appointment_id = a.id
		LEFT JOIN 
			payments p ON p.appointment_id = a.id AND p.payment_status = 'completed'
		WHERE 
			a.date BETWEEN $1 AND $2
			AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
		GROUP BY 
			d.user_id, u.full_name
		ORDER BY 
			appointments DESC
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var s StaffPerformanceStats

		if doctorID, ok := result["doctor_id"].(int64); ok {
			s.DoctorID = doctorID
		}
		if doctorName, ok := result["doctor_name"].(string); ok {
			s.DoctorName = doctorName
		}
		if appointments, ok := result["appointments"].(int); ok {
			s.Appointments = appointments
		}
		if avgDuration, ok := result["avg_duration_minutes"].(float64); ok {
			s.AvgDuration = avgDuration
		}
		if revenue, ok := result["revenue"].(float64); ok {
			s.Revenue = revenue
		}

		stats = append(stats, s)
	}

	return stats, nil
}

func (h *Handler) getAverageTimeStats(ctx *gin.Context, startDate, endDate time.Time) (TimeStats, error) {
	var stats TimeStats

	results, err := h.store.ExecStatementMany(ctx, `
		WITH appointment_times AS (
			SELECT 
				a.id,
				EXTRACT(EPOCH FROM (a.checked_in_time - a.scheduled_time)) / 60 as waiting_time,
				EXTRACT(EPOCH FROM (c.end_time - c.start_time)) / 60 as examination_time,
				EXTRACT(EPOCH FROM (c.end_time - a.scheduled_time)) / 60 as total_time
			FROM 
				appointments a
			JOIN 
				consultations c ON c.appointment_id = a.id
			WHERE 
				a.date BETWEEN $1 AND $2
				AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
				AND a.checked_in_time IS NOT NULL
		)
		SELECT 
			AVG(waiting_time) as avg_waiting_time,
			AVG(examination_time) as avg_examination_time,
			AVG(total_time) as avg_total_time
		FROM 
			appointment_times
	`, []interface{}{startDate, endDate})

	if err != nil {
		return TimeStats{}, err
	}

	for _, result := range results {
		if avgWaitingTime, ok := result["avg_waiting_time"].(float64); ok {
			stats.AvgWaitingTime = avgWaitingTime
		}
		if avgExaminationTime, ok := result["avg_examination_time"].(float64); ok {
			stats.AvgExaminationTime = avgExaminationTime
		}
		if avgTotalTime, ok := result["avg_total_time"].(float64); ok {
			stats.AvgTotalTime = avgTotalTime
		}
	}

	return stats, nil
}

func (h *Handler) getServiceEfficiency(ctx *gin.Context, startDate, endDate time.Time) ([]ServiceEfficiencyStats, error) {
	var stats []ServiceEfficiencyStats

	results, err := h.store.ExecStatementMany(ctx, `
		SELECT 
			s.id as service_id,
			s.name as service_name,
			AVG(EXTRACT(EPOCH FROM (c.end_time - c.start_time)) / 60) as avg_duration_minutes,
			s.duration as standard_duration,
			(s.duration / NULLIF(AVG(EXTRACT(EPOCH FROM (c.end_time - c.start_time)) / 60), 0)) * 100 as efficiency_percentage
		FROM 
			services s
		JOIN 
			appointments a ON a.service_id = s.id
		JOIN 
			consultations c ON c.appointment_id = a.id
		WHERE 
			a.date BETWEEN $1 AND $2
			AND a.state_id = (SELECT id FROM states WHERE state = 'Completed')
		GROUP BY 
			s.id, s.name, s.duration
		ORDER BY 
			efficiency_percentage DESC
	`, []interface{}{startDate, endDate})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var s ServiceEfficiencyStats

		if serviceID, ok := result["service_id"].(int64); ok {
			s.ServiceID = serviceID
		}
		if serviceName, ok := result["service_name"].(string); ok {
			s.ServiceName = serviceName
		}
		if avgDuration, ok := result["avg_duration_minutes"].(float64); ok {
			s.AvgDuration = avgDuration
		}

		stats = append(stats, s)
	}

	return stats, nil
}

func (h *Handler) getCustomerSatisfaction(ctx *gin.Context, startDate, endDate time.Time) (float64, error) {
	var satisfaction float64

	// This assumes you have a rating field in consultations or a separate feedback table
	// Adjust according to your actual schema
	results, err := h.store.ExecStatementMany(ctx, `
		SELECT 
			COALESCE(AVG(c.rating), 0) as avg_rating
		FROM 
			consultations c
		JOIN 
			appointments a ON c.appointment_id = a.id
		WHERE 
			a.date BETWEEN $1 AND $2
			AND c.rating IS NOT NULL
	`, []interface{}{startDate, endDate})

	if err != nil {
		return 0, err
	}

	for _, result := range results {
		if avgRating, ok := result["avg_rating"].(float64); ok {
			satisfaction = avgRating
		}
	}

	return satisfaction, nil
}

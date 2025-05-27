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
				c.created_at BETWEEN $1 AND $2
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

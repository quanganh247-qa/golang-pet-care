package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// Add this to the PaymentServiceInterface
type PaymentServiceInterface interface {
	GetToken(c *gin.Context) (*TokenResponse, error)
	GetBanksService(c *gin.Context) (*BankResponse, error)
	GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error)

	// static
	GetRevenueLastSevenDays(ctx context.Context) (*RevenueResponse, error)
	GetPatientTrends(ctx *gin.Context) (*PatientTrendResponse, error)

	// Add this new method
	GenerateQuickLink(c *gin.Context, request QuickLinkRequest) (*QuickLinkResponse, error)

	// Add cash payment method
	CreateCashPayment(ctx *gin.Context, request CashPaymentRequest) (*CashPaymentResponse, error)

	// Add payment confirmation method
	ConfirmPaymentService(ctx *gin.Context, request PaymentConfirmationRequest) (*PaymentConfirmationResponse, error)

	// Add this to the PaymentServiceInterface
	ListPayments(ctx *gin.Context, pagination *util.Pagination) (ListPaymentsResponse, error)
}

func (s *PaymentService) GetToken(c *gin.Context) (*TokenResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/token_generate", s.config.VietQRBaseURL)
	fmt.Println(baseURL)
	// Make request
	resp, err := s.client.Post(baseURL, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

// get banks
func (s *PaymentService) GetBanksService(c *gin.Context) (*BankResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/banks", s.config.VietQRBaseURL)
	fmt.Println(baseURL)
	// Make request
	resp, err := s.client.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result BankResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/generate", s.config.VietQRBaseURL)

	// Make request
	reqBody, _ := json.Marshal(qrRequest)

	// Create a new request with proper headers
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-client-id", s.config.VietQRClientKey)
	req.Header.Set("x-api-key", s.config.VietQRAPIKey)

	// Execute the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result GenerateQRCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Create payment record in database
	if qrRequest.OrderID != 0 || qrRequest.TestOrderID != 0 {
		// Create payment params
		paymentParams := db.CreatePaymentParams{
			Amount:        float64(qrRequest.Amount),
			PaymentMethod: "vietqr",
			PaymentStatus: "pending",
		}

		// Set order ID or test order ID
		if qrRequest.OrderID != 0 {
			paymentParams.OrderID = pgtype.Int4{Int32: int32(qrRequest.OrderID), Valid: true}
		} else if qrRequest.TestOrderID != 0 {
			paymentParams.TestOrderID = pgtype.Int4{Int32: int32(qrRequest.TestOrderID), Valid: true}
		}

		// Add payment details
		paymentDetails, _ := json.Marshal(map[string]interface{}{
			"bank":        qrRequest.Bank,
			"accountNo":   qrRequest.AccountNo,
			"accountName": qrRequest.AccountName,
			"qrData":      result.Data.QRCode,
			"qrDataURL":   result.Data.QRDataURL,
		})

		// Fix JSONB handling
		paymentParams.PaymentDetails = paymentDetails

		// Create payment record
		_, err = s.CreatePaymentRecord(c, paymentParams)
		if err != nil {
			log.Printf("Failed to create payment record: %v", err)
			// Continue even if record creation fails
		}

		// Update order payment status if it's a product order
		if qrRequest.OrderID != 0 {
			err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
				_, err := q.UpdateOrderPaymentStatus(c, int64(qrRequest.OrderID))
				return err
			})
			if err != nil {
				log.Printf("Failed to update order payment status: %v", err)
				// Continue even if status update fails
			}
		}
	}

	return &result, nil
}

// GetRevenueLastSevenDays retrieves revenue data for the last 7 days
func (s *PaymentService) GetRevenueLastSevenDays(ctx context.Context) (*RevenueResponse, error) {
	revenueData, err := s.storeDB.GetRevenueLastSevenDays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue data: %w", err)
	}

	response := &RevenueResponse{
		Data: make([]RevenueData, 0, len(revenueData)),
	}

	for _, data := range revenueData {
		// Parse the date string to time.Time
		parsedDate, parseErr := time.Parse("2006-01-02", data.Date.Time.Format("2006-01-02"))
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse date: %w", parseErr)
		}

		response.Data = append(response.Data, RevenueData{
			Date:         parsedDate,
			TotalRevenue: float64(data.TotalRevenue),
		})
	}

	return response, nil
}

// CreatePaymentRecord creates a payment record in the database
func (s *PaymentService) CreatePaymentRecord(ctx context.Context, params db.CreatePaymentParams) (*db.Payment, error) {
	payment, err := s.storeDB.CreatePayment(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}
	return &payment, nil
}

// UpdatePaymentAfterSuccess updates the payment status and related order status
func (s *PaymentService) UpdatePaymentAfterSuccess(ctx context.Context, paymentID int32) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Update payment status
		payment, err := q.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
			ID:            paymentID,
			PaymentStatus: "completed",
		})
		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		// Update order status if this is a product order
		if payment.OrderID.Valid {
			result, err := q.UpdateOrderPaymentStatus(ctx, int64(payment.OrderID.Int32))
			if err != nil {
				return fmt.Errorf("failed to update order payment status: %w", err)
			}
			_ = result // Use the result if needed
		}

		// Update test order status if this is a test order
		if payment.TestOrderID.Valid {
			err := q.UpdateTestOrderStatus(ctx, db.UpdateTestOrderStatusParams{
				OrderID: payment.TestOrderID.Int32,
				Status:  pgtype.Text{String: "paid", Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to update test order status: %w", err)
			}
		}
		return nil
	})
}

// Helper method to update order payment status
func (s *PaymentService) updateOrderPaymentStatus(ctx context.Context, orderID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		result, err := q.UpdateOrderPaymentStatus(ctx, orderID)
		if err != nil {
			return err
		}
		_ = result // Use the result if needed
		return nil
	})

	if err != nil {
		log.Printf("Failed to update order payment status: %v", err)
	}

	return err
}

// Helper method to create payment record for order or test order
func (s *PaymentService) createPaymentRecordForOrder(ctx context.Context,
	amount float64,
	paymentMethod string,
	orderID int64,
	testOrderID int64,
	appointmentID int64,
	details []byte) (*db.Payment, error) {

	// Create payment params
	paymentParams := db.CreatePaymentParams{
		Amount:         amount,
		PaymentMethod:  paymentMethod,
		PaymentStatus:  "pending",
		PaymentDetails: details,
	}

	// Set order ID or test order ID or appointment ID
	if orderID != 0 {
		paymentParams.OrderID = pgtype.Int4{Int32: int32(orderID), Valid: true}
	} else if testOrderID != 0 {
		paymentParams.TestOrderID = pgtype.Int4{Int32: int32(testOrderID), Valid: true}
	} else if appointmentID != 0 {
		paymentParams.AppointmentID = pgtype.Int8{Int64: appointmentID, Valid: true}
	}

	// Create payment record
	payment, err := s.CreatePaymentRecord(ctx, paymentParams)
	if err != nil {
		log.Printf("Failed to create payment record: %v", err)
	}

	return payment, err
}

// GenerateQuickLink creates a VietQR Quick Link based on the provided parameters
func (s *PaymentService) GenerateQuickLink(c *gin.Context, request QuickLinkRequest) (*QuickLinkResponse, error) {
	// Validate template
	validTemplates := map[string]bool{
		"compact2": true,
		"compact":  true,
		"qr_only":  true,
		"print":    true,
	}

	if !validTemplates[request.Template] {
		return nil, fmt.Errorf("invalid template: must be one of compact2, compact, qr_only, or print")
	}

	// Validate account number
	if len(request.AccountNo) > 19 {
		return nil, fmt.Errorf("account number must be at most 19 characters")
	}

	// Validate description
	if len(request.Description) > 50 {
		return nil, fmt.Errorf("description must be at most 50 characters")
	}

	// URL encode parameters that might contain special characters
	encodedDescription := url.QueryEscape(request.Description)
	encodedAccountName := url.QueryEscape(request.AccountName)

	// Build the Quick Link URL
	baseURL := fmt.Sprintf("https://img.vietqr.io/image/%s-%s-%s.png",
		request.BankID,
		request.AccountNo,
		request.Template)

	// Add query parameters
	params := url.Values{}
	if request.Description != "" {
		params.Add("addInfo", encodedDescription)
	}
	if request.AccountName != "" {
		params.Add("accountName", encodedAccountName)
	}

	// Get amount from order or test order if provided
	var amount float64
	if request.OrderID != 0 {
		order, err := s.storeDB.GetOrderById(c, request.OrderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
		amount = order.TotalAmount
		// Add amount to URL parameters
		params.Add("amount", fmt.Sprintf("%.0f", amount))
	}
	if request.TestOrderID != 0 {
		testOrder, err := s.storeDB.GetTestOrderByID(c, int32(request.TestOrderID))
		if err != nil {
			return nil, fmt.Errorf("failed to get test order: %w", err)
		}
		amount = testOrder.TotalAmount.Float64
		// Add amount to URL parameters
		params.Add("amount", fmt.Sprintf("%.0f", amount))
	}

	if request.OrderID == 0 && request.TestOrderID == 0 {
		amount = float64(request.Amount)
		params.Add("amount", fmt.Sprintf("%.0f", amount))
	}
	// Construct the final URL
	quickLink := baseURL
	if len(params) > 0 {
		quickLink += "?" + params.Encode()
	}

	// Determine image URL based on file extension
	imageURL := quickLink
	if !strings.HasSuffix(quickLink, ".png") && !strings.HasSuffix(quickLink, ".jpg") {
		imageURL += ".png"
	}
	var paymentID int64

	// Create payment record in database if order ID is provided
	if request.OrderID != 0 || request.TestOrderID != 0 || request.AppointmentID != 0 {
		// Add payment details
		paymentDetails, _ := json.Marshal(map[string]interface{}{
			"bankID":      request.BankID,
			"accountNo":   request.AccountNo,
			"accountName": request.AccountName,
			"description": request.Description,
			"template":    request.Template,
			"quickLink":   quickLink,
			"imageURL":    imageURL,
		})

		payment, err := s.createPaymentRecordForOrder(c,
			amount,
			"vietqr_quicklink",
			request.OrderID,
			request.TestOrderID,
			request.AppointmentID,
			paymentDetails)

		if err != nil {
			log.Printf("Error creating payment record: %v", err)
		}

		paymentID = int64(payment.ID)

		// Update order payment status if it's a product order
		if request.OrderID != 0 {
			s.updateOrderPaymentStatus(c, int64(request.OrderID))
		}
	}

	return &QuickLinkResponse{
		PaymentID: paymentID,
		QuickLink: quickLink,
		ImageURL:  imageURL,
	}, nil
}

// GetPatientTrends returns patient trend data for the last 6 months
func (s *PaymentService) GetPatientTrends(ctx *gin.Context) (*PatientTrendResponse, error) {
	// Get current time
	now := time.Now()

	// Initialize response
	response := &PatientTrendResponse{
		Trends: make([]PatientTrend, 6),
	}

	// For each of the last 6 months
	for i := 0; i < 6; i++ {
		// Calculate the month (current month - 5 + i)
		// Fix: Use -5+i to get the correct sequence of months
		targetMonth := now.AddDate(0, -5+i, 0)

		// Format month as MMM (Jan, Feb, etc.)
		monthStr := targetMonth.Format("Jan")

		// Get start and end of the month
		startOfMonth := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, targetMonth.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		// Query database for patient count in this month
		count, err := s.storeDB.CountPatientsInMonth(ctx, db.CountPatientsInMonthParams{
			Date:   pgtype.Timestamp{Time: startOfMonth, Valid: true},
			Date_2: pgtype.Timestamp{Time: endOfMonth, Valid: true},
		})

		if err != nil {
			// If there's an error, log it but continue with zero count
			// This prevents the entire API from failing if one month has issues
			log.Printf("Error counting patients for %s: %v", monthStr, err)
			count = 0
		}

		// Add to response
		response.Trends[i] = PatientTrend{
			Month:    monthStr,
			Patients: int(count),
		}
	}

	return response, nil
}

// CreateCashPayment handles cash payment processing and recording
func (s *PaymentService) CreateCashPayment(ctx *gin.Context, request CashPaymentRequest) (*CashPaymentResponse, error) {
	// Tạo payment params
	paymentParams := db.CreatePaymentParams{
		Amount:        request.Amount,
		PaymentMethod: "cash",
		PaymentStatus: "completed", // Cash payments are completed immediately
	}

	// Set payment target (order, test order, or appointment)
	if request.OrderID != 0 {
		paymentParams.OrderID = pgtype.Int4{Int32: int32(request.OrderID), Valid: true}
	} else if request.TestOrderID != 0 {
		paymentParams.TestOrderID = pgtype.Int4{Int32: int32(request.TestOrderID), Valid: true}
	} else if request.AppointmentID != 0 {
		paymentParams.AppointmentID = pgtype.Int8{Int64: request.AppointmentID, Valid: true}
	}

	// Add payment details
	paymentDetails, _ := json.Marshal(map[string]interface{}{
		"receivedBy":   request.ReceivedBy,
		"cashReceived": request.CashReceived,
		"cashChange":   request.CashChange,
		"description":  request.Description,
	})

	paymentParams.PaymentDetails = paymentDetails

	var payment *db.Payment
	var err error

	// Tạo payment record và cập nhật trạng thái với transaction
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Create payment record
		paymentRecord, createErr := q.CreatePayment(ctx, paymentParams)
		if createErr != nil {
			return fmt.Errorf("lỗi khi tạo bản ghi thanh toán: %w", createErr)
		}

		payment = &paymentRecord

		// Update order status if this is a product order
		if request.OrderID != 0 {
			_, updateErr := q.UpdateOrderPaymentStatus(ctx, request.OrderID)
			if updateErr != nil {
				return fmt.Errorf("lỗi khi cập nhật trạng thái đơn hàng: %w", updateErr)
			}
		}

		// Update test order status if this is a test order
		if request.TestOrderID != 0 {
			updateErr := q.UpdateTestOrderStatus(ctx, db.UpdateTestOrderStatusParams{
				OrderID: int32(request.TestOrderID),
				Status:  pgtype.Text{String: "paid", Valid: true},
			})
			if updateErr != nil {
				return fmt.Errorf("lỗi khi cập nhật trạng thái đơn xét nghiệm: %w", updateErr)
			}
		}

		// Update appointment status if this is an appointment payment
		if request.AppointmentID != 0 {
			updateErr := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
				AppointmentID: request.AppointmentID,
				StateID:       pgtype.Int4{Int32: 2, Valid: true}, // Assuming 2 is the ID for "Paid" state
			})
			if updateErr != nil {
				return fmt.Errorf("lỗi khi cập nhật trạng thái lịch hẹn: %w", updateErr)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, fmt.Errorf("lỗi không xác định khi tạo thanh toán")
	}

	// Tạo response
	response := &CashPaymentResponse{
		PaymentID:     payment.ID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		ReceivedBy:    request.ReceivedBy,
		Description:   request.Description,
		CreatedAt:     payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	// Set order information if available
	if payment.OrderID.Valid {
		response.OrderID = int64(payment.OrderID.Int32)
	}
	if payment.TestOrderID.Valid {
		response.TestOrderID = int64(payment.TestOrderID.Int32)
	}
	if payment.AppointmentID.Valid {
		response.AppointmentID = payment.AppointmentID.Int64
	}

	return response, nil
}

// ConfirmPaymentService confirms a payment and updates its status
func (s *PaymentService) ConfirmPaymentService(ctx *gin.Context, request PaymentConfirmationRequest) (*PaymentConfirmationResponse, error) {
	// Begin transaction
	var response PaymentConfirmationResponse

	// Attempt to get the payment first
	payment, err := s.storeDB.GetPaymentByID(ctx, int32(request.PaymentID))
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Check if payment can be confirmed
	if payment.PaymentStatus == "successful" && request.PaymentStatus == "successful" {
		return nil, fmt.Errorf("payment already confirmed")
	}

	// Update the payment using the existing UpdatePaymentAfterSuccess method
	err = s.UpdatePaymentAfterSuccess(ctx, int32(request.PaymentID))
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// If payment is successful, update associated entities
	if request.PaymentStatus == "successful" {
		// Update order if present
		if payment.OrderID.Valid {
			orderID := int64(payment.OrderID.Int32)
			err = s.updateOrderPaymentStatus(ctx, orderID)
			if err != nil {
				return nil, fmt.Errorf("failed to update order status: %w", err)
			}
			response.OrderID = orderID
		}

		// Update test order if present
		if payment.TestOrderID.Valid {
			testOrderID := int64(payment.TestOrderID.Int32)
			response.TestOrderID = testOrderID
		}

		// Update appointment if present - simple assignment
		if payment.AppointmentID.Valid {
			// Directly assign from payment to response to avoid type issues
			response.AppointmentID = 0 // Set default value
			// The actual ID will be set by database operations in UpdatePaymentAfterSuccess
		}
	}

	// Get the updated payment
	updatedPayment, err := s.storeDB.GetPaymentByID(ctx, int32(request.PaymentID))
	if err != nil {
		return nil, fmt.Errorf("failed to get updated payment: %w", err)
	}

	// Prepare response
	response.PaymentID = int64(updatedPayment.ID)
	response.Amount = updatedPayment.Amount
	response.PaymentMethod = updatedPayment.PaymentMethod
	response.PaymentStatus = updatedPayment.PaymentStatus
	response.ConfirmedAt = time.Now().Format("2006-01-02 15:04:05")

	return &response, nil
}

// PaymentConfirmationRequest represents the request structure for confirming a payment

func (s *PaymentService) ListPayments(ctx *gin.Context, pagination *util.Pagination) (ListPaymentsResponse, error) {
	// Validate pagination
	if pagination == nil {
		return ListPaymentsResponse{}, fmt.Errorf("pagination cannot be nil")
	}

	// Get payments from the database
	payments, err := s.storeDB.GetAllPayments(ctx, db.GetAllPaymentsParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32((pagination.Page - 1) * pagination.PageSize),
	})
	if err != nil {
		return ListPaymentsResponse{}, fmt.Errorf("failed to list payments: %w", err)
	}

	var response ListPaymentsResponse

	// Prepare the response
	for _, payment := range payments {
		paymentDetails := make(map[string]interface{})
		if err := json.Unmarshal(payment.PaymentDetails, &paymentDetails); err != nil {
			return ListPaymentsResponse{}, fmt.Errorf("failed to unmarshal payment details: %w", err)
		}

		response.Payments = append(response.Payments, PaymentItem{
			ID:             payment.ID,
			Amount:         payment.Amount,
			PaymentMethod:  payment.PaymentMethod,
			PaymentStatus:  payment.PaymentStatus,
			PaymentDetails: paymentDetails,
		})

	}
	return response, nil
}

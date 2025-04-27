package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/payOSHQ/payos-lib-golang"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// Add this to the PaymentServiceInterface
type PaymentServiceInterface interface {
	GetToken(c *gin.Context) (*TokenResponse, error)
	GetBanksService(c *gin.Context) (*BankResponse, error)
	GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error)

	GenerateOauthToken(c *gin.Context) (*OauthTokenResponse, error)

	createPayPalOrder(accessToken string, orderRequest OrderRequest) (*PayPalOrderResponse, error)
	capturePayPalOrder(accessToken string, orderID string) (*OrderCaptureResponse, error)
	getOrderDetails(accessToken string, orderID string) (map[string]interface{}, error)
	updateOrder(accessToken string, orderID string, updates []OrderUpdateRequest) (map[string]interface{}, error)
	trackOrder(accessToken string, orderID string) (map[string]interface{}, error)
	getPayPalAccessToken() (string, error)

	// payos
	createPaymentLink(ctx *gin.Context, request CreatePaymentLinkRequest) (string, error)

	// static
	GetRevenueLastSevenDays(ctx context.Context) (*RevenueResponse, error)
	GetPatientTrends(ctx *gin.Context) (*PatientTrendResponse, error)

	// Add this new method
	GenerateQuickLink(c *gin.Context, request QuickLinkRequest) (*QuickLinkResponse, error)
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

// func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
// 	// Build base URL

// 	baseURL := fmt.Sprintf("%s/generate", s.config.VietQRBaseURL)

// 	// Make request
// 	reqBody, _ := json.Marshal(qrRequest)
// 	resp, err := s.client.Post(baseURL, "application/json", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %v", err)
// 	}
// 	// Thêm các Header cần thiết
// 	resp.Header.Set("x-client-id", s.config.VietQRClientKey)
// 	resp.Header.Set("x-api-key", s.config.VietQRAPIKey)
// 	defer resp.Body.Close()

// 	// Check response status
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
// 	}

// 	// Parse response
// 	var result GenerateQRCodeResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %v", err)
// 	}

// 	if qrRequest.OrderID != 0 {
// 		err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
// 			// UpdateOrderPaymentStatus
// 			_, err := q.UpdateOrderPaymentStatus(c, int64(qrRequest.OrderID))
// 			if err != nil {
// 				return err
// 			}
// 			return nil

// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &result, nil
// }

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

// generate oauth token
func (s *PaymentService) GenerateOauthToken(c *gin.Context) (*OauthTokenResponse, error) {
	return nil, nil
}

// Create a PayPal order
func (s *PaymentService) createPayPalOrder(accessToken string, orderRequest OrderRequest) (*PayPalOrderResponse, error) {
	url := s.config.PayPalBaseURL + "/v2/checkout/orders"
	method := "POST"

	// Initialize totals
	var itemTotal, shippingCost, taxAmount, totalAmount float64

	// Parse the base amount if provided
	if orderRequest.Amount != "" {
		fmt.Sscanf(orderRequest.Amount, "%f", &totalAmount)
	}

	// Prepare purchase units
	purchaseUnit := PayPalOrderPurchaseUnit{
		Description: orderRequest.Description,
		ReferenceID: "default",
		Amount: PayPalOrderAmount{
			CurrencyCode: orderRequest.Currency,
			Value:        fmt.Sprintf("%.2f", totalAmount),
		},
	}

	// Add items if provided
	if len(orderRequest.Items) > 0 {
		// Create items array
		items := make([]PayPalOrderItem, 0, len(orderRequest.Items))

		// Calculate item total from items
		for _, item := range orderRequest.Items {
			var unitPrice float64
			fmt.Sscanf(item.UnitPrice, "%f", &unitPrice)
			itemTotal += unitPrice * float64(item.Quantity)

			// Add item to items array
			items = append(items, PayPalOrderItem{
				Name:        item.Name,
				Description: item.Description,
				Quantity:    fmt.Sprintf("%d", item.Quantity),
				UnitAmount: PayPalOrderAmount{
					CurrencyCode: orderRequest.Currency,
					Value:        item.UnitPrice,
				},
				SKU:      item.SKU,
				Category: "PHYSICAL_GOODS",
			})
		}

		// Add items to purchase unit
		purchaseUnit.Items = items

		// Parse shipping and tax costs if provided
		if orderRequest.ShippingCost != "" {
			fmt.Sscanf(orderRequest.ShippingCost, "%f", &shippingCost)
		}
		if orderRequest.TaxAmount != "" {
			fmt.Sscanf(orderRequest.TaxAmount, "%f", &taxAmount)
		}

		// Calculate the total amount properly
		totalAmount = itemTotal + shippingCost + taxAmount

		// Create amount breakdown
		purchaseUnit.Amount.Breakdown = &PayPalOrderAmountBreakdown{
			ItemTotal: &PayPalOrderAmount{
				CurrencyCode: orderRequest.Currency,
				Value:        fmt.Sprintf("%.2f", itemTotal),
			},
		}

		// Add shipping cost if provided
		if shippingCost > 0 {
			purchaseUnit.Amount.Breakdown.Shipping = &PayPalOrderAmount{
				CurrencyCode: orderRequest.Currency,
				Value:        fmt.Sprintf("%.2f", shippingCost),
			}
		}

		// Add tax amount if provided
		if taxAmount > 0 {
			purchaseUnit.Amount.Breakdown.TaxTotal = &PayPalOrderAmount{
				CurrencyCode: orderRequest.Currency,
				Value:        fmt.Sprintf("%.2f", taxAmount),
			}
		}
	}

	// Set the final amount with the calculated total
	purchaseUnit.Amount = PayPalOrderAmount{
		CurrencyCode: orderRequest.Currency,
		Value:        fmt.Sprintf("%.2f", totalAmount),
		Breakdown:    purchaseUnit.Amount.Breakdown,
	}

	// Create application context
	merchantName := "Test Store"
	if orderRequest.MerchantName != "" {
		merchantName = orderRequest.MerchantName
	}

	// Create PayPal order request
	paypalRequest := PayPalOrderCreateRequest{
		Intent:        "CAPTURE",
		PurchaseUnits: []PayPalOrderPurchaseUnit{purchaseUnit},
		ApplicationContext: PayPalApplicationContext{
			ReturnURL:          "http://localhost:3000/success",
			CancelURL:          "http://localhost:3000/cancel",
			BrandName:          merchantName,
			UserAction:         "PAY_NOW",
			ShippingPreference: "NO_SHIPPING",
		},
	}

	jsonData, err := json.Marshal(paypalRequest)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create PayPal order: %s", string(body))
	}

	var orderResponse PayPalOrderResponse
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		return nil, err
	}

	return &orderResponse, nil
}

// Capture a PayPal order (finalize the payment)
func (s *PaymentService) capturePayPalOrder(accessToken string, orderID string) (*OrderCaptureResponse, error) {
	url := s.config.PayPalBaseURL + "/v2/checkout/orders/" + orderID + "/capture"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to capture PayPal order: %s", string(body))
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, err
	}

	// Extract relevant information
	captureResponse := &OrderCaptureResponse{
		OrderID: orderID,
		Status:  responseMap["status"].(string),
	}

	return captureResponse, nil
}

// Get order details from PayPal
func (s *PaymentService) getOrderDetails(accessToken string, orderID string) (map[string]interface{}, error) {
	url := s.config.PayPalBaseURL + "/v2/checkout/orders/" + orderID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get order details: %s", string(body))
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// Update order details on PayPal
func (s *PaymentService) updateOrder(accessToken string, orderID string, updates []OrderUpdateRequest) (map[string]interface{}, error) {
	url := s.config.PayPalBaseURL + "/v2/checkout/orders/" + orderID
	method := "PATCH"

	jsonData, err := json.Marshal(updates)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to update order: %s", string(body))
	}

	// Get updated order details
	return s.getOrderDetails(accessToken, orderID)
}

// Track an order
func (s *PaymentService) trackOrder(accessToken string, orderID string) (map[string]interface{}, error) {
	url := s.config.PayPalBaseURL + "/v2/checkout/orders/" + orderID + "/track"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to track order: %s", string(body))
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// Get PayPal OAuth2 token
func (s *PaymentService) getPayPalAccessToken() (string, error) {
	url := s.config.PayPalBaseURL + "/v1/oauth2/token"
	method := "POST"

	payload := bytes.NewBuffer([]byte("grant_type=client_credentials"))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.config.PayPalClientID, s.config.PayPalClientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse PayPalTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (s *PaymentService) createPaymentLink(ctx *gin.Context, request CreatePaymentLinkRequest) (string, error) {

	var (
		orderCode   = time.Now().UnixNano() / int64(time.Millisecond)
		amount      int
		items       []payos.Item
		description string
	)

	// Xử lý trường hợp có OrderID
	if request.MobilePayment == true {

		cartItems, err := s.storeDB.GetCartItemsByUserId(ctx, request.UserID)
		if err != nil {
			return "", err
		}
		if cartItems == nil {
			return "", errors.New("cart not found")
		}
		// Tạo danh sách sản phẩm từ giỏ hàng
		items = make([]payos.Item, len(cartItems))
		for i, item := range cartItems {
			items[i] = payos.Item{
				Name:     item.ProductName,
				Quantity: int(item.Quantity.Int32),
				Price:    int(item.TotalPrice.Float64),
			}
		}

	} else { // Trường hợp không có OrderID
		if len(request.Items) == 0 {
			return "", errors.New("items are required when order ID is not provided")
		}

		amount = calculateTotalAmount(request.Items)
		items = request.Items
		description = request.Description
	}
	// Tạo checkout request
	paymentLinkRequest := payos.CheckoutRequestType{
		OrderCode:   orderCode,
		Amount:      amount,
		Description: description,
		Items:       items,
		CancelUrl:   "",
		ReturnUrl:   "",
	}

	paymentLinkResponse, err := payos.CreatePaymentLink(paymentLinkRequest)
	if err != nil {
		log.Printf("Create payment link failed: %v", err)
		return "", fmt.Errorf("payment link creation failed: %w", err)
	}

	return paymentLinkResponse.CheckoutUrl, nil
}

// Hàm tính tổng tiền từ items
func calculateTotalAmount(items []payos.Item) int {
	var total int
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return total
}

// Add this to the existing payment.service.go file

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
func (s *PaymentService) UpdatePaymentAfterSuccess(ctx context.Context, paymentID int32, transactionID string) error {
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
			_, err = q.UpdateOrderPaymentStatus(ctx, int64(payment.OrderID.Int32))
			if err != nil {
				return fmt.Errorf("failed to update order payment status: %w", err)
			}
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

		// if payment.AppointmentID.Valid {
		// 	err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
		// 		AppointmentID: payment.AppointmentID.Int64,
		// 		StateID:       pgtype.Int4{Int32: 2, Valid: true},
		// 	})
		// 	if err != nil {
		// 		return fmt.Errorf("failed to update appointment status: %w", err)
		// 	}
		// }

		return nil
	})
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

	// Create payment record in database if order ID is provided
	if request.OrderID != 0 || request.TestOrderID != 0 || request.AppointmentID != 0 {

		// Create payment params
		paymentParams := db.CreatePaymentParams{
			Amount:        float64(amount),
			PaymentMethod: "vietqr_quicklink",
			PaymentStatus: "pending",
		}

		// Set order ID or test order ID
		if request.OrderID != 0 {
			paymentParams.OrderID = pgtype.Int4{Int32: int32(request.OrderID), Valid: true}
		} else if request.TestOrderID != 0 {
			paymentParams.TestOrderID = pgtype.Int4{Int32: int32(request.TestOrderID), Valid: true}
		} else if request.AppointmentID != 0 {
			paymentParams.AppointmentID = pgtype.Int8{Int64: request.AppointmentID, Valid: true}
		}

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

		// Set payment details
		paymentParams.PaymentDetails = paymentDetails

		// Create payment record
		_, err := s.CreatePaymentRecord(c, paymentParams)
		if err != nil {
			log.Printf("Failed to create payment record: %v", err)
			// Continue even if record creation fails
		}

		// Update order payment status if it's a product order
		if request.OrderID != 0 {
			err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
				_, err := q.UpdateOrderPaymentStatus(c, int64(request.OrderID))
				return err
			})
			if err != nil {
				log.Printf("Failed to update order payment status: %v", err)
				// Continue even if status update fails
			}
		}
	}

	return &QuickLinkResponse{
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

package payment

import (
	"bytes"
	"encoding/json"
<<<<<<< HEAD
<<<<<<< HEAD
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

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
}

func (s *PaymentService) GetToken(c *gin.Context) (*TokenResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/token_generate", s.config.VietQRBaseURL)
=======
=======
	"errors"
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

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
}

func (s *PaymentService) GetToken(c *gin.Context) (*TokenResponse, error) {
	// Build base URL
<<<<<<< HEAD
<<<<<<< HEAD
	baseURL := fmt.Sprintf("%s/token_generate", s.config.BaseURL)
>>>>>>> c449ffc (feat: cart api)
=======
	baseURL := fmt.Sprintf("%s/token_generate", s.config.PaymentBaseURL)
>>>>>>> e859654 (Elastic search)
=======
	baseURL := fmt.Sprintf("%s/token_generate", s.config.VietQRBaseURL)
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
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
<<<<<<< HEAD
<<<<<<< HEAD
func (s *PaymentService) GetBanksService(c *gin.Context) (*BankResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/banks", s.config.VietQRBaseURL)
=======
func (s *VietQRService) GetBanksService(c *gin.Context) (*BankResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/banks", s.config.BaseURL)
>>>>>>> c449ffc (feat: cart api)
=======
func (s *PaymentService) GetBanksService(c *gin.Context) (*BankResponse, error) {
	// Build base URL
<<<<<<< HEAD
	baseURL := fmt.Sprintf("%s/banks", s.config.PaymentBaseURL)
>>>>>>> e859654 (Elastic search)
=======
	baseURL := fmt.Sprintf("%s/banks", s.config.VietQRBaseURL)
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
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

<<<<<<< HEAD
<<<<<<< HEAD
func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
	// Build base URL

	baseURL := fmt.Sprintf("%s/generate", s.config.VietQRBaseURL)
=======
// generate qr
func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
	// Build base URL

<<<<<<< HEAD
	baseURL := fmt.Sprintf("%s/generate", s.config.BaseURL)
>>>>>>> c449ffc (feat: cart api)
=======
	baseURL := fmt.Sprintf("%s/generate", s.config.PaymentBaseURL)
>>>>>>> e859654 (Elastic search)
=======
func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
	// Build base URL

	baseURL := fmt.Sprintf("%s/generate", s.config.VietQRBaseURL)
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)

	// Make request
	reqBody, _ := json.Marshal(qrRequest)
	resp, err := s.client.Post(baseURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	// Thêm các Header cần thiết
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	resp.Header.Set("x-client-id", s.config.VietQRClientKey)
	resp.Header.Set("x-api-key", s.config.VietQRAPIKey)
=======
	resp.Header.Set("x-client-id", s.config.ClientKey)
	resp.Header.Set("x-api-key", s.config.APIKey)
>>>>>>> c449ffc (feat: cart api)
=======
	resp.Header.Set("x-client-id", s.config.PaymentClientKey)
	resp.Header.Set("x-api-key", s.config.PaymentAPIKey)
>>>>>>> e859654 (Elastic search)
=======
	resp.Header.Set("x-client-id", s.config.VietQRClientKey)
	resp.Header.Set("x-api-key", s.config.VietQRAPIKey)
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result GenerateQRCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b0fe977 (place order and make payment)
	err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
		// UpdateOrderPaymentStatus
		_, err := q.UpdateOrderPaymentStatus(c, int64(qrRequest.OrderID))
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// generate oauth token
func (s *PaymentService) GenerateOauthToken(c *gin.Context) (*OauthTokenResponse, error) {
	return nil, nil
}
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ada3717 (Docker file)

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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)

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
<<<<<<< HEAD
=======
	return &result, nil
}
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)

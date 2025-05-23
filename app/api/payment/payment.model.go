package payment

import (
	"net/http"
	"time"

	"github.com/payOSHQ/payos-lib-golang"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PaymentConfig struct {
	VietQRAPIKey    string
	VietQRClientKey string
	VietQRBaseURL   string
}

// GoongService handles interactions with VierQR Maps API
type PaymentService struct {
	config  *PaymentConfig
	client  *http.Client
	storeDB db.Store
}

type PaymentApi struct {
	controller PaymentControllerInterface
}

type PaymentController struct {
	service PaymentServiceInterface
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type BankResponse struct {
	Code string `json:"code"` // Mã trạng thái trả về
	Desc string `json:"desc"` // Mô tả
	Data []Bank `json:"data"` // Danh sách ngân hàng
}

type Bank struct {
	ID                int    `json:"id"`                // ID của ngân hàng
	Name              string `json:"name"`              // Tên đầy đủ của ngân hàng
	Code              string `json:"code"`              // Mã ngân hàng
	Bin               string `json:"bin"`               // Mã BIN của ngân hàng
	ShortName         string `json:"shortName"`         // Tên viết tắt của ngân hàng
	Logo              string `json:"logo"`              // URL logo của ngân hàng
	TransferSupported int    `json:"transferSupported"` // Hỗ trợ chuyển khoản (1: Có, 0: Không)
	LookupSupported   int    `json:"lookupSupported"`   // Hỗ trợ tra cứu (1: Có, 0: Không)
}

type QRRequest struct {
	Amount      int64  `json:"amount"`
	AccountName string `json:"accountName"`
	AccountNo   string `json:"accountNo"`
	AcqId       string `json:"acqId"`
	Template    string `json:"template"`
	Bank        string `json:"bank"`
	AddInfo     string `json:"addInfo"`
	Format      string `json:"format"`
	OrderID     int64  `json:"order_id"`
	TestOrderID int64  `json:"test_order_id"`
	Medicine    bool   `json:"medicine"`
}

type GenerateQRCodeResponse struct {
	Code string         `json:"code"` // Mã trạng thái trả về
	Desc string         `json:"desc"` // Mô tả trạng thái
	Data GenerateQRData `json:"data"` // Dữ liệu trả về
}

type GenerateQRData struct {
	AcpID       int    `json:"acpId"`       // Mã ngân hàng
	AccountName string `json:"accountName"` // Tên tài khoản
	QRCode      string `json:"qrCode"`      // Dữ liệu QR code
	QRDataURL   string `json:"qrDataURL"`   // Dữ liệu QR code dạng base64
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Environment variables
var (
	PayPalClientID     string
	PayPalClientSecret string
	PayPalBaseURL      string
)

// PayPalTokenResponse represents the response format for OAuth token
type PayPalTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// PayPalOrderCreateRequest represents order creation request
type PayPalOrderCreateRequest struct {
	Intent             string                    `json:"intent"`
	PurchaseUnits      []PayPalOrderPurchaseUnit `json:"purchase_units"`
	ApplicationContext PayPalApplicationContext  `json:"application_context"`
}

// PayPalApplicationContext represents the application context
type PayPalApplicationContext struct {
	ReturnURL          string `json:"return_url"`
	CancelURL          string `json:"cancel_url"`
	BrandName          string `json:"brand_name,omitempty"`
	UserAction         string `json:"user_action,omitempty"`
	ShippingPreference string `json:"shipping_preference,omitempty"`
}

// PayPalOrderPurchaseUnit represents a purchase unit
type PayPalOrderPurchaseUnit struct {
	Amount      PayPalOrderAmount `json:"amount"`
	Description string            `json:"description,omitempty"`
	ReferenceID string            `json:"reference_id,omitempty"`
	Items       []PayPalOrderItem `json:"items,omitempty"`
}

// PayPalOrderItem represents an item in a purchase unit
type PayPalOrderItem struct {
	Name        string            `json:"name"`
	Quantity    string            `json:"quantity"`
	UnitAmount  PayPalOrderAmount `json:"unit_amount"`
	Description string            `json:"description,omitempty"`
	SKU         string            `json:"sku,omitempty"`
	Category    string            `json:"category,omitempty"`
}

// PayPalOrderAmount represents a monetary amount
type PayPalOrderAmount struct {
	CurrencyCode string                      `json:"currency_code"`
	Value        string                      `json:"value"`
	Breakdown    *PayPalOrderAmountBreakdown `json:"breakdown,omitempty"`
}

// PayPalOrderAmountBreakdown represents a breakdown of a monetary amount
type PayPalOrderAmountBreakdown struct {
	ItemTotal        *PayPalOrderAmount `json:"item_total,omitempty"`
	Shipping         *PayPalOrderAmount `json:"shipping,omitempty"`
	Handling         *PayPalOrderAmount `json:"handling,omitempty"`
	TaxTotal         *PayPalOrderAmount `json:"tax_total,omitempty"`
	ShippingDiscount *PayPalOrderAmount `json:"shipping_discount,omitempty"`
	Discount         *PayPalOrderAmount `json:"discount,omitempty"`
}

// PayPalOrderResponse represents order creation response
type PayPalOrderResponse struct {
	ID     string                    `json:"id"`
	Status string                    `json:"status"`
	Links  []PayPalOrderResponseLink `json:"links"`
}

// PayPalOrderResponseLink represents a HATEOAS link
type PayPalOrderResponseLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// OrderRequest represents the request from client
type OrderRequest struct {
	Amount       string      `json:"amount"`
	Currency     string      `json:"currency"`
	Description  string      `json:"description"`
	Items        []OrderItem `json:"items,omitempty"`
	ShippingCost string      `json:"shipping_cost,omitempty"`
	TaxAmount    string      `json:"tax_amount,omitempty"`
	MerchantName string      `json:"merchant_name,omitempty"`
}

// OrderItem represents an item in the order
type OrderItem struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Quantity    int    `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	SKU         string `json:"sku,omitempty"`
}

// OrderCaptureResponse represents the client response format
type OrderCaptureResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// OrderUpdateRequest represents the request to update an order
type OrderUpdateRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type CreatePaymentLinkRequest struct {
	MobilePayment bool         // option 2: Use for order_id
	Items         []payos.Item // option 1: using without order_id
	Description   string
	UserID        int64
}

// PatientTrend represents monthly patient statistics
type PatientTrend struct {
	Month    string `json:"month"`    // Month in MMM format (e.g., "Jan", "Feb")
	Patients int    `json:"patients"` // Number of patients in that month
}

// PatientTrendResponse represents the response for patient trend data
type PatientTrendResponse struct {
	Trends []PatientTrend `json:"trends"` // Array of patient trends for multiple months
}

// RevenueData represents daily revenue information
type RevenueData struct {
	Date         time.Time `json:"date"`
	TotalRevenue float64   `json:"total_revenue"`
}

// RevenueResponse represents the response for revenue data
type RevenueResponse struct {
	Data []RevenueData `json:"data"`
}

// QuickLinkRequest represents the request to generate a VietQR Quick Link
type QuickLinkRequest struct {
	BankID      string `json:"bank_id" binding:"required"`    // Mã BIN hoặc tên viết tắt của ngân hàng
	AccountNo   string `json:"account_no" binding:"required"` // Số tài khoản người nhận
	Template    string `json:"template" binding:"required"`   // Template hiển thị (compact2, compact, qr_only, print)
	Description string `json:"description,omitempty"`         // Nội dung chuyển khoản (tối đa 50 ký tự)
	AccountName string `json:"account_name,omitempty"`        // Tên người thụ hưởng
	Amount      int64  `json:"amount,omitempty"`              // Số tiền chuyển khoản (tối đa 12 chữ số)

	// Internal fields for application use
	OrderID       int64 `json:"order_id,omitempty"`       // ID đơn hàng trong hệ thống
	TestOrderID   int64 `json:"test_order_id,omitempty"`  // ID đơn hàng xét nghiệm trong hệ thống
	AppointmentID int64 `json:"appointment_id,omitempty"` // ID lịch hẹn trong hệ thống
}

// QuickLinkResponse represents the response from generating a VietQR Quick Link
type QuickLinkResponse struct {
	PaymentID int64  `json:"payment_id"` // ID của thanh toán
	QuickLink string `json:"quick_link"` // URL của Quick Link
	ImageURL  string `json:"image_url"`  // URL của ảnh QR code
}

// CashPaymentRequest represents the request to create a cash payment
type CashPaymentRequest struct {
	Amount        float64 `json:"amount" binding:"required"`        // Số tiền thanh toán
	Description   string  `json:"description,omitempty"`            // Mô tả thanh toán
	OrderID       int64   `json:"order_id,omitempty"`               // ID đơn hàng nếu thanh toán cho đơn hàng
	TestOrderID   int64   `json:"test_order_id,omitempty"`          // ID đơn xét nghiệm nếu thanh toán cho xét nghiệm
	AppointmentID int64   `json:"appointment_id,omitempty"`         // ID lịch hẹn nếu thanh toán cho lịch hẹn
	ReceivedBy    string  `json:"received_by" binding:"required"`   // Người nhận tiền
	CashReceived  float64 `json:"cash_received" binding:"required"` // Số tiền nhận được
	CashChange    float64 `json:"cash_change" binding:"omitempty"`  // Số tiền thối lại
}

// CashPaymentResponse represents the response after creating a cash payment
type CashPaymentResponse struct {
	PaymentID     int32   `json:"payment_id"`               // ID của thanh toán
	Amount        float64 `json:"amount"`                   // Số tiền thanh toán
	PaymentMethod string  `json:"payment_method"`           // Phương thức thanh toán (Cash)
	PaymentStatus string  `json:"payment_status"`           // Trạng thái thanh toán
	OrderID       int64   `json:"order_id,omitempty"`       // ID đơn hàng (nếu có)
	TestOrderID   int64   `json:"test_order_id,omitempty"`  // ID đơn xét nghiệm (nếu có)
	AppointmentID int64   `json:"appointment_id,omitempty"` // ID lịch hẹn (nếu có)
	ReceivedBy    string  `json:"received_by"`              // Người nhận tiền
	CreatedAt     string  `json:"created_at"`               // Thời gian tạo
	Description   string  `json:"description,omitempty"`    // Mô tả
}

// PaymentConfirmationRequest represents the request to confirm a payment
type PaymentConfirmationRequest struct {
	AppointmentID int64  `json:"appointment_id,omitempty"`          // ID of the appointment to confirm
	PaymentID     int64  `json:"payment_id" binding:"required"`     // ID of the payment to confirm
	PaymentStatus string `json:"payment_status" binding:"required"` // New payment status (successful, failed)
	Notes         string `json:"notes,omitempty"`                   // Optional notes about the confirmation
}

// PaymentConfirmationResponse represents the response after confirming a payment
type PaymentConfirmationResponse struct {
	PaymentID     int64   `json:"payment_id"`               // ID of the confirmed payment
	OrderID       int64   `json:"order_id,omitempty"`       // Associated order ID, if any
	TestOrderID   int64   `json:"test_order_id,omitempty"`  // Associated test order ID, if any
	AppointmentID int64   `json:"appointment_id,omitempty"` // Associated appointment ID, if any
	Amount        float64 `json:"amount"`                   // Payment amount
	PaymentMethod string  `json:"payment_method"`           // Payment method used
	PaymentStatus string  `json:"payment_status"`           // Updated payment status
	TransactionID string  `json:"transaction_id"`           // Transaction ID
	ConfirmedAt   string  `json:"confirmed_at"`             // Time of confirmation
}

// PaymentItem represents a single payment in the list
type PaymentItem struct {
	ID             int32                  `json:"id"`
	Amount         float64                `json:"amount"`
	PaymentMethod  string                 `json:"payment_method"`
	PaymentStatus  string                 `json:"payment_status"`
	OrderID        *int32                 `json:"order_id,omitempty"`
	TestOrderID    *int32                 `json:"test_order_id,omitempty"`
	AppointmentID  *int64                 `json:"appointment_id,omitempty"`
	TransactionID  string                 `json:"transaction_id,omitempty"`
	PaymentDetails map[string]interface{} `json:"payment_details,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// ListPaymentsResponse represents the paginated response for listing payments
type ListPaymentsResponse struct {
	Payments []PaymentItem `json:"payments"`
}

type PaymentResponse struct {
	PaymentID     int64   `json:"payment_id"`               // ID của thanh toán
	Amount        float64 `json:"amount"`                   // Số tiền thanh toán
	PaymentMethod string  `json:"payment_method"`           // Phương thức thanh toán
	PaymentStatus string  `json:"payment_status"`           // Trạng thái thanh toán
	OrderID       int64   `json:"order_id,omitempty"`       // ID đơn hàng (nếu có)
	TestOrderID   int64   `json:"test_order_id,omitempty"`  // ID đơn xét nghiệm (nếu có)
	AppointmentID int64   `json:"appointment_id,omitempty"` // ID lịch hẹn (nếu có)
	ReceivedBy    string  `json:"received_by"`              // Người nhận tiền
	CreatedAt     string  `json:"created_at"`               // Thời gian tạo
	Description   string  `json:"description,omitempty"`    // Mô tả
	TransactionID string  `json:"transaction_id"`           // ID giao dịch
	ConfirmedAt   string  `json:"confirmed_at"`             // Thời gian xác nhận
	Notes         string  `json:"notes,omitempty"`          // Ghi chú xác nhận
}

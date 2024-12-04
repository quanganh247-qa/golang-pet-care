package payment

<<<<<<< HEAD
import (
	"net/http"

	"github.com/payOSHQ/payos-lib-golang"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PaymentConfig struct {
	VietQRAPIKey    string
	VietQRClientKey string
	VietQRBaseURL   string

	PayPalClientID     string
	PayPalClientSecret string
	PayPalBaseURL      string
	PayPalEnvironment  string

	PayOSAPIKey     string
	PayOSClientKey  string
	PayOSChecsumKey string
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
=======
import "net/http"

// GoongConfig contains configuration for VierQR Maps API
type VietQRConfig struct {
	APIKey    string
	ClientKey string
	BaseURL   string
}

// GoongService handles interactions with VierQR Maps API
type VietQRService struct {
	config *VietQRConfig
	client *http.Client
}

type VietQRApi struct {
	controller VietQRControllerInterface
}

type VietQRController struct {
	service VietQRServiceInterface
>>>>>>> c449ffc (feat: cart api)
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
	Amount      int    `json:"amount"`
	AccountName string `json:"accountName"`
	AccountNo   string `json:"accountNo"`
	AcqId       string `json:"acqId"`
	Template    string `json:"template"`
	Bank        string `json:"bank"`
	AddInfo     string `json:"addInfo"`
	Format      string `json:"format"`
<<<<<<< HEAD
	OrderID     int64  `json:"order_id"`
=======
>>>>>>> c449ffc (feat: cart api)
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
<<<<<<< HEAD

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
=======
>>>>>>> c449ffc (feat: cart api)

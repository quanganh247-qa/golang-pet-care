package payment

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

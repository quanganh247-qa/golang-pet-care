package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PaymentServiceInterface interface {
	GetToken(c *gin.Context) (*TokenResponse, error)
	GetBanksService(c *gin.Context) (*BankResponse, error)
	GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error)

	GenerateOauthToken(c *gin.Context) (*OauthTokenResponse, error)
}

func (s *PaymentService) GetToken(c *gin.Context) (*TokenResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/token_generate", s.config.PaymentBaseURL)
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
	baseURL := fmt.Sprintf("%s/banks", s.config.PaymentBaseURL)
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

// generate qr
func (s *PaymentService) GenerateQRService(c *gin.Context, qrRequest QRRequest) (*GenerateQRCodeResponse, error) {
	// Build base URL

	baseURL := fmt.Sprintf("%s/generate", s.config.PaymentBaseURL)

	// Make request
	reqBody, _ := json.Marshal(qrRequest)
	resp, err := s.client.Post(baseURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	// Thêm các Header cần thiết
	resp.Header.Set("x-client-id", s.config.PaymentClientKey)
	resp.Header.Set("x-api-key", s.config.PaymentAPIKey)
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

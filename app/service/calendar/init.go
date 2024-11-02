package calendar

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

var (
	conf *oauth2.Config
)

// InitConfig khởi tạo OAuth2 config
func InitConfig(config *util.Config) {
	conf = &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}
}

// GetAuthURL tạo URL để user authorize
func GetAuthURL() string {
	// Tạo state ngẫu nhiên cho bảo mật
	state := generateRandomState()
	return conf.AuthCodeURL(state)
}

// HandleCallback xử lý OAuth callback và lưu token
func HandleCallback(code string, username string) (*http.Client, error) {
	ctx := context.Background()

	// Exchange authorization code để lấy token
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	// Lưu token
	if err := saveToken(username, token); err != nil {
		return nil, fmt.Errorf("failed to save token: %v", err)
	}

	// Tạo HTTP client với token
	client := conf.Client(ctx, token)
	return client, nil
}

// GetClient lấy client từ token đã lưu hoặc tạo mới
func GetClient(token TokenInfo, username string) (*http.Client, error) {
	// convert token info to oauth2.Token
	oauth2Token := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	// Kiểm tra token hết hạn
	if oauth2Token.Expiry.Before(time.Now()) {
		// Tự động refresh token nếu có refresh token
		if token.RefreshToken != "" {
			newToken, err := conf.TokenSource(context.Background(), oauth2Token).Token()
			if err != nil {
				return nil, fmt.Errorf("failed to refresh token: %v", err)
			}
			// Lưu token mới
			if err := saveToken(username, oauth2Token); err != nil {
				return nil, fmt.Errorf("failed to save refreshed token: %v", err)
			}
			oauth2Token = newToken
		} else {
			return nil, fmt.Errorf("token expired and no refresh token available")
		}
	}

	return conf.Client(context.Background(), oauth2Token), nil
}

// saveToken lưu token vào file
func saveToken(username string, token *oauth2.Token) error {
	tokenInfo := TokenInfo{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	_, err := db.StoreDB.InsertTokenInfo(context.Background(), db.InsertTokenInfoParams{
		AccessToken:  tokenInfo.AccessToken,
		TokenType:    tokenInfo.TokenType,
		UserName:     username,
		RefreshToken: pgtype.Text{String: tokenInfo.RefreshToken, Valid: true},
		Expiry:       tokenInfo.Expiry,
	})
	if err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}

// generateRandomState tạo chuỗi ngẫu nhiên cho state
func generateRandomState() string {
	random := string(util.RandomString(16))
	return random
}

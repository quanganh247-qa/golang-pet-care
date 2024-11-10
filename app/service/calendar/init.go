package calendar

import (
	"fmt"
	"sync"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Custom errors
var (
	ErrStateExpired    = fmt.Errorf("oauth state expired")
	ErrStateInvalid    = fmt.Errorf("invalid oauth state")
	ErrTokenExpired    = fmt.Errorf("token expired")
	ErrNoRefreshToken  = fmt.Errorf("no refresh token available")
	ErrTokenSaveFailed = fmt.Errorf("failed to save token")
)

type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

type GoogleCalendarService struct {
	oauthConfig *oauth2.Config
	dbStore     db.Store
	stateTTL    time.Duration
	tokenCache  *sync.Map // cache để lưu tokens
	refreshMux  *sync.Map // mutex cho mỗi user để tránh concurrent refresh
}

type OAuthState struct {
	State     string    `json:"state"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func InitConfig(config *util.Config) *GoogleCalendarService {
	conf := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleCalendarService{
		oauthConfig: conf,
		dbStore:     db.StoreDB,
		stateTTL:    15 * time.Minute, // State expires after 15 minutes
		tokenCache:  &sync.Map{},
		refreshMux:  &sync.Map{},
	}
}

// // GetAuthURL generates OAuth URL with secure state
// func (s *GoogleCalendarService) GetAuthURL(username string) (string, error) {
// 	state := &OAuthState{
// 		State:     util.RandomString(32), // Longer state for better security
// 		Username:  username,
// 		CreatedAt: time.Now(),
// 	}

// 	// Save state to database
// 	if err := s.saveState(state); err != nil {
// 		return "", fmt.Errorf("failed to save state: %w", err)
// 	}

// 	return s.oauthConfig.AuthCodeURL(state.State), nil
// }

// // HandleCallback processes OAuth callback with state validation
// func (s *GoogleCalendarService) HandleCallback(ctx context.Context, code, state string) (*http.Client, error) {
// 	// Verify state
// 	savedState, err := s.verifyState(state)
// 	if err != nil {
// 		return nil, fmt.Errorf("state verification failed: %w", err)
// 	}

// 	// Exchange code for token
// 	token, err := s.oauthConfig.Exchange(ctx, code)
// 	if err != nil {
// 		return nil, fmt.Errorf("token exchange failed: %w", err)
// 	}

// 	// Save token
// 	if err := s.saveToken(ctx, savedState.Username, token); err != nil {
// 		return nil, fmt.Errorf("failed to save token: %w", err)
// 	}

// 	// Cache token
// 	s.cacheToken(savedState.Username, token)

// 	return s.oauthConfig.Client(ctx, token), nil
// }

// // GetClient gets or refreshes token with concurrency control
// func (s *GoogleCalendarService) GetClient(ctx context.Context, username string) (*http.Client, error) {
// 	// Try to get token from cache first
// 	if cachedToken, ok := s.getTokenFromCache(username); ok {
// 		if cachedToken.Valid() {
// 			return s.oauthConfig.Client(ctx, cachedToken), nil
// 		}
// 	}

// 	// Get token from database
// 	tokenInfo, err := s.db.GetTokenInfo(ctx, username)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get token: %w", err)
// 	}

// 	oauth2Token := &oauth2.Token{
// 		AccessToken:  tokenInfo.AccessToken,
// 		TokenType:    tokenInfo.TokenType,
// 		RefreshToken: tokenInfo.RefreshToken.String,
// 		Expiry:       tokenInfo.Expiry,
// 	}

// 	// Check if token needs refresh
// 	if !oauth2Token.Valid() {
// 		// Get or create mutex for this user
// 		mux, _ := s.refreshMux.LoadOrStore(username, &sync.Mutex{})
// 		mutex := mux.(*sync.Mutex)

// 		// Lock to prevent concurrent refresh
// 		mutex.Lock()
// 		defer mutex.Unlock()

// 		// Check again after acquiring lock
// 		if cachedToken, ok := s.getTokenFromCache(username); ok && cachedToken.Valid() {
// 			return s.oauthConfig.Client(ctx, cachedToken), nil
// 		}

// 		// Refresh token
// 		newToken, err := s.refreshToken(ctx, oauth2Token, username)
// 		if err != nil {
// 			return nil, fmt.Errorf("token refresh failed: %w", err)
// 		}
// 		oauth2Token = newToken
// 	}

// 	return s.oauthConfig.Client(ctx, oauth2Token), nil
// }

// // saveToken lưu token vào file
// func saveToken(username string, token *oauth2.Token) error {
// 	tokenInfo := TokenInfo{
// 		AccessToken:  token.AccessToken,
// 		TokenType:    token.TokenType,
// 		RefreshToken: token.RefreshToken,
// 		Expiry:       token.Expiry,
// 	}

// 	_, err := db.StoreDB.InsertTokenInfo(context.Background(), db.InsertTokenInfoParams{
// 		AccessToken:  tokenInfo.AccessToken,
// 		TokenType:    tokenInfo.TokenType,
// 		UserName:     username,
// 		RefreshToken: pgtype.Text{String: tokenInfo.RefreshToken, Valid: true},
// 		Expiry:       tokenInfo.Expiry,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to save token: %v", err)
// 	}

// 	return nil
// }

// // generateRandomState tạo chuỗi ngẫu nhiên cho state
// func generateRandomState() string {
// 	random := string(util.RandomString(16))
// 	return random
// }

// // refreshToken handles token refresh with proper error handling
// func (s *GoogleCalendarService) refreshToken(ctx context.Context, token *oauth2.Token, username string) (*oauth2.Token, error) {
// 	if token.RefreshToken == "" {
// 		return nil, ErrNoRefreshToken
// 	}

// 	tokenSource := s.oauthConfig.TokenSource(ctx, token)
// 	newToken, err := tokenSource.Token()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to refresh token: %w", err)
// 	}

// 	if err := s.saveToken(ctx, username, newToken); err != nil {
// 		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
// 	}

// 	s.cacheToken(username, newToken)
// 	return newToken, nil
// }

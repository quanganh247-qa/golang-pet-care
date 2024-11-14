package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

func (c *ClientType) LoadCacheByKey(key string, result interface{}, duration time.Duration) (cacheData func(interface{})) {
	err := c.GetWithBackground(key, result) // Attempt to load data from cache
	cacheData = func(i interface{}) {       // Define a closure to set cache
		err = c.SetWithBackground(key, i, duration) // Store new data in cache
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", key, err) // Log error if setting fails
		}
	}
	return // Return closure
}

type userInfo struct {
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Address         string `json:"address"`
	Role            string `json:"role"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	DataImage       string `json:"data_image"`
	OriginalImage   string `json:"original_image"`
}

func (c *ClientType) UserInfoLoadCache(username string) (*userInfo, error) {
	userKey := fmt.Sprintf("%s:%s", USER_INFO_KEY, username)
	userInformation := userInfo{}
	err := c.GetWithBackground(userKey, userInformation)
	if err != nil {
		log.Printf("Error when get cache for key %s: %v", userKey, err)
		userData, err := db.StoreDB.GetUser(ctxRedis, username)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("Không tìm thấy user với username = %s", username)
			}
			return nil, err
		}

		userRes := userInfo{
			UserID:        userData.ID,
			Username:      userData.Username,
			Email:         userData.Email,
			FullName:      userData.FullName,
			DataImage:     string(userData.DataImage),
			OriginalImage: userData.OriginalImage.String,
			PhoneNumber:   userData.PhoneNumber.String,
			Address:       userData.Address.String,
			Role:          userData.Role.String,
		}
		err = c.SetWithBackground(userKey, &userRes, time.Hour*12)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", userKey, err)
		}
		return &userRes, nil
	}
	return &userInformation, nil
}

// func (c *ClientType) TokenUserInfoLoadCache(username string) (*calendar.TokenInfo, error) {
// 	userKey := fmt.Sprintf("%s:%s", TOKEN_USER_INFO_KEY, username)
// 	userInformation := calendar.TokenInfo{}
// 	err := c.GetWithBackground(userKey, userInformation)
// 	if err != nil {
// 		log.Printf("Error when get cache for key %s: %v", userKey, err)
// 		tokenInfo, err := db.StoreDB.GetTokenInfo(ctxRedis, username)

// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, fmt.Errorf("Không tìm thấy user với username = %s", username)
// 			}
// 			return nil, err
// 		}

// 		userRes := calendar.TokenInfo{
// 			AccessToken:  tokenInfo.AccessToken,
// 			TokenType:    tokenInfo.TokenType,
// 			RefreshToken: tokenInfo.RefreshToken.String,
// 			Expiry:       tokenInfo.Expiry,
// 		}
// 		err = c.SetWithBackground(userKey, &userRes, time.Hour*12)
// 		if err != nil {
// 			log.Printf("Error when set cache for key %s: %v", userKey, err)
// 		}
// 		return &userRes, nil
// 	}
// 	return &userInformation, nil
// }

func (client *ClientType) RemoveUserInfoCache(username string) {
	userInfoKey := fmt.Sprintf("%s:%s", USER_INFO_KEY, username)
	client.RemoveCacheByKey(userInfoKey)
}
func (client *ClientType) RemoveTokenUserInfoCache(username string) {
	tokenUserInfoKey := fmt.Sprintf("%s:%s", TOKEN_USER_INFO_KEY, username)
	client.RemoveCacheByKey(tokenUserInfoKey)
}

func (client *ClientType) ClearUserInfoCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", USER_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}
func (client *ClientType) ClearTokenUserInfoCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", TOKEN_USER_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

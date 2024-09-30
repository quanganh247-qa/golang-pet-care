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
	Username int64  `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	PlanType string `json:"planType"`
	// Phone    string `json:"phone"`
	// EmployeeCode string          `json:"employeeCode"`
	// Avatar       string          `json:"avatar"`
	// Permissions map[string]bool `json:"permissions"`
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
		// var perms map[string]bool
		// err = json.Unmarshal([]byte(userData.))

		userRes := userInfo{
			Username: userData.ID,
			Email:    userData.Email,
			FullName: userData.FullName,
			PlanType: userData.PlanType.String,
		}
		err = c.SetWithBackground(userKey, &userRes, time.Hour*12)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", userKey, err)
		}
		return &userRes, nil
	}
	return &userInformation, nil
}

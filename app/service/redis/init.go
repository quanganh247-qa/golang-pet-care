package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/redis/go-redis/v9"
)

type ClientType struct {
	RedisClient *redis.Client
	Debug       bool // Debug flag to control verbose logging
}

var (
	Client   *ClientType
	ctxRedis = context.Background()
)

func InitRedis(config util.Config) {
	options := &redis.Options{
		Addr: config.RedisAddress,
		// Username:     config.RedisUsername,
		// Password:     config.RedisPassword,
		PoolSize:     10, // Số kết nối tối đa trong pool
		MinIdleConns: 1,  // Số kết nối nhàn rỗi tối thiểu
	}

	Client = &ClientType{
		RedisClient: redis.NewClient(options),
		Debug:       config.DebugMode, // Set debug mode based on configuration
	}
}

func (client *ClientType) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	dataValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	setValue := client.RedisClient.Set(ctx, key, dataValue, duration)
	if setValue.Err() != nil {
		return setValue.Err()
	}
	return nil
}

func (client *ClientType) Get(ctx context.Context, key string, result interface{}) error {
	valueRe := client.RedisClient.Get(ctx, key)
	if valueRe.Err() != nil {
		return valueRe.Err()
	}
	err := json.Unmarshal([]byte(valueRe.Val()), result)
	if err != nil {
		return err
	}
	return nil
}

func (client *ClientType) SetWithBackground(key string, value interface{}, duration time.Duration) error {
	dataValue, err := json.Marshal(value) // Serialize value to JSON
	if err != nil {
		if client.Debug {
			log.Printf("REDIS MARSHAL ERROR: %s - %v", key, err)
		}
		return err // Return error if serialization fails
	}
	cmdRes := client.RedisClient.Set(ctxRedis, key, dataValue, duration) // Store in Redis
	if cmdRes.Err() != nil {
		if client.Debug {
			log.Printf("REDIS SET ERROR: %s - %v", key, cmdRes.Err())
		}
		return cmdRes.Err() // Return error if Redis command fails
	}
	if client.Debug {
		log.Printf("REDIS SET: %s (TTL: %v)", key, duration)
	}
	return nil // No error, return nil
}

func (client *ClientType) GetWithBackground(key string, result interface{}) error {
	cmdRes := client.RedisClient.Get(ctxRedis, key) // Retrieve value from Redis
	if cmdRes.Err() != nil {
		if client.Debug {
			log.Printf("REDIS MISS: %s - %v", key, cmdRes.Err())
		}
		return cmdRes.Err() // Return error if Redis command fails
	}
	err := json.Unmarshal([]byte(cmdRes.Val()), result) // Deserialize JSON to result
	if err != nil {
		if client.Debug {
			log.Printf("REDIS UNMARSHAL ERROR: %s - %v", key, err)
		}
		return err // Return error if deserialization fails
	}
	if client.Debug {
		log.Printf("REDIS HIT: %s", key)
	}
	return nil // No error, return nil
}

func (client *ClientType) RemoveCacheByKey(key string) error {
	cmdRes := client.RedisClient.Del(ctxRedis, key)
	if cmdRes.Err() != nil {
		return cmdRes.Err()
	}
	return nil
}

func (client *ClientType) RemoveCacheBySubString(pattern string) error {
	cmdRes := client.RedisClient.Keys(ctxRedis, pattern)
	if cmdRes.Err() != nil {
		return cmdRes.Err()
	}
	for _, key := range cmdRes.Val() {
		err := client.RemoveCacheByKey(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *ClientType) StoreOTPInRedis(userIdentifier string, otp int64) error {
	const OTPExpiryTime = 5 * time.Minute
	otpKey := fmt.Sprintf("OTP-%s", userIdentifier)
	return client.SetWithBackground(otpKey, otp, OTPExpiryTime)
}

func (client *ClientType) ReadOTPFromRedis(userIdentifier string) (int64, error) {
	otpKey := fmt.Sprintf("OTP-%s", userIdentifier)

	var otp int64
	err := client.GetWithBackground(otpKey, &otp)
	if err != nil {

		return 0, fmt.Errorf("failed to get OTP: %w", err)
	}

	if otp == 0 {
		return 0, fmt.Errorf("invalid OTP value for %s", otpKey)
	}

	return otp, nil
}

// Delete OTP from Redis
func (client *ClientType) DeleteOTPFromRedis(userIdentifier string) error {
	otpKey := fmt.Sprintf("OTP-%s", userIdentifier)
	return client.RemoveCacheByKey(otpKey)
}

// ClearAllCache removes all keys from Redis cache
func (client *ClientType) ClearAllCache() error {
	// Get all keys using KEYS * pattern
	keys, err := client.RedisClient.Keys(ctxRedis, "*").Result()
	if err != nil {
		return fmt.Errorf("failed to get all keys: %w", err)
	}

	// If there are keys to delete, use pipeline for better performance
	if len(keys) > 0 {
		pipeline := client.RedisClient.Pipeline()
		for _, key := range keys {
			pipeline.Del(ctxRedis, key)
		}

		// Execute the pipeline commands
		_, err = pipeline.Exec(ctxRedis)
		if err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}

		if client.Debug {
			log.Printf("Cleared %d keys from Redis cache", len(keys))
		}
	}

	return nil
}

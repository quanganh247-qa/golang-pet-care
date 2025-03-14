package redis

import (
	"context"
	"encoding/json"
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"log"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 98e9e45 (ratelimit and recovery function)
=======
	"log"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 98e9e45 (ratelimit and recovery function)
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientType struct {
	RedisClient *redis.Client
}

var (
	Client   *ClientType
	ctxRedis = context.Background()
)

func InitRedis(address string) {
	Client = &ClientType{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
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
		return err // Return error if serialization fails
	}
	cmdRes := client.RedisClient.Set(ctxRedis, key, dataValue, duration) // Store in Redis
	if cmdRes.Err() != nil {
		return cmdRes.Err() // Return error if Redis command fails
	}
	return nil // No error, return nil
}

func (client *ClientType) GetWithBackground(key string, result interface{}) error {
	cmdRes := client.RedisClient.Get(ctxRedis, key) // Retrieve value from Redis
	if cmdRes.Err() != nil {
		return cmdRes.Err() // Return error if Redis command fails
	}
	err := json.Unmarshal([]byte(cmdRes.Val()), result) // Deserialize JSON to result
	if err != nil {
		return err // Return error if deserialization fails
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

func (client *ClientType) RemoveCacheBySubString(stringPattern string) error {
	cmdRes := client.RedisClient.Keys(ctxRedis, stringPattern)
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

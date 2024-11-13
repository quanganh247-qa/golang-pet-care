package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func InitRedis(address string) error {
	Client = &ClientType{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	// Test the connection
	pong, err := Client.RedisClient.Ping(ctxRedis).Result() // Fixed: Added RedisClient
	if err != nil {
		log.Fatal(err) // Handle error and exit the program
		return fmt.Errorf("cannot connect to redis: %v", err)
	}

	log.Printf("Redis connected successfully at %s (Ping response: %s)", address, pong)
	return nil
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

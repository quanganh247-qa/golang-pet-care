package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

const (
	CACHE_MIDDLEWARE_KEY = "cache:middleware"
)

// Request body reader that doesn't consume the body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheMiddleware creates a middleware that caches API responses in Redis
// cacheDuration: how long to cache the response
// keyPrefix: a prefix to use for the cache key
// methods: HTTP methods to cache (e.g. []string{"GET", "POST"})
func CacheMiddleware(cacheDuration time.Duration, keyPrefix string, methods []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache specified methods
		methodAllowed := false
		for _, m := range methods {
			if c.Request.Method == m {
				methodAllowed = true
				break
			}
		}

		if !methodAllowed {
			c.Next()
			return
		}

		// Create a unique cache key based on the full request
		key := generateCacheKey(c, keyPrefix)

		// Try to get response from cache
		var cachedResponse string
		err := redis.Client.GetWithBackground(key, &cachedResponse)
		if err == nil {
			// Cache hit
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(cachedResponse))
			c.Abort()
			return
		}

		// Cache miss, capture the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// Store response in cache if status is 200 OK
		if c.Writer.Status() == http.StatusOK {
			responseBody := blw.body.String()
			redis.Client.SetWithBackground(key, responseBody, cacheDuration)
			c.Header("X-Cache", "MISS")
		}
	}
}

// Generate a unique cache key based on the request
func generateCacheKey(c *gin.Context, keyPrefix string) string {
	// Start with the key prefix and URL path
	path := c.Request.URL.Path
	url := c.Request.URL.String() // includes query params

	var dataToHash string

	if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
		// For methods with body, include the request body in the key
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		// Restore the body for later use
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Add body to the hash
		dataToHash = fmt.Sprintf("%s:%s:%s", keyPrefix, url, string(bodyBytes))
	} else {
		// For GET requests, just use the URL (which includes query params)
		dataToHash = fmt.Sprintf("%s:%s", keyPrefix, url)
	}

	// Add user identifier for personalized content (if authenticated)
	if userID, exists := c.Get("user_id"); exists {
		dataToHash = fmt.Sprintf("%s:%v", dataToHash, userID)
	}

	// Create hash for the key
	hash := sha256.Sum256([]byte(dataToHash))
	hashString := hex.EncodeToString(hash[:])

	// Format: cache:middleware:prefix:path:hash
	pathKey := strings.ReplaceAll(path, "/", ":")
	return fmt.Sprintf("%s:%s:%s:%s", CACHE_MIDDLEWARE_KEY, keyPrefix, pathKey, hashString)
}

// InvalidateCache invalidates cache for a specific prefix
func InvalidateCache(keyPrefix string) {
	pattern := fmt.Sprintf("%s:%s:*", CACHE_MIDDLEWARE_KEY, keyPrefix)
	redis.Client.RemoveCacheBySubString(pattern)
}

package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func (s *UserService) checkRateLimit(ctx *gin.Context, userIdentifier string) error {
// 	const MaxRequestsPerSecond = 0
// 	rateLimitkey := fmt.Sprintf(constants.OTP_REQUEST_RATE, userIdentifier)
// 	rate, err := s.redis.Incr(ctx, rateLimitKey).Result()
// 	if err != nil {
// 		// s.logger.ErrorWithContext(ctx, "[ERROR] incrementing OTP request rate", zap.Error(err))
// 		return err
// 	}
// 	if rate > MaxRequestsPerSecond {
// 		// s.logger.InfoWithContext(ctx, "User has exceeded the maximum requests per second")
// 		return errors.New("You have exceeded the maximum requests per second.")
// 	}
// 	s.redis.Expire(ctx, rateLimitKey, time.Second)
// 	return nil
// // }
// func (s *UserService) incrementOTPRequestCount(ctx *gin.Context, userIdentifier string) error {
// 	const (
// 		MaxUnverifiedOTPs   = 5
// 		OTPRequestBlockTime = 10 * time.Second
// 	)
// 	requestKey := fmt.Sprintf("OTP_REQUEST_COUNT:%s", userIdentifier)
// 	requests, err := s.redis.IncrBy(ctx, requestKey, 1).Result()
// 	if err != nil {
// 		// s.logger.ErrorWithContext(ctx, "[ERROR] incrementing OTP request count", zap.Error(err))
// 		return err
// 	}
// 	if requests >= MaxUnverifiedOTPs {
// 		err = s.redis.Set(ctx, requestKey, "0", OTPRequestBlockTime).Err()
// 		if err != nil {
// 			// s.logger.ErrorWithContext(ctx, "[ERROR] resetting request count", zap.Error(err))
// 			return err
// 		}
// 		err = s.redis.Set(ctx, fmt.Sprintf("OTP_REQUEST_BLOCK:%s", userIdentifier), "1", OTPRequestBlockTime).Err()
// 		if err != nil {
// 			// s.logger.ErrorWithContext(ctx, "[ERROR] blocking user from requesting OTPs", zap.Error(err))
// 			return err
// 		}
// 		// s.logger.InfoWithContext(ctx, "User has been blocked from requesting OTPs")
// 	}
// 	return nil
// }

// func generateOTP() (string, error) {
// 	const otpLength = 6
// 	var digits = []rune("0123456789")
// 	otp := make([]rune, otpLength)
// 	for i := range otp {
// 		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
// 		if err != nil {
// 			return "", err
// 		}
// 		otp[i] = digits[num.Int64()]
// 	}
// 	return string(otp), nil
// }

func (s *UserService) storeOTPInRedis(ctx context.Context, userIdentifier string, otp int64) error {
	const OTPExpiryTime = 5 * time.Minute
	otpkey := fmt.Sprintf("OTP: %s", userIdentifier)
	// err := s.redis.Set(ctx, otpkey, otp, OTPExpiryTime)
	err := s.redis.SetWithBackground(otpkey, otp, OTPExpiryTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) readOTPFromRedis(ctx context.Context, userIdentifier string) error {
	otpKey := fmt.Sprintf("OTP:%s", userIdentifier)
	err := s.redis.GetWithBackground(otpKey, nil)
	if err != nil {
		return err
	}
	return nil
}

// func (s *UserService) sendOTPToUser(contact, otp string) error {
// 	// Implement actual SMS or email sending logic here
// 	// Placeholder implementation:
// 	fmt.Printf("Sending OTP %s to %s\n", otp, contact)
// 	return nil
// }

func (s *UserService) verifyOTP(ctx *gin.Context) {
	userIdentifier := ctx.PostForm("user_identifier")
	providedOTP := ctx.PostForm("otp")

	if userIdentifier == "" || providedOTP == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user identifier and OTP are required"})
		return
	}

	err := s.readOTPFromRedis(ctx, userIdentifier)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read OTP"})
		return
	}

	// OTP is valid, proceed with further processing
	ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

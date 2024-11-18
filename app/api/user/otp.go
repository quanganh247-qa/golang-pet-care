package user

// func (s *UserService) checkUserBlockStatus(ctx *gin.Context, userIdentifier string) error {
// 	blockKey := fmt.Sprintf(constants.OTP_REQUEST_BLOCK, userIdentifier)
// 	isBlocked, err := s.redis.Exists(ctx, blockKey).Result()
// 	if err != nil {
// 		// s.logger.ErrorWithContext(ctx, "[ERROR] checking user block status", zap.Error(err))
// 		return err
// 	}
// 	if isBlocked == 1 {
// 		// s.logger.InfoWithContext(ctx, "User is blocked from requesting OTPs")
// 		return errors.New("You are temporarily blocked due to multiple unverified OTP requests.")
// 	}
// 	return nil
// }
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

// func (s *UserService) storeOTPInRedis(ctx *gin.Context, userIdentifier string, otp string) error {
// 	const OTPExpiryTime = 5 * time.Minute
// 	otpkey := fmt.Sprintf("OTP:%s", userIdentifier)
// 	err := s.redis.Set(ctx, otpkey, otp, OTPExpiryTime)
// 	if err != nil {
// 		// s.logger.ErrorWithContext(ctx, "[ERROR] setting OTP in Redis", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }

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
// func (s *UserService) sendOTPToUser(contact, otp string) error {
// 	// Implement actual SMS or email sending logic here
// 	// Placeholder implementation:
// 	fmt.Printf("Sending OTP %s to %s\n", otp, contact)
// 	return nil
// }

// func (s *UserService) readOTPFromRedis(ctx *gin.Context, userIdentifier string) (string, error) {
// 	otpKey := fmt.Sprintf("OTP:%s", userIdentifier)
// 	otp, err := s.redis.Get(ctx, otpKey).Result()
// 	if err != nil {
// 		return "", err
// 	}
// 	return otp, nil
// }

// func (s *UserService) verifyOTP(ctx *gin.Context) {
// 	userIdentifier := ctx.PostForm("user_identifier")
// 	providedOTP := ctx.PostForm("otp")

// 	if userIdentifier == "" || providedOTP == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user identifier and OTP are required"})
// 		return
// 	}

// 	storedOTP, err := s.readOTPFromRedis(ctx, userIdentifier)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read OTP"})
// 		return
// 	}

// 	if storedOTP != providedOTP {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
// 		return
// 	}

// 	// OTP is valid, proceed with further processing
// 	ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
// }

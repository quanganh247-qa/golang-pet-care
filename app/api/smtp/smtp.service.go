package smtp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// SMTPConfigService handles SMTP configuration operations
type SMTPConfigService struct {
	config util.Config
	store  db.Store
}

// NewSMTPConfigService creates a new SMTP config service
func NewSMTPConfigService(config util.Config, store db.Store) *SMTPConfigService {
	return &SMTPConfigService{
		config: config,
		store:  store,
	}
}

// CreateSMTPConfig creates a new SMTP configuration
func (s *SMTPConfigService) CreateSMTPConfig(ctx *gin.Context) {
	var req CreateSMTPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If this is the default config, we need to update other configs
	if req.IsDefault {
		err := s.store.SetAsDefaultSMTPConfig(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default config status"})
			return
		}
	}

	// Set default values if not provided
	smtpHost := req.SMTPHost
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := req.SMTPPort
	if smtpPort == "" {
		smtpPort = "587"
	}

	// Create the SMTP config
	arg := db.CreateSMTPConfigParams{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		SmtpHost:  smtpHost,
		SmtpPort:  smtpPort,
		IsDefault: req.IsDefault,
	}

	smtpConfig, err := s.store.CreateSMTPConfig(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMTP configuration"})
		return
	}

	response := SMTPConfigResponse{
		ID:        smtpConfig.ID,
		Name:      smtpConfig.Name,
		Email:     smtpConfig.Email,
		SMTPHost:  smtpConfig.SmtpHost,
		SMTPPort:  smtpConfig.SmtpPort,
		IsDefault: smtpConfig.IsDefault,
		CreatedAt: smtpConfig.CreatedAt.Time,
		UpdatedAt: smtpConfig.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetSMTPConfig gets a SMTP configuration by ID
func (s *SMTPConfigService) GetSMTPConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	var configID int64
	fmt.Sscanf(id, "%d", &configID)

	smtpConfig, err := s.store.GetSMTPConfig(ctx, configID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}

	response := SMTPConfigResponse{
		ID:        smtpConfig.ID,
		Name:      smtpConfig.Name,
		Email:     smtpConfig.Email,
		SMTPHost:  smtpConfig.SmtpHost,
		SMTPPort:  smtpConfig.SmtpPort,
		IsDefault: smtpConfig.IsDefault,
		CreatedAt: smtpConfig.CreatedAt.Time,
		UpdatedAt: smtpConfig.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetDefaultSMTPConfig gets the default SMTP configuration
func (s *SMTPConfigService) GetDefaultSMTPConfig(ctx *gin.Context) {
	smtpConfig, err := s.store.GetDefaultSMTPConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Default SMTP configuration not found"})
		return
	}

	response := SMTPConfigResponse{
		ID:        smtpConfig.ID,
		Name:      smtpConfig.Name,
		Email:     smtpConfig.Email,
		SMTPHost:  smtpConfig.SmtpHost,
		SMTPPort:  smtpConfig.SmtpPort,
		IsDefault: smtpConfig.IsDefault,
		CreatedAt: smtpConfig.CreatedAt.Time,
		UpdatedAt: smtpConfig.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

// ListSMTPConfigs lists all SMTP configurations
func (s *SMTPConfigService) ListSMTPConfigs(ctx *gin.Context) {
	smtpConfigs, err := s.store.ListSMTPConfigs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list SMTP configurations"})
		return
	}

	var response []SMTPConfigResponse
	for _, config := range smtpConfigs {
		response = append(response, SMTPConfigResponse{
			ID:        config.ID,
			Name:      config.Name,
			Email:     config.Email,
			SMTPHost:  config.SmtpHost,
			SMTPPort:  config.SmtpPort,
			IsDefault: config.IsDefault,
			CreatedAt: config.CreatedAt.Time,
			UpdatedAt: config.UpdatedAt.Time,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateSMTPConfig updates a SMTP configuration
func (s *SMTPConfigService) UpdateSMTPConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	var configID int64
	fmt.Sscanf(id, "%d", &configID)

	var req UpdateSMTPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing config to check if we need to update the password
	existingConfig, err := s.store.GetSMTPConfig(ctx, configID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}

	// Use existing password if not provided
	password := req.Password
	if password == "" {
		password = existingConfig.Password
	}

	// If this is being set as default, update other configs
	if req.IsDefault && !existingConfig.IsDefault {
		err := s.store.SetAsDefaultSMTPConfig(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default config status"})
			return
		}
	}

	// Set default values if not provided
	smtpHost := req.SMTPHost
	if smtpHost == "" {
		smtpHost = existingConfig.SmtpHost
	}

	smtpPort := req.SMTPPort
	if smtpPort == "" {
		smtpPort = existingConfig.SmtpPort
	}

	// Update the SMTP config
	arg := db.UpdateSMTPConfigParams{
		ID:        configID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  password,
		SmtpHost:  smtpHost,
		SmtpPort:  smtpPort,
		IsDefault: req.IsDefault,
	}

	smtpConfig, err := s.store.UpdateSMTPConfig(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMTP configuration"})
		return
	}

	response := SMTPConfigResponse{
		ID:        smtpConfig.ID,
		Name:      smtpConfig.Name,
		Email:     smtpConfig.Email,
		SMTPHost:  smtpConfig.SmtpHost,
		SMTPPort:  smtpConfig.SmtpPort,
		IsDefault: smtpConfig.IsDefault,
		CreatedAt: smtpConfig.CreatedAt.Time,
		UpdatedAt: smtpConfig.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteSMTPConfig deletes a SMTP configuration
func (s *SMTPConfigService) DeleteSMTPConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	var configID int64
	fmt.Sscanf(id, "%d", &configID)

	// Check if it's the default config
	config, err := s.store.GetSMTPConfig(ctx, configID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}

	if config.IsDefault {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the default SMTP configuration"})
		return
	}

	err = s.store.DeleteSMTPConfig(ctx, configID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMTP configuration"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "SMTP configuration deleted successfully"})
}

// SetDefaultSMTPConfig sets a SMTP configuration as default
func (s *SMTPConfigService) SetDefaultSMTPConfig(ctx *gin.Context) {
	var req SetDefaultSMTPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First, reset all configs to non-default
	err := s.store.SetAsDefaultSMTPConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default config status"})
		return
	}

	// Get the existing config
	existingConfig, err := s.store.GetSMTPConfig(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}

	// Update the config to be default
	arg := db.UpdateSMTPConfigParams{
		ID:        req.ID,
		Name:      existingConfig.Name,
		Email:     existingConfig.Email,
		Password:  existingConfig.Password,
		SmtpHost:  existingConfig.SmtpHost,
		SmtpPort:  existingConfig.SmtpPort,
		IsDefault: true,
	}

	smtpConfig, err := s.store.UpdateSMTPConfig(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default SMTP configuration"})
		return
	}

	response := SMTPConfigResponse{
		ID:        smtpConfig.ID,
		Name:      smtpConfig.Name,
		Email:     smtpConfig.Email,
		SMTPHost:  smtpConfig.SmtpHost,
		SMTPPort:  smtpConfig.SmtpPort,
		IsDefault: smtpConfig.IsDefault,
		CreatedAt: smtpConfig.CreatedAt.Time,
		UpdatedAt: smtpConfig.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

// TestSMTPConfig tests a SMTP configuration by sending a test email
func (s *SMTPConfigService) TestSMTPConfig(ctx *gin.Context) {
	var req TestSMTPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing config
	existingConfig, err := s.store.GetSMTPConfig(ctx, req.SMTPId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}
	// Create a test email sender
	emailSender := mail.NewCustomSMTPSender(existingConfig.Name, existingConfig.Email, existingConfig.Password, existingConfig.SmtpHost, existingConfig.SmtpPort)

	// Try to send a test email
	subject := "SMTP Configuration Test"
	content := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body>
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #333;">SMTP Configuration Test</h2>
			<p>This is a test email to verify your SMTP configuration.</p>
			<p>If you're receiving this email, it means your SMTP configuration is working correctly.</p>
			<p>Configuration details:</p>
			<ul>
				<li>Name: %s</li>
				<li>Email: %s</li>
				<li>SMTP Host: %s</li>
				<li>SMTP Port: %s</li>
				<li>Test Time: %s</li>
			</ul>
			<p>Best regards,<br>Pet Care App Team</p>
		</div>
	</body>
	</html>`, existingConfig.Name, existingConfig.Email, existingConfig.SmtpHost, existingConfig.SmtpPort, time.Now().Format(time.RFC3339))

	to := []string{req.TestEmail}

	err = emailSender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, TestSMTPConfigResponse{
			Success: false,
			Message: fmt.Sprintf("SMTP configuration test failed: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, TestSMTPConfigResponse{
		Success: true,
		Message: "SMTP configuration test successful. Test email sent.",
	})
}

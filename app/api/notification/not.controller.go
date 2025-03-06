package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type NotificationControllerInterface interface {
	InsertNotification(ctx *gin.Context)
	ListNotificationsByUsername(ctx *gin.Context)
	DeleteNotification(ctx *gin.Context)
	SubscribeToNotifications(ctx *gin.Context)
	UnsubscribeFromNotifications(ctx *gin.Context)
	SendRealTimeNotification(ctx *gin.Context)
	MarkAsRead(ctx *gin.Context)
}

func (c *NotificationController) InsertNotification(ctx *gin.Context) {
	var arg NotificationRequest
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	notification, err := c.service.InsertNotificationService(ctx, arg, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notification)
}

func (c *NotificationController) ListNotificationsByUsername(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	notifications, err := c.service.ListNotificationsByUsernameService(ctx, authPayload.Username, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}

func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = c.service.DeleteNotificationService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

// Thêm controller để đăng ký/hủy đăng ký thông báo

func (c *NotificationController) SubscribeToNotifications(ctx *gin.Context) {
	var req SubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := ctx.GetString("username")

	err := c.service.SubscribeToNotifications(ctx, username, req.Topics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
}

func (c *NotificationController) UnsubscribeFromNotifications(ctx *gin.Context) {
	var req SubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := ctx.GetString("username")

	err := c.service.UnsubscribeFromNotifications(ctx, username, req.Topics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func (c *NotificationController) SendRealTimeNotification(ctx *gin.Context) {
	var req NotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := &Notification{
		Username:    req.Username,
		Title:       req.Title,
		Content:     req.Content,
		NotifyType:  req.NotifyType,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
	}

	err := c.service.SendRealTimeNotification(ctx, notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	notificationID := ctx.Param("id")
	notificationIDInt, err := strconv.ParseInt(notificationID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	err = c.service.MarkAsRead(ctx, notificationIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

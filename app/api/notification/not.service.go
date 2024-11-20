package notification

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type NotServiceInterface interface {
	InsertNotificationService(ctx *gin.Context, arg NotificationRequest, username string) (*NotificationResponse, error)
	ListNotificationsByUsernameService(ctx *gin.Context, username string, pagination *util.Pagination) ([]NotificationResponse, error)
	DeleteNotificationService(ctx *gin.Context, username string) error
}

func (s *NotService) InsertNotificationService(ctx *gin.Context, arg NotificationRequest, username string) (*NotificationResponse, error) {

	var res db.Notification

	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", arg.DateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DateTime: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err = q.InsertNotification(ctx, db.InsertNotificationParams{
			Username:    username,
			Title:       arg.Title,
			Description: pgtype.Text{String: arg.Description, Valid: true},
			Datetime:    pgtype.Timestamp{Time: parsedTime, Valid: true},
		})

		if err != nil {
			return fmt.Errorf("failed to insert notification: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &NotificationResponse{
		NotificationID: res.Notificationid,
		Username:       username,
		Title:          res.Title,
		Description:    res.Description.String,
		DateTime:       parsedTime,
	}, nil
}

func (s *NotService) ListNotificationsByUsernameService(ctx *gin.Context, username string, pagination *util.Pagination) ([]NotificationResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize
	notifications, err := s.storeDB.GetNotificationsByUsername(ctx, db.GetNotificationsByUsernameParams{
		Username: username,
		Limit:    int32(pagination.PageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	var res []NotificationResponse
	for _, notification := range notifications {
		res = append(res, NotificationResponse{
			NotificationID: notification.Notificationid,
			Username:       notification.Username,
			Title:          notification.Title,
			Description:    notification.Description.String,
			DateTime:       notification.Datetime.Time,
			IsRead:         notification.IsRead.Bool,
		})
	}
	return res, nil
}

// delete notification by username
func (s *NotService) DeleteNotificationService(ctx *gin.Context, username string) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeleteAllNotificationsByUser(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to delete notification: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction delete notification failed: %w", err)
	}
	return nil
}

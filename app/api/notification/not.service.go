package notification

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type NotificationServiceInterface interface {
	InsertNotificationService(ctx context.Context, arg NotificationRequest, username string) (*NotificationResponse, error)
	ListNotificationsByUsernameService(ctx context.Context, username string, pagination *util.Pagination) ([]*NotificationResponse, error)
	DeleteNotificationService(ctx context.Context, username string) error
	SendRealTimeNotification(ctx context.Context, notification *Notification) error
	SubscribeToNotifications(ctx context.Context, username string, topics []string) error
	UnsubscribeFromNotifications(ctx context.Context, username string, topics []string) error
	MarkAsRead(ctx context.Context, notificationID int64) error
}

func (s *NotificationService) InsertNotificationService(ctx context.Context, arg NotificationRequest, username string) (*NotificationResponse, error) {
	notification, err := s.storeDB.CreatetNotification(ctx, db.CreatetNotificationParams{
		Username:    username,
		Title:       arg.Title,
		Content:     pgtype.Text{String: arg.Content, Valid: true},
		NotifyType:  pgtype.Text{String: arg.NotifyType, Valid: true},
		RelatedID:   pgtype.Int4{Int32: int32(arg.RelatedID), Valid: true},
		RelatedType: pgtype.Text{String: arg.RelatedType, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if s.wsHub != nil {
		notificationJSON, _ := json.Marshal(notification)
		s.wsHub.SendToUser(username, notificationJSON)
	}

	return &NotificationResponse{
		ID:          notification.ID,
		Username:    notification.Username,
		Title:       notification.Title,
		Content:     notification.Content.String,
		NotifyType:  notification.NotifyType.String,
		RelatedID:   notification.RelatedID.Int32,
		RelatedType: notification.RelatedType.String,
		IsRead:      notification.IsRead.Bool,
		DateTime:    notification.Datetime.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *NotificationService) ListNotificationsByUsernameService(ctx context.Context, username string, pagination *util.Pagination) ([]*NotificationResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	notifications, err := s.storeDB.ListNotificationsByUsername(ctx, db.ListNotificationsByUsernameParams{
		Username: username,
		Limit:    int32(pagination.PageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var response []*NotificationResponse
	for _, n := range notifications {
		response = append(response, &NotificationResponse{
			ID:          n.ID,
			Username:    n.Username,
			Title:       n.Title,
			Content:     n.Content.String,
			NotifyType:  n.NotifyType.String,
			RelatedID:   n.RelatedID.Int32,
			RelatedType: n.RelatedType.String,
			IsRead:      n.IsRead.Bool,
			DateTime:    n.Datetime.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return response, nil
}

func (s *NotificationService) DeleteNotificationService(ctx context.Context, username string) error {
	return s.storeDB.DeleteNotificationsByUsername(ctx, username)
}

func (s *NotificationService) SendRealTimeNotification(ctx context.Context, notification *Notification) error {
	dbNotification, err := s.storeDB.CreatetNotification(ctx, db.CreatetNotificationParams{
		Username:    notification.Username,
		Title:       notification.Title,
		Content:     pgtype.Text{String: notification.Content, Valid: true},
		NotifyType:  pgtype.Text{String: notification.NotifyType, Valid: true},
		RelatedID:   pgtype.Int4{Int32: int32(notification.RelatedID), Valid: true},
		RelatedType: pgtype.Text{String: notification.RelatedType, Valid: true},
	})
	if err != nil {
		return err
	}

	if s.wsHub != nil {
		notificationJSON, _ := json.Marshal(dbNotification)
		s.wsHub.SendToUser(notification.Username, notificationJSON)
	}

	return nil
}

func (s *NotificationService) SubscribeToNotifications(ctx context.Context, username string, topics []string) error {
	for _, topic := range topics {
		err := s.storeDB.CreateNotificationPreference(ctx, db.CreateNotificationPreferenceParams{
			Username: username,
			Topic:    topic,
			Enabled:  pgtype.Bool{Bool: true, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) UnsubscribeFromNotifications(ctx context.Context, username string, topics []string) error {
	for _, topic := range topics {
		err := s.storeDB.UpdateNotificationPreference(ctx, db.UpdateNotificationPreferenceParams{
			Username: username,
			Topic:    topic,
			Enabled:  pgtype.Bool{Bool: false, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID int64) error {
	return s.storeDB.MarkNotificationAsRead(ctx, notificationID)
}

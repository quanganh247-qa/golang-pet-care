package websocket

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type NotificationService struct {
	storeDB db.Store
	ws      *WSClientManager
}

func NewNotificationService(store db.Store, ws *WSClientManager) *NotificationService {
	return &NotificationService{
		storeDB: store,
		ws:      ws,
	}
}

package rooms

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type RoomService struct {
	storeDB db.Store
}

type RoomApi struct {
	controller RoomControllerInterface
}

type RoomController struct {
	service RoomServiceInterface
}

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type UpdateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status string `json:"status" binding:"required,oneof=available occupied maintenance"`
}

type RoomResponse struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Status               string     `json:"status"`
	Type                 string     `json:"type"`
	CurrentAppointmentID *int64     `json:"current_appointment_id,omitempty"`
	AvailableAt          *time.Time `json:"available_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type ListRoomsResponse struct {
	Rooms      []RoomResponse `json:"rooms"`
	TotalCount int32          `json:"total_count"`
}

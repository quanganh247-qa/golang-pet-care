package rooms

import (
	"context"
	"fmt"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type RoomServiceInterface interface {
	CreateRoom(ctx context.Context, username string, req CreateRoomRequest) (*RoomResponse, error)
	GetRoomByID(ctx context.Context, roomID int64) (*RoomResponse, error)
	ListRooms(ctx context.Context, pagination *util.Pagination) ([]RoomResponse, error)
	UpdateRoom(ctx context.Context, username string, roomID int64, req UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, username string, roomID int64) error
}

func (s *RoomService) CreateRoom(ctx context.Context, username string, req CreateRoomRequest) (*RoomResponse, error) {
	user, err := s.storeDB.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if user.Role.String != "admin" {
		return nil, fmt.Errorf("unauthorized: only admin can create rooms")
	}

	room, err := s.storeDB.CreateRoom(ctx, db.CreateRoomParams{
		Name: req.Name,
		Type: req.Type,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %v", err)
	}

	return &RoomResponse{
		ID:     room.ID,
		Name:   room.Name,
		Status: room.Status.String,
	}, nil
}

func (s *RoomService) GetRoomByID(ctx context.Context, roomID int64) (*RoomResponse, error) {
	room, err := s.storeDB.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	return &RoomResponse{
		ID:   room.ID,
		Name: room.Name,
	}, nil
}

func (s *RoomService) ListRooms(ctx context.Context, pagination *util.Pagination) ([]RoomResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	rooms, err := s.storeDB.GetAvailableRooms(ctx, db.GetAvailableRoomsParams{
		Offset: int32(offset),
		Limit:  int32(pagination.PageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %v", err)
	}

	response := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		response[i] = RoomResponse{
			ID:   room.ID,
			Name: room.Name,
			Type: room.Type,
		}
	}

	return response, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, username string, roomID int64, req UpdateRoomRequest) error {
	user, err := s.storeDB.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.Role.String != "admin" {
		return fmt.Errorf("unauthorized: only admin can update rooms")
	}

	err = s.storeDB.UpdateRoom(ctx, db.UpdateRoomParams{
		ID:   roomID,
		Name: req.Name,
		Type: req.Type,
	})
	if err != nil {
		return fmt.Errorf("failed to update room: %v", err)
	}

	return nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, username string, roomID int64) error {
	user, err := s.storeDB.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.Role.String != "admin" {
		return fmt.Errorf("unauthorized: only admin can delete rooms")
	}

	err = s.storeDB.DeleteRoom(ctx, roomID)
	if err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}

	return nil
}

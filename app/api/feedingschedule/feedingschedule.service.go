package feedingschedule

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type FeedingScheduleServiceInterface interface {
	CreateFeedingSchedule(ctx *gin.Context, petID int64, req createFeedingScheduleRequest) (*createFeedingScheduleResponse, error)
	GetFeedingScheduleByPetID(ctx *gin.Context, petID int64) ([]createFeedingScheduleResponse, error)
	ListActiveFeedingSchedules(ctx *gin.Context) ([]createFeedingScheduleResponse, error)
	UpdateFeedingSchedule(ctx *gin.Context, id int64, req updateFeedingScheduleRequest) error
	DeleteFeedingSchedule(ctx context.Context, id int64) error
}

type FeedingScheduleService struct {
	storeDB db.Store
}

func (s *FeedingScheduleService) CreateFeedingSchedule(ctx *gin.Context, petID int64, req createFeedingScheduleRequest) (*createFeedingScheduleResponse, error) {
	res, err := s.storeDB.CreateFeedingSchedule(ctx, db.CreateFeedingScheduleParams{
		Petid:     pgtype.Int8{Int64: petID, Valid: true},
		MealTime:  pgtype.Time{Microseconds: req.MealTime.UnixNano() / 1000, Valid: true},
		FoodType:  req.FoodType,
		Quantity:  req.Quantity,
		Frequency: req.Frequency,
		Lastfed:   pgtype.Timestamp{Time: req.LastFed, Valid: true},
		Notes:     pgtype.Text{String: req.Notes, Valid: true},
		IsActive:  pgtype.Bool{Bool: req.IsActive, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create feeding schedule: %w", err)
	}

	return &createFeedingScheduleResponse{
		FeedingScheduleID: res.FeedingScheduleID,
		PetID:             res.Petid,
		MealTime:          res.MealTime,
		FoodType:          res.FoodType,
		Quantity:          res.Quantity,
		Frequency:         res.Frequency,
		LastFed:           res.Lastfed,
		Notes:             res.Notes,
		IsActive:          res.IsActive,
	}, nil
}

func (s *FeedingScheduleService) GetFeedingScheduleByPetID(ctx *gin.Context, petID int64) ([]createFeedingScheduleResponse, error) {
	res, err := s.storeDB.GetFeedingScheduleByPetID(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feeding schedules: %w", err)
	}

	var schedules []createFeedingScheduleResponse
	for _, r := range res {
		schedules = append(schedules, createFeedingScheduleResponse{
			FeedingScheduleID: r.FeedingScheduleID,
			PetID:             r.Petid,
			MealTime:          r.MealTime,
			FoodType:          r.FoodType,
			Quantity:          r.Quantity,
			Frequency:         r.Frequency,
			LastFed:           r.Lastfed,
			Notes:             r.Notes,
			IsActive:          r.IsActive,
		})
	}

	return schedules, nil
}

func (s *FeedingScheduleService) ListActiveFeedingSchedules(ctx *gin.Context) ([]createFeedingScheduleResponse, error) {
	res, err := s.storeDB.ListActiveFeedingSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active feeding schedules: %w", err)
	}

	var schedules []createFeedingScheduleResponse
	for _, r := range res {
		schedules = append(schedules, createFeedingScheduleResponse{
			FeedingScheduleID: r.FeedingScheduleID,
			PetID:             r.Petid,
			MealTime:          r.MealTime,
			FoodType:          r.FoodType,
			Quantity:          r.Quantity,
			Frequency:         r.Frequency,
			LastFed:           r.Lastfed,
			Notes:             r.Notes,
			IsActive:          r.IsActive,
		})
	}

	return schedules, nil
}

func (s *FeedingScheduleService) UpdateFeedingSchedule(ctx *gin.Context, id int64, req updateFeedingScheduleRequest) error {
	err := s.storeDB.UpdateFeedingSchedule(ctx, db.UpdateFeedingScheduleParams{
		FeedingScheduleID: id,
		MealTime:          req.MealTime,
		FoodType:          req.FoodType,
		Quantity:          req.Quantity,
		Frequency:         req.Frequency,
		Lastfed:           req.LastFed,
		Notes:             req.Notes,
		IsActive:          req.IsActive,
	})
	if err != nil {
		return fmt.Errorf("failed to update feeding schedule: %w", err)
	}
	return nil
}

func (s *FeedingScheduleService) DeleteFeedingSchedule(ctx context.Context, id int64) error {
	return s.storeDB.DeleteFeedingSchedule(ctx, id)
}

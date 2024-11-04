package feedingschedule

import "time"

type createFeedingScheduleRequest struct {
	PetID     int64     `json:"pet_id"`
	MealTime  time.Time `json:"meal_time"`
	FoodType  string    `json:"food_type"`
	Quantity  float64   `json:"quantity"`
	Frequency string    `json:"frequency"`
	LastFed   time.Time `json:"last_fed"`
	Notes     string    `json:"notes"`
	IsActive  bool      `json:"is_active"`
}

type updateFeedingScheduleRequest struct {
	FeedingScheduleID int64     `json:"feeding_schedule_id"`
	MealTime          time.Time `json:"meal_time"`
	FoodType          string    `json:"food_type"`
	Quantity          float64   `json:"quantity"`
	Frequency         string    `json:"frequency"`
	LastFed           time.Time `json:"last_fed"`
	Notes             string    `json:"notes"`
	IsActive          bool      `json:"is_active"`
}

type createFeedingScheduleResponse struct {
	FeedingScheduleID int64     `json:"feeding_schedule_id"`
	PetID             int64     `json:"pet_id"`
	MealTime          time.Time `json:"meal_time"`
	FoodType          string    `json:"food_type"`
	Quantity          float64   `json:"quantity"`
	Frequency         string    `json:"frequency"`
	LastFed           time.Time `json:"last_fed"`
	Notes             string    `json:"notes"`
	IsActive          bool      `json:"is_active"`
}

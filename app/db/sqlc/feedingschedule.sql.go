// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: feedingschedule.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFeedingSchedule = `-- name: CreateFeedingSchedule :one
INSERT INTO FeedingSchedule (petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive
`

type CreateFeedingScheduleParams struct {
	Petid     pgtype.Int8      `json:"petid"`
	Mealtime  pgtype.Time      `json:"mealtime"`
	Foodtype  string           `json:"foodtype"`
	Quantity  float64          `json:"quantity"`
	Frequency string           `json:"frequency"`
	Lastfed   pgtype.Timestamp `json:"lastfed"`
	Notes     pgtype.Text      `json:"notes"`
	Isactive  pgtype.Bool      `json:"isactive"`
}

func (q *Queries) CreateFeedingSchedule(ctx context.Context, arg CreateFeedingScheduleParams) (Feedingschedule, error) {
	row := q.db.QueryRow(ctx, createFeedingSchedule,
		arg.Petid,
		arg.Mealtime,
		arg.Foodtype,
		arg.Quantity,
		arg.Frequency,
		arg.Lastfed,
		arg.Notes,
		arg.Isactive,
	)
	var i Feedingschedule
	err := row.Scan(
		&i.Feedingscheduleid,
		&i.Petid,
		&i.Mealtime,
		&i.Foodtype,
		&i.Quantity,
		&i.Frequency,
		&i.Lastfed,
		&i.Notes,
		&i.Isactive,
	)
	return i, err
}

const deleteFeedingSchedule = `-- name: DeleteFeedingSchedule :exec
DELETE FROM FeedingSchedule
WHERE feedingScheduleID = $1
`

func (q *Queries) DeleteFeedingSchedule(ctx context.Context, feedingscheduleid int64) error {
	_, err := q.db.Exec(ctx, deleteFeedingSchedule, feedingscheduleid)
	return err
}

const getFeedingScheduleByPetID = `-- name: GetFeedingScheduleByPetID :many
SELECT feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive
FROM FeedingSchedule
WHERE petID = $1
ORDER BY feedingScheduleID
`

func (q *Queries) GetFeedingScheduleByPetID(ctx context.Context, petid pgtype.Int8) ([]Feedingschedule, error) {
	rows, err := q.db.Query(ctx, getFeedingScheduleByPetID, petid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Feedingschedule{}
	for rows.Next() {
		var i Feedingschedule
		if err := rows.Scan(
			&i.Feedingscheduleid,
			&i.Petid,
			&i.Mealtime,
			&i.Foodtype,
			&i.Quantity,
			&i.Frequency,
			&i.Lastfed,
			&i.Notes,
			&i.Isactive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listActiveFeedingSchedules = `-- name: ListActiveFeedingSchedules :many
SELECT feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive
FROM FeedingSchedule
WHERE isActive = true
ORDER BY petID
`

func (q *Queries) ListActiveFeedingSchedules(ctx context.Context) ([]Feedingschedule, error) {
	rows, err := q.db.Query(ctx, listActiveFeedingSchedules)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Feedingschedule{}
	for rows.Next() {
		var i Feedingschedule
		if err := rows.Scan(
			&i.Feedingscheduleid,
			&i.Petid,
			&i.Mealtime,
			&i.Foodtype,
			&i.Quantity,
			&i.Frequency,
			&i.Lastfed,
			&i.Notes,
			&i.Isactive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFeedingSchedule = `-- name: UpdateFeedingSchedule :exec
UPDATE FeedingSchedule
SET mealTime = $2, foodType = $3, quantity = $4, frequency = $5, lastFed = $6, notes = $7, isActive = $8
WHERE feedingScheduleID = $1
`

type UpdateFeedingScheduleParams struct {
	Feedingscheduleid int64            `json:"feedingscheduleid"`
	Mealtime          pgtype.Time      `json:"mealtime"`
	Foodtype          string           `json:"foodtype"`
	Quantity          float64          `json:"quantity"`
	Frequency         string           `json:"frequency"`
	Lastfed           pgtype.Timestamp `json:"lastfed"`
	Notes             pgtype.Text      `json:"notes"`
	Isactive          pgtype.Bool      `json:"isactive"`
}

func (q *Queries) UpdateFeedingSchedule(ctx context.Context, arg UpdateFeedingScheduleParams) error {
	_, err := q.db.Exec(ctx, updateFeedingSchedule,
		arg.Feedingscheduleid,
		arg.Mealtime,
		arg.Foodtype,
		arg.Quantity,
		arg.Frequency,
		arg.Lastfed,
		arg.Notes,
		arg.Isactive,
	)
	return err
}
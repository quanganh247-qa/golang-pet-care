// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: service.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createService = `-- name: CreateService :one
INSERT INTO Service (
  typeID,
  name,
  price,
  duration,
  description,
  isAvailable
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING serviceid, typeid, name, price, duration, description, isavailable, removed_at
`

type CreateServiceParams struct {
	Typeid      pgtype.Int8     `json:"typeid"`
	Name        string          `json:"name"`
	Price       pgtype.Float8   `json:"price"`
	Duration    pgtype.Interval `json:"duration"`
	Description pgtype.Text     `json:"description"`
	Isavailable pgtype.Bool     `json:"isavailable"`
}

func (q *Queries) CreateService(ctx context.Context, arg CreateServiceParams) (Service, error) {
	row := q.db.QueryRow(ctx, createService,
		arg.Typeid,
		arg.Name,
		arg.Price,
		arg.Duration,
		arg.Description,
		arg.Isavailable,
	)
	var i Service
	err := row.Scan(
		&i.Serviceid,
		&i.Typeid,
		&i.Name,
		&i.Price,
		&i.Duration,
		&i.Description,
		&i.Isavailable,
		&i.RemovedAt,
	)
	return i, err
}

const deleteService = `-- name: DeleteService :exec
DELETE FROM Service WHERE serviceID = $1
`

func (q *Queries) DeleteService(ctx context.Context, serviceid int64) error {
	_, err := q.db.Exec(ctx, deleteService, serviceid)
	return err
}

const getAllServices = `-- name: GetAllServices :many
SELECT serviceid, typeid, name, price, duration, description, isavailable, removed_at FROM Service ORDER BY name LIMIT $1 OFFSET $2
`

type GetAllServicesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllServices(ctx context.Context, arg GetAllServicesParams) ([]Service, error) {
	rows, err := q.db.Query(ctx, getAllServices, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.Serviceid,
			&i.Typeid,
			&i.Name,
			&i.Price,
			&i.Duration,
			&i.Description,
			&i.Isavailable,
			&i.RemovedAt,
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

const getServiceByID = `-- name: GetServiceByID :one
SELECT serviceid, typeid, name, price, duration, description, isavailable, removed_at FROM Service 
WHERE serviceID = $1 LIMIT 1
`

func (q *Queries) GetServiceByID(ctx context.Context, serviceid int64) (Service, error) {
	row := q.db.QueryRow(ctx, getServiceByID, serviceid)
	var i Service
	err := row.Scan(
		&i.Serviceid,
		&i.Typeid,
		&i.Name,
		&i.Price,
		&i.Duration,
		&i.Description,
		&i.Isavailable,
		&i.RemovedAt,
	)
	return i, err
}

const updateService = `-- name: UpdateService :exec
UPDATE Service Set
  typeID = $2,
  name = $3,
  price = $4,
  duration = $5,
  description = $6,
  isAvailable = $7
WHERE serviceID = $1
`

type UpdateServiceParams struct {
	Serviceid   int64           `json:"serviceid"`
	Typeid      pgtype.Int8     `json:"typeid"`
	Name        string          `json:"name"`
	Price       pgtype.Float8   `json:"price"`
	Duration    pgtype.Interval `json:"duration"`
	Description pgtype.Text     `json:"description"`
	Isavailable pgtype.Bool     `json:"isavailable"`
}

func (q *Queries) UpdateService(ctx context.Context, arg UpdateServiceParams) error {
	_, err := q.db.Exec(ctx, updateService,
		arg.Serviceid,
		arg.Typeid,
		arg.Name,
		arg.Price,
		arg.Duration,
		arg.Description,
		arg.Isavailable,
	)
	return err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: page.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPage = `-- name: CreatePage :one
INSERT INTO pages (name, content, project_id, slug, created_at, updated_at)
VALUES ($1, $2, $3, $4,now(), now())
RETURNING id, project_id, name, slug, created_at, updated_at, content, category_name, component_code, removed_at
`

type CreatePageParams struct {
	Name      string      `json:"name"`
	Content   pgtype.Text `json:"content"`
	ProjectID pgtype.Int4 `json:"project_id"`
	Slug      string      `json:"slug"`
}

func (q *Queries) CreatePage(ctx context.Context, arg CreatePageParams) (Page, error) {
	row := q.db.QueryRow(ctx, createPage,
		arg.Name,
		arg.Content,
		arg.ProjectID,
		arg.Slug,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Name,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Content,
		&i.CategoryName,
		&i.ComponentCode,
		&i.RemovedAt,
	)
	return i, err
}

const deletePage = `-- name: DeletePage :one
DELETE FROM pages WHERE id = $1 RETURNING id, project_id, name, slug, created_at, updated_at, content, category_name, component_code, removed_at
`

func (q *Queries) DeletePage(ctx context.Context, id int64) (Page, error) {
	row := q.db.QueryRow(ctx, deletePage, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Name,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Content,
		&i.CategoryName,
		&i.ComponentCode,
		&i.RemovedAt,
	)
	return i, err
}

const getPages = `-- name: GetPages :many
SELECT id, project_id, name, slug, created_at, updated_at, content, category_name, component_code, removed_at FROM pages ORDER BY created_at DESC
`

func (q *Queries) GetPages(ctx context.Context) ([]Page, error) {
	rows, err := q.db.Query(ctx, getPages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Page{}
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.Name,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Content,
			&i.CategoryName,
			&i.ComponentCode,
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

const updatePage = `-- name: UpdatePage :one
UPDATE pages
SET name = $1, content = $2, updated_at = $3
WHERE id = $4
RETURNING id, project_id, name, slug, created_at, updated_at, content, category_name, component_code, removed_at
`

type UpdatePageParams struct {
	Name      string             `json:"name"`
	Content   pgtype.Text        `json:"content"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	ID        int64              `json:"id"`
}

func (q *Queries) UpdatePage(ctx context.Context, arg UpdatePageParams) (Page, error) {
	row := q.db.QueryRow(ctx, updatePage,
		arg.Name,
		arg.Content,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Name,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Content,
		&i.CategoryName,
		&i.ComponentCode,
		&i.RemovedAt,
	)
	return i, err
}

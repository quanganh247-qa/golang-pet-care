// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: invoice.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countInvoices = `-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices
`

func (q *Queries) CountInvoices(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countInvoices)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createInvoice = `-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number,
    amount,
    date,
    due_date,
    status,
    description,
    customer_name,
    type,
    appointment_id,
    test_order_id,
    order_id    
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING id, invoice_number, amount, date, due_date, status, description, customer_name, created_at, type, appointment_id, test_order_id, order_id
`

type CreateInvoiceParams struct {
	InvoiceNumber string      `json:"invoice_number"`
	Amount        float64     `json:"amount"`
	Date          pgtype.Date `json:"date"`
	DueDate       pgtype.Date `json:"due_date"`
	Status        string      `json:"status"`
	Description   pgtype.Text `json:"description"`
	CustomerName  pgtype.Text `json:"customer_name"`
	Type          pgtype.Text `json:"type"`
	AppointmentID pgtype.Int8 `json:"appointment_id"`
	TestOrderID   pgtype.Int8 `json:"test_order_id"`
	OrderID       pgtype.Int8 `json:"order_id"`
}

func (q *Queries) CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (Invoice, error) {
	row := q.db.QueryRow(ctx, createInvoice,
		arg.InvoiceNumber,
		arg.Amount,
		arg.Date,
		arg.DueDate,
		arg.Status,
		arg.Description,
		arg.CustomerName,
		arg.Type,
		arg.AppointmentID,
		arg.TestOrderID,
		arg.OrderID,
	)
	var i Invoice
	err := row.Scan(
		&i.ID,
		&i.InvoiceNumber,
		&i.Amount,
		&i.Date,
		&i.DueDate,
		&i.Status,
		&i.Description,
		&i.CustomerName,
		&i.CreatedAt,
		&i.Type,
		&i.AppointmentID,
		&i.TestOrderID,
		&i.OrderID,
	)
	return i, err
}

const createInvoiceItem = `-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (
    invoice_id,
    name,
    price,
    quantity
) VALUES (
    $1, $2, $3, $4
) RETURNING id, invoice_id, name, price, quantity
`

type CreateInvoiceItemParams struct {
	InvoiceID int32   `json:"invoice_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int32   `json:"quantity"`
}

func (q *Queries) CreateInvoiceItem(ctx context.Context, arg CreateInvoiceItemParams) (InvoiceItem, error) {
	row := q.db.QueryRow(ctx, createInvoiceItem,
		arg.InvoiceID,
		arg.Name,
		arg.Price,
		arg.Quantity,
	)
	var i InvoiceItem
	err := row.Scan(
		&i.ID,
		&i.InvoiceID,
		&i.Name,
		&i.Price,
		&i.Quantity,
	)
	return i, err
}

const deleteInvoice = `-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE id = $1
`

func (q *Queries) DeleteInvoice(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteInvoice, id)
	return err
}

const deleteInvoiceItem = `-- name: DeleteInvoiceItem :exec
DELETE FROM invoice_items
WHERE id = $1
`

func (q *Queries) DeleteInvoiceItem(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteInvoiceItem, id)
	return err
}

const getInvoiceByID = `-- name: GetInvoiceByID :one
SELECT id, invoice_number, amount, date, due_date, status, description, customer_name, created_at, type, appointment_id, test_order_id, order_id FROM invoices
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetInvoiceByID(ctx context.Context, id int32) (Invoice, error) {
	row := q.db.QueryRow(ctx, getInvoiceByID, id)
	var i Invoice
	err := row.Scan(
		&i.ID,
		&i.InvoiceNumber,
		&i.Amount,
		&i.Date,
		&i.DueDate,
		&i.Status,
		&i.Description,
		&i.CustomerName,
		&i.CreatedAt,
		&i.Type,
		&i.AppointmentID,
		&i.TestOrderID,
		&i.OrderID,
	)
	return i, err
}

const getInvoiceByNumber = `-- name: GetInvoiceByNumber :one
SELECT id, invoice_number, amount, date, due_date, status, description, customer_name, created_at, type, appointment_id, test_order_id, order_id FROM invoices
WHERE invoice_number = $1 LIMIT 1
`

func (q *Queries) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (Invoice, error) {
	row := q.db.QueryRow(ctx, getInvoiceByNumber, invoiceNumber)
	var i Invoice
	err := row.Scan(
		&i.ID,
		&i.InvoiceNumber,
		&i.Amount,
		&i.Date,
		&i.DueDate,
		&i.Status,
		&i.Description,
		&i.CustomerName,
		&i.CreatedAt,
		&i.Type,
		&i.AppointmentID,
		&i.TestOrderID,
		&i.OrderID,
	)
	return i, err
}

const getInvoiceItemByID = `-- name: GetInvoiceItemByID :one
SELECT id, invoice_id, name, price, quantity FROM invoice_items
WHERE id = $1
`

func (q *Queries) GetInvoiceItemByID(ctx context.Context, id int32) (InvoiceItem, error) {
	row := q.db.QueryRow(ctx, getInvoiceItemByID, id)
	var i InvoiceItem
	err := row.Scan(
		&i.ID,
		&i.InvoiceID,
		&i.Name,
		&i.Price,
		&i.Quantity,
	)
	return i, err
}

const getInvoiceItems = `-- name: GetInvoiceItems :many
SELECT id, invoice_id, name, price, quantity FROM invoice_items
WHERE invoice_id = $1
`

func (q *Queries) GetInvoiceItems(ctx context.Context, invoiceID int32) ([]InvoiceItem, error) {
	rows, err := q.db.Query(ctx, getInvoiceItems, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []InvoiceItem{}
	for rows.Next() {
		var i InvoiceItem
		if err := rows.Scan(
			&i.ID,
			&i.InvoiceID,
			&i.Name,
			&i.Price,
			&i.Quantity,
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

const getInvoiceWithItems = `-- name: GetInvoiceWithItems :many
SELECT 
    i.id,
    i.invoice_number,
    i.amount,
    i.date,
    i.due_date,
    i.status,
    i.description,
    i.customer_name,
    i.created_at,
    ii.id as item_id,
    ii.name as item_name,
    ii.price as item_price,
    ii.quantity as item_quantity
FROM invoices i
LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
WHERE i.id = $1
`

type GetInvoiceWithItemsRow struct {
	ID            int32            `json:"id"`
	InvoiceNumber string           `json:"invoice_number"`
	Amount        float64          `json:"amount"`
	Date          pgtype.Date      `json:"date"`
	DueDate       pgtype.Date      `json:"due_date"`
	Status        string           `json:"status"`
	Description   pgtype.Text      `json:"description"`
	CustomerName  pgtype.Text      `json:"customer_name"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	ItemID        pgtype.Int4      `json:"item_id"`
	ItemName      pgtype.Text      `json:"item_name"`
	ItemPrice     pgtype.Float8    `json:"item_price"`
	ItemQuantity  pgtype.Int4      `json:"item_quantity"`
}

func (q *Queries) GetInvoiceWithItems(ctx context.Context, id int32) ([]GetInvoiceWithItemsRow, error) {
	rows, err := q.db.Query(ctx, getInvoiceWithItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetInvoiceWithItemsRow{}
	for rows.Next() {
		var i GetInvoiceWithItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.InvoiceNumber,
			&i.Amount,
			&i.Date,
			&i.DueDate,
			&i.Status,
			&i.Description,
			&i.CustomerName,
			&i.CreatedAt,
			&i.ItemID,
			&i.ItemName,
			&i.ItemPrice,
			&i.ItemQuantity,
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

const listInvoices = `-- name: ListInvoices :many
SELECT id, invoice_number, amount, date, due_date, status, description, customer_name, created_at, type, appointment_id, test_order_id, order_id FROM invoices
ORDER BY date DESC
LIMIT $1
OFFSET $2
`

type ListInvoicesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListInvoices(ctx context.Context, arg ListInvoicesParams) ([]Invoice, error) {
	rows, err := q.db.Query(ctx, listInvoices, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.ID,
			&i.InvoiceNumber,
			&i.Amount,
			&i.Date,
			&i.DueDate,
			&i.Status,
			&i.Description,
			&i.CustomerName,
			&i.CreatedAt,
			&i.Type,
			&i.AppointmentID,
			&i.TestOrderID,
			&i.OrderID,
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

const updateInvoiceAmount = `-- name: UpdateInvoiceAmount :exec
UPDATE invoices
SET 
    amount = (
        SELECT COALESCE(SUM(price * quantity), 0) 
        FROM invoice_items 
        WHERE invoice_id = $1
    )
WHERE id = $1
`

func (q *Queries) UpdateInvoiceAmount(ctx context.Context, invoiceID int32) error {
	_, err := q.db.Exec(ctx, updateInvoiceAmount, invoiceID)
	return err
}

const updateInvoiceItem = `-- name: UpdateInvoiceItem :one
UPDATE invoice_items
SET 
    name = COALESCE($2, name),
    price = COALESCE($3, price),
    quantity = COALESCE($4, quantity)
WHERE id = $1
RETURNING id, invoice_id, name, price, quantity
`

type UpdateInvoiceItemParams struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

func (q *Queries) UpdateInvoiceItem(ctx context.Context, arg UpdateInvoiceItemParams) (InvoiceItem, error) {
	row := q.db.QueryRow(ctx, updateInvoiceItem,
		arg.ID,
		arg.Name,
		arg.Price,
		arg.Quantity,
	)
	var i InvoiceItem
	err := row.Scan(
		&i.ID,
		&i.InvoiceID,
		&i.Name,
		&i.Price,
		&i.Quantity,
	)
	return i, err
}

const updateInvoiceStatus = `-- name: UpdateInvoiceStatus :exec
UPDATE invoices
SET status = $2
WHERE id = $1
`

type UpdateInvoiceStatusParams struct {
	ID     int32  `json:"id"`
	Status string `json:"status"`
}

func (q *Queries) UpdateInvoiceStatus(ctx context.Context, arg UpdateInvoiceStatusParams) error {
	_, err := q.db.Exec(ctx, updateInvoiceStatus, arg.ID, arg.Status)
	return err
}

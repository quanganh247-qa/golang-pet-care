// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateServiceType(ctx context.Context, servicetypename pgtype.Text) (Servicetype, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)

package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Connection struct {
	Close func()
}

func Init(config util.Config) (*Connection, error) {
	_, err := token.NewJWTMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker %w", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	db.InitStore(connPool)

	conn := &Connection{
		Close: func() {
			connPool.Close()
			fmt.Println("Database connection pool closed.")
		},
	}
	return conn, nil
}

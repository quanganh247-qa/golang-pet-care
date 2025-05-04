package device_token

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DeviceTokenServiceInterface interface {
	InsertToken(ctx context.Context, req DVTRequest, username string) (*DVTResponse, error)
	DeleteDevicetToken(ctx context.Context, username string, token string) error
	GetDeviceTokenByUsername(ctx context.Context, username string) ([]DVTResponse, error)
}

func (s *DeviceTokenService) InsertToken(ctx context.Context, req DVTRequest, username string) (*DVTResponse, error) {

	lastUseAt, expiredAt, err := util.ParseStringToTime(req.LastUsedAt, req.ExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}
	token, err := s.storeDB.InsertDeviceToken(ctx, db.InsertDeviceTokenParams{
		Username:   username,
		Token:      req.Token,
		DeviceType: pgtype.Text{String: req.DeviceType, Valid: true},
		LastUsedAt: pgtype.Timestamp{Time: lastUseAt, Valid: true},
		ExpiredAt:  pgtype.Timestamp{Time: expiredAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert token: %w", err)
	}

	return &DVTResponse{
		ID:         token.ID,
		Username:   token.Username,
		Token:      token.Token,
		DeviceType: token.DeviceType.String,
		CreatedAt:  token.CreatedAt.Time.Format(time.RFC3339),
		LastUsedAt: token.LastUsedAt.Time.Format(time.RFC3339),
		ExpiredAt:  token.ExpiredAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *DeviceTokenService) DeleteDevicetToken(ctx context.Context, username string, token string) error {

	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err = q.DeleteDeviceToken(ctx, db.DeleteDeviceTokenParams{
			Username: username,
			Token:    token,
		})
		if err != nil {
			return fmt.Errorf("failed to delete token: %w", err)

		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return err

}

func (s *DeviceTokenService) GetDeviceTokenByUsername(ctx context.Context, username string) ([]DVTResponse, error) {
	tokens, err := s.storeDB.GetDeviceTokenByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get device token: %w", err)
	}

	var response []DVTResponse
	for _, token := range tokens {
		response = append(response, DVTResponse{
			ID:         token.ID,
			Username:   token.Username,
			Token:      token.Token,
			DeviceType: token.DeviceType.String,
			CreatedAt:  token.CreatedAt.Time.Format(time.RFC3339),
			LastUsedAt: token.LastUsedAt.Time.Format(time.RFC3339),
			ExpiredAt:  token.ExpiredAt.Time.Format(time.RFC3339),
		})
	}

	return response, nil
}

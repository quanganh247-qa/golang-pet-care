package calendar

// type StateStore interface {
// 	SaveState(ctx context.Context, state *OAuthState) error
// 	GetState(ctx context.Context, stateString string) (*OAuthState, error)
// 	DeleteState(ctx context.Context, stateString string) error
// }

// func (s *GoogleCalendarService) saveState(ctx context.Context, state *OAuthState) error {
// 	err := s.dbStore.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		_, err := q.SaveState(ctx, db.SaveStateParams{
// 			State:     state.State,
// 			Username:  state.Username,
// 			CreatedAt: pgtype.Timestamp{Time: state.CreatedAt, Valid: true},
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to save state: %w", err)
// 		}
// 		return nil
// 	})

// 	return fmt.Errorf("transaction failed: %w", err)
// }

// func (s *GoogleCalendarService) GetState(ctx context.Context, stateString string) (*OAuthState, error) {
// 	var state OAuthState
// 	err := s.dbStore.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		res, err := q.GetState(ctx, stateString)
// 		if err != nil {
// 			return fmt.Errorf("failed to get state: %w", err)
// 		}
// 		state = OAuthState{
// 			State:     res.State,
// 			Username:  res.Username,
// 			CreatedAt: res.CreatedAt.Time,
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("transaction failed: %w", err)
// 	}

// 	return &state, nil
// }

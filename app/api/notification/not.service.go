package notification

type NotServiceInterface interface {
}

// Insert notification
// func (s *NotService) InsertNotification(ctx context.Context, arg NotificationRequest) (*NotificationResponse, error) {

// 	dueDate := time.Time.Format(arg.DueDate, time.RFC3339)

// 	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		_, err := q.InsertNotification(ctx, db.InsertNotificationParams{
// 			Petid:            pgtype.Int8{Int64: arg.PetID, Valid: true},
// 			Title:            arg.Title,
// 			Body:             pgtype.Text{String: arg.Body, Valid: true},
// 			Duedate:          pgtype.Timestamp{Time: dueDate, Status: pgtype.Present},
// 			RepeatInterval:   arg.RepeatInterval,
// 			NotificationSent: arg.NotificationSent,
// 		})
// 		return err
// 	})
// }

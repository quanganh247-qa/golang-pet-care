package reminder

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ReminderServiceInterface interface {
	CreateReminderService(ctx *gin.Context, arg CreateReminderRequest) (*db.Reminder, error)
}

func (s *RemiderService) CreateReminderService(ctx *gin.Context, arg CreateReminderRequest) (*db.Reminder, error) {

	var res db.Reminder
	fmt.Println(arg.Petid.Int, arg.Userid.Int, arg.Remindertype, arg.Reminderdate, arg.Description, arg.Status)
	parsedTime, err := time.Parse("2006-01-02", arg.Reminderdate)
	if err != nil {
		return nil, fmt.Errorf("cannot parse reminder date: %v", err)
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err = q.CreateReminder(ctx, db.CreateReminderParams{
			Userid:       pgtype.Int8{Int64: arg.Userid.Int, Valid: true},
			Petid:        pgtype.Int8{Int64: arg.Petid.Int, Valid: true},
			Remindertype: arg.Remindertype,
			Reminderdate: pgtype.Date{Time: parsedTime, Valid: true},
			Description:  pgtype.Text{String: arg.Description.String, Valid: true},
			Status:       pgtype.Text{String: arg.Status.String, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("cannot create reminder: %v", err)
		}
		return nil
	})

	return &db.Reminder{
		Reminderid:   res.Reminderid,
		Userid:       res.Userid,
		Petid:        res.Petid,
		Remindertype: res.Remindertype,
		Reminderdate: res.Reminderdate,
		Description:  res.Description,
		Status:       res.Status,
	}, nil
}

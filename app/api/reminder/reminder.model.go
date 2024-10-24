package reminder

import (
	"github.com/jackc/pgx/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type RemiderService struct {
	storeDB db.Store
}

type ReminderController struct {
	reminderService ReminderServiceInterface
}
type ReminderApi struct {
	controller ReminderControllerInterface
}

type CreateReminderRequest struct {
	Userid       pgtype.Int8 `json:"userid"`
	Petid        pgtype.Int8 `json:"petid"`
	Remindertype string      `json:"remindertype"`
	Reminderdate string      `json:"reminderdate"`
	Description  pgtype.Text `json:"description"`
	Status       pgtype.Text `json:"status"`
}

type ReminderResponse struct {
	Reminderid int64  `json:"reminderid"`
	Pet        db.Pet `json:"pet"`
}

// api -> controller -> service -> model -> db

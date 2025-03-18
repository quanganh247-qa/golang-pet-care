package queue

import (
	"fmt"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type QueueServiceInterface interface {
	GetQueueService(ctx *gin.Context) ([]QueueItem, error)
	UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error
}

func (s *QueueService) GetQueueService(ctx *gin.Context) ([]QueueItem, error) {
	appointments, err := s.storeDB.GetAppointmentsQueue(ctx) // Assuming 1 is the state ID for waiting/arrived
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}

	var queueItems []QueueItem
	for _, appointment := range appointments {

		a, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointment detail: %v", err)
		}

		doc, err := s.storeDB.GetDoctor(ctx, a.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %v", err)
		}

		// Calculate waiting time
		waitingSince := appointment.Date.Time
		actualWaitTime := time.Since(waitingSince).Round(time.Minute).String()

		queueItem := QueueItem{
			ID:              appointment.AppointmentID,
			PatientName:     a.PetName.String,
			Status:          a.StateName.String, // You might want to get this from a separate status field
			AppointmentType: a.ServiceName.String,
			Doctor:          doc.Name,
			WaitingSince:    waitingSince.Format("3:04 PM"),
			ActualWaitTime:  actualWaitTime,
		}

		queueItems = append(queueItems, queueItem)
	}

	// Sort queue items by priority (high first) and waiting time
	sortQueueItems(queueItems)

	return queueItems, nil
}

func (s *QueueService) UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error {
	// Update appointment status
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// You might need to add a new query to update appointment status
		// For now, we'll just update the state_id
		stateID := int32(1) // You should map status to appropriate state_id
		if status == "in_progress" {
			stateID = 2
		} else if status == "completed" {
			stateID = 3
		}

		err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: id,
			StateID:       pgtype.Int4{Int32: stateID, Valid: true},
		})
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to update appointment status: %v", err)
	}

	return nil
}

// // Helper function to sort queue items
// func sortQueueItems(items []QueueItem) {
// 	// Sort by priority (high first) and then by waiting time
// 	sort.Slice(items, func(i, j int) bool {
// 		// If priorities are the same, sort by waiting time
// 		return items[i].WaitingSince < items[j].WaitingSince
// 	})
// }

// Helper function to sort queue items
func sortQueueItems(items []QueueItem) {
	// Sort by priority (high first) and then by waiting time
	sort.Slice(items, func(i, j int) bool {
		// If priorities are different, sort by priority (high first)
		if items[i].Priority != items[j].Priority {
			return items[i].Priority > items[j].Priority
		}
		// If priorities are the same, sort by waiting time (longer wait first)
		return items[i].WaitingSince < items[j].WaitingSince
	})
}

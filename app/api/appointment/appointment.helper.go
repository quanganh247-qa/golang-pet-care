package appointment

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

// Redis key constants for appointments
const (
	APPOINTMENT_INFO_KEY      = "APPOINTMENT_INFO"
	APPOINTMENT_LIST_USER_KEY = "APPOINTMENT_LIST_USER"
	APPOINTMENT_LIST_DOC_KEY  = "APPOINTMENT_LIST_DOC"
	APPOINTMENT_SOAP_KEY      = "APPOINTMENT_SOAP"
	APPOINTMENT_QUEUE_KEY     = "APPOINTMENT_QUEUE"
	TIME_SLOTS_KEY            = "TIME_SLOTS"
	LOCK_PREFIX               = "LOCK"
)

// Helper function to sort queue items
func sortQueueItems(items []QueueItem) {
	// Define priority order
	priorityOrder := map[string]int{
		"Normal": 1,
		"High":   2,
	}
	// Sort by priority (high first) and then by waiting time
	sort.Slice(items, func(i, j int) bool {
		// If priorities are different, sort by priority (high first)
		if items[i].Priority != items[j].Priority {
			return priorityOrder[items[i].Priority] > priorityOrder[items[j].Priority]
		}
		// If priorities are the same, sort by waiting time (longer wait first)
		return items[i].WaitingSince < items[j].WaitingSince
	})
}

// isObjectiveEmpty checks if the ObjectiveData struct is empty
func isObjectiveEmpty(obj ObjectiveData) bool {
	// Add appropriate checks based on your ObjectiveData struct fields
	// This is an example - adjust according to your actual struct fields
	return obj == ObjectiveData{} // Compare with empty struct
}

// GetAppointmentCacheKey generates a Redis key for a specific appointment
func GetAppointmentCacheKey(id int64) string {
	return fmt.Sprintf("%s:%d", APPOINTMENT_INFO_KEY, id)
}

// GetAppointmentListByUserCacheKey generates a Redis key for a user's appointments
func GetAppointmentListByUserCacheKey(username string) string {
	return fmt.Sprintf("%s:%s", APPOINTMENT_LIST_USER_KEY, username)
}

// GetAppointmentListByDoctorCacheKey generates a Redis key for a doctor's appointments
func GetAppointmentListByDoctorCacheKey(doctorID int64) string {
	return fmt.Sprintf("%s:%d", APPOINTMENT_LIST_DOC_KEY, doctorID)
}

// GetAppointmentSOAPCacheKey generates a Redis key for an appointment's SOAP
func GetAppointmentSOAPCacheKey(appointmentID int64) string {
	return fmt.Sprintf("%s:%d", APPOINTMENT_SOAP_KEY, appointmentID)
}

// GetQueueCacheKey generates a Redis key for the appointment queue
func GetQueueCacheKey(username string) string {
	return fmt.Sprintf("%s:%s", APPOINTMENT_QUEUE_KEY, username)
}

// GetTimeSlotsKey generates a Redis key for time slots
func GetTimeSlotsKey(doctorID int64, date string) string {
	return fmt.Sprintf("%s:%d:%s", TIME_SLOTS_KEY, doctorID, date)
}

// GetLockKey generates a distributed lock key for concurrency control
func GetLockKey(resourceID string) string {
	return fmt.Sprintf("%s:%s", LOCK_PREFIX, resourceID)
}

// AcquireLock attempts to acquire a distributed lock
// Returns true if lock acquired successfully, false otherwise
func AcquireLock(resourceID string, ttl time.Duration) bool {
	if redis.Client == nil {
		return true // If Redis is not available, proceed without locking
	}

	lockKey := GetLockKey(resourceID)
	// Use Redis SET NX (Only set if key doesn't exist)
	success, err := redis.Client.RedisClient.SetNX(context.Background(), lockKey, "locked", ttl).Result()
	if err != nil {
		return false
	}
	return success
}

// ReleaseLock releases a previously acquired lock
func ReleaseLock(resourceID string) {
	if redis.Client == nil {
		return
	}

	lockKey := GetLockKey(resourceID)
	redis.Client.RemoveCacheByKey(lockKey)
}

// ClearAppointmentCache clears all appointment-related caches
func ClearAppointmentCache() {
	if redis.Client == nil {
		return
	}

	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", APPOINTMENT_INFO_KEY))
	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", APPOINTMENT_LIST_USER_KEY))
	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", APPOINTMENT_LIST_DOC_KEY))
	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", APPOINTMENT_SOAP_KEY))
	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", APPOINTMENT_QUEUE_KEY))
	redis.Client.RemoveCacheBySubString(fmt.Sprintf("%s*", TIME_SLOTS_KEY))
}

// ClearAppointmentCacheByID clears cache for a specific appointment
func ClearAppointmentCacheByID(appointmentID int64) {
	if redis.Client == nil {
		return
	}

	redis.Client.RemoveCacheByKey(GetAppointmentCacheKey(appointmentID))
	redis.Client.RemoveCacheByKey(GetAppointmentSOAPCacheKey(appointmentID))
}

// SetUserRole thiết lập vai trò cho người dùng
func (nm *NotificationManager) SetUserRole(username string, role string) {
	nm.roleMutex.Lock()
	defer nm.roleMutex.Unlock()
	nm.userRoles[username] = role
}

// GetUserRole lấy vai trò của người dùng
func (nm *NotificationManager) GetUserRole(username string) (string, bool) {
	nm.roleMutex.RLock()
	defer nm.roleMutex.RUnlock()
	role, exists := nm.userRoles[username]
	return role, exists
}

// AddMissedNotification lưu thông báo vào danh sách thông báo bị bỏ lỡ
func (nm *NotificationManager) AddMissedNotification(username string, notification AppointmentNotification) {
	nm.missedMutex.Lock()
	defer nm.missedMutex.Unlock()

	if _, exists := nm.missedNotifications[username]; !exists {
		nm.missedNotifications[username] = make([]AppointmentNotification, 0)
	}
	nm.missedNotifications[username] = append(nm.missedNotifications[username], notification)
}

// GetMissedNotifications lấy danh sách thông báo bị bỏ lỡ của một người dùng
func (nm *NotificationManager) GetMissedNotifications(username string) []AppointmentNotification {
	nm.missedMutex.RLock()
	defer nm.missedMutex.RUnlock()

	if notifications, exists := nm.missedNotifications[username]; exists {
		return notifications
	}
	return []AppointmentNotification{}
}

// ClearMissedNotifications xóa tất cả thông báo bị bỏ lỡ của một người dùng
func (nm *NotificationManager) ClearMissedNotifications(username string) {
	nm.missedMutex.Lock()
	defer nm.missedMutex.Unlock()

	if _, exists := nm.missedNotifications[username]; exists {
		nm.missedNotifications[username] = []AppointmentNotification{}
	}
}

// MoveMissedToNotifications chuyển tất cả thông báo bị bỏ lỡ của một người dùng vào danh sách thông báo thông thường
func (nm *NotificationManager) MoveMissedToNotifications(username string) {
	// Lấy danh sách thông báo bị bỏ lỡ
	missedNotifications := nm.GetMissedNotifications(username)
	if len(missedNotifications) == 0 {
		return
	}

	// Thêm vào danh sách thông báo thông thường
	for _, notification := range missedNotifications {
		nm.AddNotification(username, notification)
	}

	// Xóa danh sách thông báo bị bỏ lỡ
	nm.ClearMissedNotifications(username)
}

// Thêm các phương thức cho NotificationManager
func (nm *NotificationManager) AddNotification(username string, notification AppointmentNotification) {

	ctx := context.Background()
	err := nm.SaveNotificationToDB(ctx, username, notification)
	if err != nil {
		// Log lỗi nhưng vẫn tiếp tục để thông báo được lưu trong bộ nhớ
		fmt.Printf("Lỗi khi lưu thông báo vào cơ sở dữ liệu: %v\n", err)
	}

	nm.notificationMutex.Lock()
	defer nm.notificationMutex.Unlock()

	// Thêm thông báo vào danh sách thông báo của người dùng
	if _, exists := nm.notifications[username]; !exists {
		nm.notifications[username] = make([]AppointmentNotification, 0)
	}
	nm.notifications[username] = append(nm.notifications[username], notification)

	// Gửi thông báo đến client nếu đang chờ
	nm.clientMutex.RLock()
	defer nm.clientMutex.RUnlock()

	if ch, exists := nm.pendingClients[username]; exists {
		// Chỉ gửi nếu kênh vẫn mở
		select {
		case ch <- notification:
			// Thông báo đã được gửi
		default:
			// Kênh đã đóng hoặc đầy
		}
	}
}

// BroadcastNotificationToAll gửi thông báo đến tất cả người dùng (giữ lại để tương thích ngược)
func (nm *NotificationManager) BroadcastNotificationToAll(notification AppointmentNotification) {
	nm.notificationMutex.RLock()
	// Tạo danh sách tất cả người dùng
	usernames := make([]string, 0, len(nm.notifications))
	for username := range nm.notifications {
		usernames = append(usernames, username)
	}
	nm.notificationMutex.RUnlock()

	// Thêm thông báo cho tất cả người dùng
	for _, username := range usernames {
		nm.AddNotification(username, notification)
	}
}

// BroadcastNotificationByRole gửi thông báo đến tất cả người dùng có vai trò cụ thể
func (nm *NotificationManager) BroadcastNotificationByRole(notification AppointmentNotification, targetRole string) {
	// Danh sách người dùng có vai trò phù hợp
	usersWithRole := nm.getUsersByRole(targetRole)

	// Gửi thông báo đến từng người dùng có vai trò phù hợp
	for _, username := range usersWithRole {
		nm.AddNotification(username, notification)
	}
}

// getUsersByRole lấy danh sách người dùng có vai trò cụ thể
func (nm *NotificationManager) getUsersByRole(role string) []string {
	nm.roleMutex.RLock()
	defer nm.roleMutex.RUnlock()

	var usernames []string

	users, err := db.StoreDB.GetUserByRole(context.Background(), pgtype.Text{String: role, Valid: true})
	if err != nil {
		return []string{}
	}

	for _, user := range users {
		if user.Role.String == role {
			usernames = append(usernames, user.Username)
		}
	}
	return usernames
}

// GetNotifications trả về tất cả thông báo của một người dùng
func (nm *NotificationManager) GetNotifications(username string) []AppointmentNotification {
	nm.notificationMutex.RLock()
	defer nm.notificationMutex.RUnlock()

	if notifications, exists := nm.notifications[username]; exists {
		return notifications
	}
	return []AppointmentNotification{}
}

// ClearNotifications xóa tất cả thông báo của một người dùng
func (nm *NotificationManager) ClearNotifications(username string) {
	nm.notificationMutex.Lock()
	defer nm.notificationMutex.Unlock()

	if _, exists := nm.notifications[username]; exists {
		nm.notifications[username] = []AppointmentNotification{}
	}
}

// WaitForNotification thiết lập long polling để chờ thông báo mới
func (nm *NotificationManager) WaitForNotification(username string, timeout time.Duration) (AppointmentNotification, bool) {
	// Tạo kênh cho client này
	notificationCh := make(chan AppointmentNotification, 1)

	// Đăng ký client vào danh sách chờ
	nm.clientMutex.Lock()
	nm.pendingClients[username] = notificationCh
	nm.clientMutex.Unlock()

	// Đảm bảo kênh bị xóa khi hàm kết thúc
	defer func() {
		nm.clientMutex.Lock()
		delete(nm.pendingClients, username)
		nm.clientMutex.Unlock()
		close(notificationCh)
	}()

	// Kiểm tra xem có thông báo mới ngay lập tức không
	nm.notificationMutex.RLock()
	if notifications, exists := nm.notifications[username]; exists && len(notifications) > 0 {
		// Có thông báo mới, trả về thông báo đầu tiên và xóa nó
		notification := notifications[0]
		nm.notificationMutex.RUnlock()

		nm.notificationMutex.Lock()
		nm.notifications[username] = notifications[1:]
		nm.notificationMutex.Unlock()

		return notification, true
	}
	nm.notificationMutex.RUnlock()

	// Không có thông báo ngay lập tức, thiết lập timeout để chờ
	select {
	case notification := <-notificationCh:
		return notification, true
	case <-time.After(timeout):
		return AppointmentNotification{}, false
	}
}

// ConvertToDBNotification chuyển đổi một AppointmentNotification thành dạng phù hợp cho DB
func (nm *NotificationManager) ConvertToDBNotification(username string, notification AppointmentNotification) db.CreatetNotificationParams {
	var content string
	if notification.Reason != "" {
		content = fmt.Sprintf("Appointment for %s with doctor %s on %s. Reason: %s",
			notification.Pet.PetName, notification.Doctor.DoctorName, notification.Date, notification.Reason)
	} else {
		content = fmt.Sprintf("Appointment for %s with doctor %s on %s",
			notification.Pet.PetName, notification.Doctor.DoctorName, notification.Date)
	}

	return db.CreatetNotificationParams{
		Username:    username,
		Title:       notification.Title,
		Content:     pgtype.Text{String: content, Valid: true},
		NotifyType:  pgtype.Text{String: "appointment", Valid: true},
		RelatedID:   pgtype.Int4{Int32: int32(notification.AppointmentID), Valid: true},
		RelatedType: pgtype.Text{String: "appointment", Valid: true},
	}
}

// SaveNotificationToDB lưu thông báo vào cơ sở dữ liệu
func (nm *NotificationManager) SaveNotificationToDB(ctx context.Context, username string, notification AppointmentNotification) error {
	params := nm.ConvertToDBNotification(username, notification)
	_, err := db.StoreDB.CreatetNotification(ctx, params)
	return err
}

// GetNotificationsFromDB lấy danh sách thông báo từ cơ sở dữ liệu
func (nm *NotificationManager) GetNotificationsFromDB(ctx context.Context, username string, limit int32, offset int32) ([]DatabaseNotification, error) {
	dbNotifications, err := db.StoreDB.ListNotification(ctx, db.ListNotificationParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var notifications []DatabaseNotification
	for _, n := range dbNotifications {

		app, err := db.StoreDB.GetAppointmentByID(ctx, int64(n.RelatedID.Int32))
		if err != nil {
			return nil, err
		}
		if app.StateID.Int32 == 1 {
			notifications = append(notifications, DatabaseNotification{
				ID:          n.ID,
				Username:    n.Username,
				Title:       n.Title,
				Content:     n.Content.String,
				NotifyType:  n.NotifyType.String,
				RelatedID:   int64(n.RelatedID.Int32),
				RelatedType: n.RelatedType.String,
				CreatedAt:   n.Datetime.Time,
				IsRead:      n.IsRead.Bool,
			})
		}

	}

	return notifications, nil
}

// ConvertDBToAppointmentNotification chuyển đổi thông báo từ DB sang AppointmentNotification
func (nm *NotificationManager) ConvertDBToAppointmentNotification(dbNotification DatabaseNotification) AppointmentNotification {
	// ID mặc định
	id := fmt.Sprintf("db-%d", dbNotification.ID)

	// Mặc định cho các trường khác
	notification := AppointmentNotification{
		ID:            id,
		Title:         dbNotification.Title,
		AppointmentID: dbNotification.RelatedID,
		// Các trường khác sẽ được fill khi cần thiết
	}

	return notification
}

// MarkNotificationAsRead đánh dấu thông báo đã đọc trong cơ sở dữ liệu
func (nm *NotificationManager) MarkNotificationAsReadInDB(ctx context.Context, notificationID int64) error {
	return db.StoreDB.MarkNotificationAsRead(ctx, notificationID)
}

// DeleteNotificationsByUsername xóa tất cả thông báo của người dùng trong DB
func (nm *NotificationManager) DeleteNotificationsByUsernameFromDB(ctx context.Context, username string) error {
	return db.StoreDB.DeleteNotificationsByUsername(ctx, username)
}

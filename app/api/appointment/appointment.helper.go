package appointment

import (
	"context"
	"fmt"
	"sort"
	"time"

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

package appointment

import (
	"testing"
	"time"
)

func TestNotificationManager_SetGetUserRole(t *testing.T) {
	nm := NewNotificationManager()

	// Test set và get role
	nm.SetUserRole("user1", "doctor")
	nm.SetUserRole("user2", "customer")
	nm.SetUserRole("user3", "admin")

	// Kiểm tra lấy role
	role1, exists1 := nm.GetUserRole("user1")
	if !exists1 || role1 != "doctor" {
		t.Errorf("Expected role 'doctor' for user1, got %s, exists: %v", role1, exists1)
	}

	role2, exists2 := nm.GetUserRole("user2")
	if !exists2 || role2 != "customer" {
		t.Errorf("Expected role 'customer' for user2, got %s, exists: %v", role2, exists2)
	}

	// Kiểm tra user không tồn tại
	_, exists3 := nm.GetUserRole("user4")
	if exists3 {
		t.Error("Expected user4 to not exist, but it does")
	}
}

func TestNotificationManager_AddNotification(t *testing.T) {
	nm := NewNotificationManager()

	// Tạo thông báo test
	notification := AppointmentNotification{
		ID:            "test-1",
		Title:         "Test Notification",
		AppointmentID: 123,
		Pet: Pet{
			PetID:   456,
			PetName: "Buddy",
		},
		Doctor: Doctor{
			DoctorID:   789,
			DoctorName: "Dr. Smith",
		},
		Reason:      "Regular checkup",
		Date:        "2025-05-03",
		ServiceName: "General Checkup",
	}

	// Thêm thông báo cho user1
	nm.AddNotification("user1", notification)

	// Kiểm tra thông báo đã được thêm
	notifications := nm.GetNotifications("user1")
	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification for user1, got %d", len(notifications))
	}

	if notifications[0].ID != "test-1" {
		t.Errorf("Expected notification ID 'test-1', got %s", notifications[0].ID)
	}

	// Kiểm tra user không tồn tại
	notifications = nm.GetNotifications("user2")
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notification for user2, got %d", len(notifications))
	}
}

func TestNotificationManager_BroadcastNotificationByRole(t *testing.T) {
	nm := NewNotificationManager()

	// Thiết lập role cho các users
	nm.SetUserRole("doctor1", "doctor")
	nm.SetUserRole("doctor2", "doctor")
	nm.SetUserRole("customer1", "customer")
	nm.SetUserRole("admin1", "admin")

	// Tạo thông báo test
	notification := AppointmentNotification{
		ID:            "test-broadcast-1",
		Title:         "Test Broadcast",
		AppointmentID: 123,
	}

	// Gửi thông báo chỉ đến những người có role "doctor"
	nm.BroadcastNotificationByRole(notification, "doctor")

	// Kiểm tra các doctor nhận được thông báo
	doctorNotifications := nm.GetNotifications("doctor1")
	if len(doctorNotifications) != 1 {
		t.Errorf("Expected 1 notification for doctor1, got %d", len(doctorNotifications))
	}

	doctorNotifications = nm.GetNotifications("doctor2")
	if len(doctorNotifications) != 1 {
		t.Errorf("Expected 1 notification for doctor2, got %d", len(doctorNotifications))
	}

	// Kiểm tra customer không nhận thông báo
	customerNotifications := nm.GetNotifications("customer1")
	if len(customerNotifications) != 0 {
		t.Errorf("Expected 0 notification for customer1, got %d", len(customerNotifications))
	}
}

func TestNotificationManager_WaitForNotificationWithExistingNotification(t *testing.T) {
	nm := NewNotificationManager()

	// Tạo thông báo test
	notification := AppointmentNotification{
		ID:            "test-wait-1",
		Title:         "Test Wait",
		AppointmentID: 123,
	}

	// Thêm thông báo cho user
	nm.AddNotification("user1", notification)

	// Chờ thông báo với timeout
	receivedNotification, hasNotification := nm.WaitForNotification("user1", 1*time.Second)

	if !hasNotification {
		t.Error("Expected to receive a notification, but didn't")
	}

	if receivedNotification.ID != "test-wait-1" {
		t.Errorf("Expected notification ID 'test-wait-1', got %s", receivedNotification.ID)
	}

	// Kiểm tra thông báo đã được xóa khỏi hàng đợi
	notifications := nm.GetNotifications("user1")
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notification left for user1, got %d", len(notifications))
	}
}

func TestNotificationManager_WaitForNotificationWithTimeout(t *testing.T) {
	nm := NewNotificationManager()

	// Chờ thông báo với timeout nhỏ
	_, hasNotification := nm.WaitForNotification("user1", 100*time.Millisecond)

	if hasNotification {
		t.Error("Expected timeout without notification, but received one")
	}
}

func TestNotificationManager_WaitForNotificationAsync(t *testing.T) {
	nm := NewNotificationManager()

	// Tạo một channel để báo hiệu khi test hoàn thành
	done := make(chan bool)

	// Chạy WaitForNotification trong một goroutine
	go func() {
		notification := AppointmentNotification{
			ID:            "test-async",
			Title:         "Test Async",
			AppointmentID: 123,
		}

		// Chờ một chút để đảm bảo goroutine khác đã bắt đầu WaitForNotification
		time.Sleep(100 * time.Millisecond)

		// Thêm thông báo cho user
		nm.AddNotification("user1", notification)
	}()

	go func() {
		// Chờ thông báo với timeout lớn hơn
		receivedNotification, hasNotification := nm.WaitForNotification("user1", 2*time.Second)

		if !hasNotification {
			t.Error("Expected to receive a notification, but didn't")
		} else if receivedNotification.ID != "test-async" {
			t.Errorf("Expected notification ID 'test-async', got %s", receivedNotification.ID)
		}

		done <- true
	}()

	// Chờ test hoàn thành hoặc timeout
	select {
	case <-done:
		// Test hoàn thành thành công
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestNotificationManager_ClearNotifications(t *testing.T) {
	nm := NewNotificationManager()

	// Tạo và thêm một số thông báo
	for i := 0; i < 3; i++ {
		notification := AppointmentNotification{
			ID:            "test-clear",
			Title:         "Test Clear",
			AppointmentID: int64(i),
		}
		nm.AddNotification("user1", notification)
	}

	// Kiểm tra thông báo đã được thêm
	notifications := nm.GetNotifications("user1")
	if len(notifications) != 3 {
		t.Errorf("Expected 3 notifications for user1, got %d", len(notifications))
	}

	// Xóa thông báo
	nm.ClearNotifications("user1")

	// Kiểm tra thông báo đã bị xóa
	notifications = nm.GetNotifications("user1")
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notification after clear, got %d", len(notifications))
	}
}

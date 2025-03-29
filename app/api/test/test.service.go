package test

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type TestServiceInterface interface {
	// CreateTest(ctx *gin.Context, petID, doctorID int64, testType string) (*TestResponse, error)
	// UpdateTestStatus(ctx *gin.Context, testID int64, status string) error
	// AddTestResult(ctx *gin.Context, testID int64, result TestResult) error
	// // GetStatusHistory(ctx *gin.Context, testID int64) ([]StatusHistory, error)
	// GetTestsByPetID(ctx *gin.Context, petID int64) ([]TestResponse, error)
	ListTests(ctx *gin.Context) (*[]Test, error)
	CreateTestOrder(ctx *gin.Context, req TestOrderRequest) error
	GetTestByID(ctx *gin.Context, id int32) (*Test, error)
	UpdateTest(ctx *gin.Context, req UpdateTestRequest) (*Test, error)
	SoftDeleteTest(ctx *gin.Context, testID string) error
}

func NewTestService(store db.Store, ws *websocket.WSClientManager) *TestService {
	return &TestService{
		storeDB: store,
		ws:      ws,
	}
}

// func (s *TestService) CreateTest(ctx *gin.Context, petID, doctorID int64, testType string) (*TestResponse, error) {
// 	test, err := s.storeDB.CreateTest(ctx, db.CreateTestParams{
// 		PetID:    pgtype.Int8{Int64: petID, Valid: true},
// 		DoctorID: pgtype.Int8{Int64: doctorID, Valid: true},
// 		TestType: testType,
// 		Status:   "Pending",
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create test: %w", err)
// 	}

// 	// Notify doctor about new test if connected
// 	doctorClientID := fmt.Sprintf("doctor_%d", doctorID)
// 	if s.ws.IsClientConnected(doctorClientID) {
// 		s.ws.SendToClient(doctorClientID, websocket.WebSocketMessage{
// 			Type:    "test_created",
// 			Message: fmt.Sprintf("New test created for pet #%d", petID),
// 			Data:    test,
// 		})
// 	}

// 	return &TestResponse{
// 		ID:        test.ID,
// 		PetID:     test.PetID.Int64,
// 		DoctorID:  test.DoctorID.Int64,
// 		TestType:  test.TestType,
// 		Status:    test.Status,
// 		CreatedAt: test.CreatedAt.Time,
// 	}, nil
// }

// func (s *TestService) UpdateTestStatus(ctx *gin.Context, testID int64, status string) error {
// 	err := s.storeDB.UpdateTestStatus(ctx, db.UpdateTestStatusParams{
// 		ID:     testID,
// 		Status: status,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to update test status: %w", err)
// 	}

// 	// Get test details to notify relevant parties
// 	test, err := s.storeDB.GetTestByID(ctx, testID)
// 	if err != nil {
// 		return err
// 	}

// 	// Notify both doctor and patient
// 	s.ws.SendToClient(fmt.Sprintf("doctor_%d", test.DoctorID.Int64), websocket.WebSocketMessage{
// 		Type:    "test_status_updated",
// 		Message: fmt.Sprintf("Test #%d status updated to %s", testID, status),
// 	})

// 	return nil
// }

// func (s *TestService) AddTestResult(ctx *gin.Context, testID int64, result TestResult) error {
// 	// Store test result
// 	_, err := s.storeDB.AddTestResult(ctx, db.AddTestResultParams{
// 		TestID:     pgtype.Int8{Int64: testID, Valid: true},
// 		Parameters: []byte(fmt.Sprintf("%v", result.Parameters)),
// 		Notes:      pgtype.Text{String: result.Notes, Valid: true},
// 		Files:      result.Files,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to add test result: %w", err)
// 	}

// 	// Get test details
// 	test, err := s.storeDB.GetTestByID(ctx, testID)
// 	if err != nil {
// 		return err
// 	}

// 	// Send real-time notification
// 	s.ws.SendToClient(fmt.Sprintf("doctor_%d", test.DoctorID.Int64), websocket.WebSocketMessage{
// 		Type:    "test_result_added",
// 		Message: fmt.Sprintf("New results available for test #%d", testID),
// 	})

// 	return nil
// }

// func (s *TestService) GetTestsByPetID(ctx *gin.Context, petID int64) ([]TestResponse, error) {

// 	res, err := s.storeDB.GetTestsByPetID(ctx, pgtype.Int8{Int64: petID, Valid: true})
// 	if err != nil {
// 		return []TestResponse{}, fmt.Errorf("failed to add test result: %w", err)
// 	}

// 	var testResults []TestResponse
// 	for _, r := range res {
// 		testResults = append(testResults, TestResponse{
// 			ID:        r.ID,
// 			PetID:     r.PetID.Int64,
// 			DoctorID:  r.DoctorID.Int64,
// 			TestType:  r.TestType,
// 			Status:    r.Status,
// 			CreatedAt: r.CreatedAt.Time,
// 		})
// 	}
// 	return testResults, nil
// }

func (s *TestService) ListTests(ctx *gin.Context) (*[]Test, error) {
	res, err := s.storeDB.ListTests(ctx)
	if err != nil {
		return &[]Test{}, fmt.Errorf("failed to add test result: %w", err)
	}
	var tests []Test
	for _, r := range res {
		tests = append(tests, Test{
			ID:             r.ID,
			Name:           r.Name,
			Price:          r.Price,
			Description:    r.Description.String,
			TurnaroundTime: r.TurnaroundTime,
			CategoryID:     r.CategoryID.String,
			TestID:         r.TestID,
			IsActive:       r.IsActive.Bool,
			CreatedAt:      r.CreatedAt.Time.String(),
			UpdatedAt:      r.UpdatedAt.Time.String(),
		})
	}
	return &tests, nil
}

func (s *TestService) CreateTestOrder(ctx *gin.Context, req TestOrderRequest) error {
	// Calculate total amount
	var totalAmount float64
	for _, testID := range req.TestIDs {
		test, err := s.storeDB.GetTestByID(ctx, int32(testID))
		if err != nil {
			return fmt.Errorf("failed to get test by ID: %w", err)
		}
		totalAmount += test.Price
	}

	err := s.storeDB.ExecWithTransaction(ctx, func(queries *db.Queries) error {

		order, err := queries.CreateTestOrder(ctx, db.CreateTestOrderParams{
			AppointmentID: pgtype.Int4{Int32: int32(req.AppointmentID), Valid: true},
			TotalAmount:   pgtype.Float8{Float64: totalAmount, Valid: true},
			Notes:         pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
		})
		if err != nil {
			return fmt.Errorf("failed to create test order: %w", err)
		}

		// Add ordered tests
		var orderedTests []OrderedTest
		for _, testID := range req.TestIDs {
			test, _ := queries.GetTestByID(ctx, int32(testID))
			// Add ordered test to database
			_, err = queries.AddOrderedTest(ctx, db.AddOrderedTestParams{
				OrderID:      pgtype.Int4{Int32: order.OrderID, Valid: true},
				TestID:       pgtype.Int4{Int32: int32(testID), Valid: true},
				PriceAtOrder: test.Price,
			})
			if err != nil {
				return fmt.Errorf("failed to add ordered test: %w", err)
			}

			orderedTests = append(orderedTests, OrderedTest{
				TestID:   testID,
				TestName: test.Name,
				Price:    test.Price,
				Status:   "pending",
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *TestService) GetTestByID(ctx *gin.Context, id int32) (*Test, error) {
	test, err := s.storeDB.GetTestByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test by ID: %w", err)
	}

	return &Test{
		ID:             test.ID,
		Name:           test.Name,
		Price:          test.Price,
		Description:    test.Description.String,
		TurnaroundTime: test.TurnaroundTime,
		CategoryID:     test.CategoryID.String,
		TestID:         test.TestID,
		IsActive:       test.IsActive.Bool,
		CreatedAt:      test.CreatedAt.Time.String(),
		UpdatedAt:      test.UpdatedAt.Time.String(),
	}, nil
}

func (s *TestService) UpdateTest(ctx *gin.Context, req UpdateTestRequest) (*Test, error) {
	// Create update params from request
	updateParams := db.UpdateTestParams{
		TestID:         req.TestID,
		Name:           req.Name,
		Description:    pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Price:          req.Price,
		TurnaroundTime: req.TurnaroundTime,
	}

	// Update test in database
	updatedTest, err := s.storeDB.UpdateTest(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update test: %w", err)
	}

	return &Test{
		ID:             updatedTest.ID,
		Name:           updatedTest.Name,
		Price:          updatedTest.Price,
		Description:    updatedTest.Description.String,
		TurnaroundTime: updatedTest.TurnaroundTime,
		CategoryID:     updatedTest.CategoryID.String,
		TestID:         updatedTest.TestID,
		IsActive:       updatedTest.IsActive.Bool,
		CreatedAt:      updatedTest.CreatedAt.Time.String(),
		UpdatedAt:      updatedTest.UpdatedAt.Time.String(),
	}, nil
}

func (s *TestService) SoftDeleteTest(ctx *gin.Context, testID string) error {
	err := s.storeDB.SoftDeleteTest(ctx, testID)
	if err != nil {
		return fmt.Errorf("failed to soft delete test: %w", err)
	}

	return nil
}

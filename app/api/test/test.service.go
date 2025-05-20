package test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type TestServiceInterface interface {
	// Updated methods for combined test/vaccine functionality
	ListItems(ctx *gin.Context, itemType ItemType) (*[]TestCategoryResponse, error)
	CreateOrder(ctx *gin.Context, req TestOrderRequest) (*db.TestOrder, error)
	GetItemByID(ctx *gin.Context, id int32) (*TestItem, error)
	UpdateItem(ctx *gin.Context, req UpdateTestRequest) (*TestItem, error)
	SoftDeleteItem(ctx *gin.Context, itemID string) error
	GetOrderedItemsByAppointment(ctx *gin.Context, appointmentID int64, itemType ItemType) (*[]OrderedItemDetail, error)
	GetAllAppointmentsWithOrders(ctx *gin.Context) (*[]AppointmentWithOrders, error)
	GetTestByAppointment(ctx *gin.Context, appointmentID int64) (*[]TestByAppointment, error)
	// Legacy methods for backward compatibility
	ListTests(ctx *gin.Context) (*[]TestCategoryResponse, error)
	GetTestByID(ctx *gin.Context, id int32) (*Test, error)
	UpdateTest(ctx *gin.Context, req UpdateTestRequest) (*Test, error)
	SoftDeleteTest(ctx *gin.Context, testID string) error
	GetOrderedTestsByAppointment(ctx *gin.Context, appointmentID int64) (*[]OrderedTestDetail, error)
	ListTestCategories(ctx *gin.Context) (*[]TestCategory, error)
}

func NewTestService(store db.Store, ws *websocket.WSClientManager) *TestService {
	return &TestService{
		storeDB: store,
		ws:      ws,
	}
}

// ListItems lists all items (tests and/or vaccines) based on optional type filter
func (s *TestService) ListItems(ctx *gin.Context, itemType ItemType) (*[]TestCategoryResponse, error) {
	// Get all items from database
	res, err := s.storeDB.ListTests(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	// Group items by category
	categoriesMap := make(map[string]*TestCategoryResponse)

	for _, r := range res {
		// Skip type filtering if itemType is "all"
		if itemType != "" && itemType != "all" {
			if ItemType(r.Type.String) != itemType && r.Type.String != "" {
				continue
			}
		}

		categoryID := r.CategoryID.String

		item := TestItem{
			ID:             r.ID,
			Name:           r.Name,
			Price:          r.Price,
			Description:    r.Description.String,
			TurnaroundTime: r.TurnaroundTime,
			CategoryID:     categoryID,
			TestID:         r.TestID,
			Type:           ItemType(r.Type.String),
			IsActive:       r.IsActive.Bool,
			CreatedAt:      r.CreatedAt.Time.String(),
			UpdatedAt:      r.UpdatedAt.Time.String(),
		}

		// If category doesn't exist in map yet, create it
		if _, exists := categoriesMap[categoryID]; !exists {
			category, err := s.storeDB.GetTestCategoryByID(ctx, categoryID)
			if err != nil {
				return nil, fmt.Errorf("failed to get test category by ID: %w", err)
			}

			categoriesMap[categoryID] = &TestCategoryResponse{
				ID:          categoryID,
				Name:        category.Name,
				Icon:        category.IconName.String,
				Description: category.Description.String,
				Items:       []TestItem{},
			}
		}

		// Add item to its category
		categoriesMap[categoryID].Items = append(categoriesMap[categoryID].Items, item)
	}

	// Convert map to slice
	categories := make([]TestCategoryResponse, 0, len(categoriesMap))
	for _, category := range categoriesMap {
		categories = append(categories, *category)
	}

	return &categories, nil
}

// For backward compatibility
func (s *TestService) ListTests(ctx *gin.Context) (*[]TestCategoryResponse, error) {
	return s.ListItems(ctx, TypeTest)
}

// CreateOrder creates a new order for tests or vaccines
func (s *TestService) CreateOrder(ctx *gin.Context, req TestOrderRequest) (*db.TestOrder, error) {
	// Calculate total amount
	var totalAmount float64
	for _, itemID := range req.ItemIDs {
		item, err := s.storeDB.GetTestByID(ctx, int32(itemID))
		if err != nil {
			return nil, fmt.Errorf("failed to get item by ID: %w", err)
		}
		totalAmount += item.Price
	}

	var order db.TestOrder
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(queries *db.Queries) error {
		order, err = queries.CreateTestOrder(ctx, db.CreateTestOrderParams{
			AppointmentID: pgtype.Int4{Int32: int32(req.AppointmentID), Valid: true},
			TotalAmount:   pgtype.Float8{Float64: totalAmount, Valid: true},
			Notes:         pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
		})
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Add ordered items
		for _, itemID := range req.ItemIDs {
			item, _ := queries.GetTestByID(ctx, int32(itemID))
			// Add ordered item to database
			_, err = queries.AddOrderedTest(ctx, db.AddOrderedTestParams{
				OrderID:      pgtype.Int4{Int32: order.OrderID, Valid: true},
				TestID:       pgtype.Int4{Int32: int32(itemID), Valid: true},
				PriceAtOrder: item.Price,
			})
			if err != nil {
				return fmt.Errorf("failed to add ordered item: %w", err)
			}
		}

		return nil
	})

	return &order, err
}

func (s *TestService) GetOrderedTestsByAppointment(ctx *gin.Context, appointmentID int64) (*[]OrderedTestDetail, error) {
	orderedTests, err := s.storeDB.GetOrderedTestsByAppointment(ctx, pgtype.Int4{Int32: int32(appointmentID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered items by appointment: %w", err)
	}

	var res []OrderedTestDetail
	for _, ot := range orderedTests {
		res = append(res, OrderedTestDetail{
			OrderedItemID: int(ot.OrderedTestID),
			ItemID:        ot.TestID,
			ItemName:      ot.TestName,
			CategoryID:    ot.CategoryID.String,
			CategoryName:  ot.CategoryName,
			Price:         ot.PriceAtOrder,
			Status:        ot.Status.String,
			Notes:         ot.Notes.String,
			OrderedDate:   ot.OrderedDate.Time.String(),
		})
	}
	return &res, nil
}

// GetItemByID gets an item (test or vaccine) by ID
func (s *TestService) GetItemByID(ctx *gin.Context, id int32) (*TestItem, error) {
	item, err := s.storeDB.GetTestByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item by ID: %w", err)
	}

	return &TestItem{
		ID:             item.ID,
		Name:           item.Name,
		Price:          item.Price,
		Description:    item.Description.String,
		TurnaroundTime: item.TurnaroundTime,
		CategoryID:     item.CategoryID.String,
		TestID:         item.TestID,
		Type:           ItemType(item.Type.String),
		IsActive:       item.IsActive.Bool,
		CreatedAt:      item.CreatedAt.Time.String(),
		UpdatedAt:      item.UpdatedAt.Time.String(),
	}, nil
}

// For backward compatibility
func (s *TestService) GetTestByID(ctx *gin.Context, id int32) (*Test, error) {
	return s.GetItemByID(ctx, id)
}

// UpdateItem updates an item (test or vaccine)
func (s *TestService) UpdateItem(ctx *gin.Context, req UpdateTestRequest) (*TestItem, error) {
	// Create update params from request
	updateParams := db.UpdateTestParams{
		TestID:         req.TestID,
		Name:           req.Name,
		Description:    pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Price:          req.Price,
		TurnaroundTime: req.TurnaroundTime,
	}

	// Update item in database
	updatedItem, err := s.storeDB.UpdateTest(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &TestItem{
		ID:             updatedItem.ID,
		Name:           updatedItem.Name,
		Price:          updatedItem.Price,
		Description:    updatedItem.Description.String,
		TurnaroundTime: updatedItem.TurnaroundTime,
		CategoryID:     updatedItem.CategoryID.String,
		TestID:         updatedItem.TestID,
		Type:           ItemType(updatedItem.Type.String),
		IsActive:       updatedItem.IsActive.Bool,
		CreatedAt:      updatedItem.CreatedAt.Time.String(),
		UpdatedAt:      updatedItem.UpdatedAt.Time.String(),
	}, nil
}

// For backward compatibility
func (s *TestService) UpdateTest(ctx *gin.Context, req UpdateTestRequest) (*Test, error) {
	// Set default type to test if not provided
	if req.Type == "" {
		req.Type = TypeTest
	}
	return s.UpdateItem(ctx, req)
}

// SoftDeleteItem soft deletes an item (test or vaccine)
func (s *TestService) SoftDeleteItem(ctx *gin.Context, itemID string) error {
	return s.storeDB.SoftDeleteTest(ctx, itemID)
}

// For backward compatibility
func (s *TestService) SoftDeleteTest(ctx *gin.Context, testID string) error {
	return s.SoftDeleteItem(ctx, testID)
}

func (s *TestService) GetAllAppointmentsWithOrders(ctx *gin.Context) (*[]AppointmentWithOrders, error) {
	appointments, err := s.storeDB.GetAllAppointmentsWithOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all appointments with orders: %w", err)
	}

	var response []AppointmentWithOrders
	for _, a := range appointments {
		var orders []TestOrderDetail

		// Handle different types of Orders data
		switch v := a.Orders.(type) {
		case []byte:
			if err := json.Unmarshal(v, &orders); err != nil {
				continue
			}
		case string:
			if err := json.Unmarshal([]byte(v), &orders); err != nil {
				continue
			}
		case []interface{}:
			// Handle case where Orders is a slice of interface{}
			for _, item := range v {
				if orderMap, ok := item.(map[string]interface{}); ok {
					order := TestOrderDetail{}

					// Extract order_id
					if orderID, ok := orderMap["order_id"].(float64); ok {
						order.OrderID = int(orderID)
					}

					// Extract status
					if status, ok := orderMap["status"].(string); ok {
						order.Status = status
					}

					// Extract total_amount
					if amount, ok := orderMap["total_amount"].(float64); ok {
						order.TotalAmount = amount
					}

					// Extract order_date
					if date, ok := orderMap["order_date"].(string); ok {
						order.OrderDate = date
					}

					// Extract tests array
					if testsData, ok := orderMap["tests"].([]interface{}); ok {
						var tests []OrderedTest
						for _, testItem := range testsData {
							if testMap, ok := testItem.(map[string]interface{}); ok {
								test := OrderedTest{}

								// Extract test_id
								if testID, ok := testMap["test_id"].(string); ok {
									intTestID, err := strconv.Atoi(testID)
									if err != nil {
										continue
									}
									test.ItemID = intTestID
								}

								// Extract test_name
								if testName, ok := testMap["test_name"].(string); ok {
									test.ItemName = testName
								}

								// Extract price
								if price, ok := testMap["price"].(float64); ok {
									test.Price = price
								}

								// Extract status
								if status, ok := testMap["status"].(string); ok {
									test.Status = status
								}

								tests = append(tests, test)
							}
						}
						order.Items = tests
					}

					orders = append(orders, order)
				}
			}
		default:
			// Try to marshal and unmarshal for other types
			data, err := json.Marshal(a.Orders)
			if err == nil {
				_ = json.Unmarshal(data, &orders)
			}
		}

		// Only add to response if there are orders
		if len(orders) > 0 {
			response = append(response, AppointmentWithOrders{
				AppointmentID:   int(a.AppointmentID),
				AppointmentDate: a.Date.Time.Format("2006-01-02"),
				PetName:         a.PetName,
				Orders:          orders,
			})
		}
	}
	return &response, nil
}

// GetOrderedItemsByAppointment gets ordered items by appointment ID with optional type filter
func (s *TestService) GetOrderedItemsByAppointment(ctx *gin.Context, appointmentID int64, itemType ItemType) (*[]OrderedItemDetail, error) {
	// Get ordered items from database
	orderedItems, err := s.storeDB.GetOrderedTestsByAppointment(ctx, pgtype.Int4{Int32: int32(appointmentID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered items by appointment: %w", err)
	}

	var res []OrderedItemDetail
	for _, item := range orderedItems {
		// Convert string ID to int32
		testID, err := strconv.ParseInt(item.TestID, 10, 32)
		if err != nil {
			continue // Skip this item if we can't parse the ID
		}

		// Get the item type from the test record
		test, err := s.storeDB.GetTestByID(ctx, int32(testID))
		if err != nil {
			continue // Skip this item if we can't get its type
		}

		// If itemType filter is provided, skip items that don't match
		if itemType != "" {
			if ItemType(test.Type.String) != itemType && test.Type.String != "" {
				continue
			}
		}

		res = append(res, OrderedItemDetail{
			OrderedItemID: int(item.OrderedTestID),
			ItemID:        item.TestID,
			ItemName:      item.TestName,
			ItemType:      ItemType(test.Type.String),
			CategoryID:    item.CategoryID.String,
			CategoryName:  item.CategoryName,
			Price:         item.PriceAtOrder,
			Status:        item.Status.String,
			Notes:         item.Notes.String,
			OrderedDate:   item.OrderedDate.Time.String(),
		})
	}
	return &res, nil
}

func (s *TestService) GetTestByAppointment(ctx *gin.Context, appointmentID int64) (*[]TestByAppointment, error) {
	var tests []TestByAppointment

	results, err := s.storeDB.ExecStatementMany(ctx, `
			SELECT DISTINCT ON (t.test_id)
				t.test_id AS test_identifier,
				t.name AS test_name,
				t.description AS test_description,
				t.price AS current_price,
				ot.price_at_order AS ordered_price,
				t.turnaround_time,
				t.type AS test_type,
				t.expiration_date,
				t.batch_number,
				tc.name AS category_name,
				tc.description AS category_description,
				tor.order_id,
				tor.order_date,
				tor.status AS order_status,
				ot.status AS test_status
			FROM test_orders tor
			JOIN ordered_tests ot ON tor.order_id = ot.order_id
			JOIN tests t ON ot.test_id = t.id
			LEFT JOIN test_categories tc ON t.category_id = tc.category_id
			WHERE tor.appointment_id = $1 AND t.type = 'vaccine'
			ORDER BY t.test_id, tor.order_date DESC;
	`, []interface{}{appointmentID})

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var test TestByAppointment

		fmt.Println(result["expiration_date"])

		if testName, ok := result["test_name"].(string); ok {
			test.TestName = testName
		}

		if testIdentifier, ok := result["test_identifier"].(string); ok {
			test.TestID = testIdentifier
		}

		if batchNumber, ok := result["batch_number"].(string); ok {
			test.BatchNumber = batchNumber
		}

		// Extract and set the expiration_date field
		if expirationDate, ok := result["expiration_date"].(time.Time); ok {
			test.ExpirationDate = expirationDate.Format("2006-01-02")
		} else if expirationDateStr, ok := result["expiration_date"].(string); ok {
			test.ExpirationDate = expirationDateStr
		}

		tests = append(tests, test)
	}

	return &tests, nil
}

// ListTestCategories lists all test categories
func (s *TestService) ListTestCategories(ctx *gin.Context) (*[]TestCategory, error) {
	categories, err := s.storeDB.ListTestCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list test categories: %w", err)
	}

	var res []TestCategory
	for _, category := range categories {
		res = append(res, TestCategory{
			ID:          category.CategoryID,
			Name:        category.Name,
			Description: category.Description.String,
		})
	}

	return &res, nil
}

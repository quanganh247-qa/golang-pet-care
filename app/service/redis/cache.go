package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

func (c *ClientType) LoadCacheByKey(key string, result interface{}, duration time.Duration) (cacheData func(interface{})) {
	err := c.GetWithBackground(key, result) // Attempt to load data from cache
	cacheData = func(i interface{}) {       // Define a closure to set cache
		err = c.SetWithBackground(key, i, duration) // Store new data in cache
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", key, err) // Log error if setting fails
		}
	}
	return // Return closure
}

type UserInfo struct {
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	HashedPassword  string `json:"hashed_password"`
	PhoneNumber     string `json:"phone_number"`
	Address         string `json:"address"`
	Role            string `json:"role"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	DataImage       string `json:"data_image"`
	OriginalImage   string `json:"original_image"`
}

func (c *ClientType) UserInfoLoadCache(username string) (*UserInfo, error) {
	userKey := fmt.Sprintf("%s:%s", USER_INFO_KEY, username)
	// log.Printf("User key: %s", userKey)
	userInformation := UserInfo{}
	err := c.GetWithBackground(userKey, &userInformation)
	if err != nil {
		// log.Printf("Error when get cache for key %s: %v", userKey, err)
		userData, err := db.StoreDB.GetUser(ctxRedis, username)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("Không tìm thấy user với username = %s", username)
			}
			return nil, err
		}

		userRes := UserInfo{
			UserID:          userData.ID,
			Username:        userData.Username,
			Email:           userData.Email,
			HashedPassword:  userData.HashedPassword,
			IsVerifiedEmail: userData.IsVerifiedEmail.Bool,
			FullName:        userData.FullName,
			DataImage:       string(userData.DataImage),
			OriginalImage:   userData.OriginalImage.String,
			PhoneNumber:     userData.PhoneNumber.String,
			Address:         userData.Address.String,
			Role:            userData.Role.String,
		}
		err = c.SetWithBackground(userKey, &userRes, time.Hour*12)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", userKey, err)
		}
		return &userRes, nil
	}
	return &userInformation, nil
}

func (client *ClientType) RemoveUserInfoCache(username string) {
	userInfoKey := fmt.Sprintf("%s:%s", USER_INFO_KEY, username)
	client.RemoveCacheByKey(userInfoKey)
}

func (client *ClientType) ClearUserInfoCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", USER_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

type PetInfo struct {
	Petid           int64   `json:"petid"`
	Username        string  `json:"username"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"` // Health notes of pet, for example: "Vaccinated", "Neutered", "Dewormed", "Flea treatment", "Veterinary treatment", "Parasite treatment", "Disease treatment", "Other treatment", "No treatment"
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	BOD             string  `json:"birth_date"`
	Weight          float64 `json:"weight"`
	DataImage       []byte  `json:"data_image"`
	OriginalImage   string  `json:"original_name"`
	MicrochipNumber string  `json:"microchip_number"`
}

func (c *ClientType) PetInfoLoadCache(petid int64) (*PetInfo, error) {
	petKey := fmt.Sprintf("%s:%d", PET_INFO_KEY, petid)
	petInfo := PetInfo{}
	err := c.GetWithBackground(petKey, &petInfo)
	if err != nil {
		res, err := db.StoreDB.GetPetByID(ctxRedis, petid)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
		pet := PetInfo{
			Petid:           res.Petid,
			Username:        res.Username,
			Name:            res.Name,
			Type:            res.Type,
			Breed:           res.Breed.String,
			Gender:          res.Gender.String,
			Healthnotes:     res.Healthnotes.String,
			MicrochipNumber: res.MicrochipNumber.String,
			BOD:             res.BirthDate.Time.Format("2006-01-02"),
			Age:             int16(res.Age.Int32),
			Weight:          res.Weight.Float64,
			DataImage:       res.DataImage,
			OriginalImage:   res.OriginalImage.String,
		}
		err = c.SetWithBackground(petKey, &pet, time.Hour*12)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", petKey, err)
		}
		return &pet, nil
	}
	return &petInfo, nil
}

func (client *ClientType) RemovePetInfoCache(petid int64) {
	petKey := fmt.Sprintf("%s:%d", PET_INFO_KEY, petid)
	client.RemoveCacheByKey(petKey)
	fmt.Println("Remove cache for key: ", petKey)
}

func (client *ClientType) ClearPetInfoCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", PET_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

type ProductInfo struct {
	ProductID     int64   `json:"product_id"`
	Description   string  `json:"description"`
	IsAvailable   bool    `json:"is_available"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Stock         int32   `json:"stock"`
	Category      string  `json:"category"`
	DataImage     []byte  `json:"data_image"`
	OriginalImage string  `json:"original_image"`
}

func (c *ClientType) ProductInfoLoadCache(productID int64) (*ProductInfo, error) {
	productKey := fmt.Sprintf("%s:%d", PRODUCT_INFO_KEY, productID)
	productInfo := ProductInfo{}
	err := c.GetWithBackground(productKey, &productInfo)
	if err != nil {
		// Cache miss, get from database
		res, err := db.StoreDB.GetProductByID(ctxRedis, productID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
		product := ProductInfo{
			ProductID:     res.ProductID,
			Name:          res.Name,
			Description:   res.Description.String,
			IsAvailable:   res.IsAvailable.Bool,
			Price:         res.Price,
			Stock:         res.StockQuantity.Int32,
			Category:      res.Category.String,
			DataImage:     res.DataImage,
			OriginalImage: res.OriginalImage.String,
		}
		// Store in cache for future requests
		err = c.SetWithBackground(productKey, &product, time.Hour*12)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", productKey, err)
		}
		return &product, nil
	}
	return &productInfo, nil
}

func (client *ClientType) RemoveProductInfoCache(productID int64) {
	productKey := fmt.Sprintf("%s:%d", PRODUCT_INFO_KEY, productID)
	client.RemoveCacheByKey(productKey)
}

func (client *ClientType) ClearProductInfoCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", PRODUCT_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

// Get all products with cache
func (c *ClientType) ProductsListLoadCache(page int32, pageSize int32) ([]ProductInfo, error) {
	// Create a key that includes pagination parameters
	listKey := fmt.Sprintf("%s:list:%d:%d", PRODUCT_INFO_KEY, page, pageSize)

	var productsList []ProductInfo
	err := c.GetWithBackground(listKey, &productsList)
	if err != nil {
		// Cache miss, get from database
		offset := (page - 1) * pageSize
		products, err := db.StoreDB.GetAllProducts(ctxRedis, db.GetAllProductsParams{
			Limit:  pageSize,
			Offset: offset,
		})
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, product := range products {
			productsList = append(productsList, ProductInfo{
				ProductID:     product.ProductID,
				Name:          product.Name,
				Price:         product.Price,
				Stock:         product.StockQuantity.Int32,
				Category:      product.Category.String,
				DataImage:     product.DataImage,
				OriginalImage: product.OriginalImage.String,
			})
		}

		// Store in cache for future requests - using shorter TTL for lists
		err = c.SetWithBackground(listKey, &productsList, time.Hour*1)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", listKey, err)
		}
	}

	return productsList, nil
}

func (client *ClientType) RemoveProductListCache(page int32, pageSize int32) {
	listKey := fmt.Sprintf("%s:list:%d:%d", PRODUCT_INFO_KEY, page, pageSize)
	client.RemoveCacheByKey(listKey)
}

// PetLogInfo is a simplified structure for caching pet logs
type PetLogInfo struct {
	LogID    int64     `json:"log_id"`
	PetID    int64     `json:"pet_id"`
	Title    string    `json:"title"`
	Notes    string    `json:"notes"`
	Datetime time.Time `json:"date_time"`
	PetName  string    `json:"pet_name"`
	PetType  string    `json:"pet_type"`
	PetBreed string    `json:"pet_breed"`
}

// Load a single pet log from cache or DB
func (c *ClientType) PetLogLoadCache(logID int64) (*PetLogInfo, error) {
	logKey := fmt.Sprintf("%s:%d", PET_LOG_KEY, logID)
	petLogInfo := PetLogInfo{}
	err := c.GetWithBackground(logKey, &petLogInfo)
	if err != nil {
		// Cache miss, get from database
		log, err := db.StoreDB.GetPetLogByID(ctxRedis, logID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}

		// Convert database result to cache structure
		logInfo := PetLogInfo{
			LogID:    log.LogID,
			PetID:    log.Petid,
			Title:    log.Title.String,
			Notes:    log.Notes.String,
			Datetime: log.Datetime.Time,
			PetName:  log.Name.String,
			PetType:  log.Type.String,
			PetBreed: log.Breed.String,
		}

		// Cache for 6 hours
		err = c.SetWithBackground(logKey, &logInfo, time.Hour*6)
		if err != nil {
			fmt.Printf("Error when set cache for key %s: %v", logKey, err)
		}
		return &logInfo, nil
	}
	return &petLogInfo, nil
}

// Load all pet logs for a pet from cache or DB
func (c *ClientType) PetLogsListByPetIDLoadCache(petID int64) ([]PetLogInfo, error) {
	listKey := fmt.Sprintf("%s:pet:%d", PET_LOG_KEY, petID)

	var petLogs []PetLogInfo
	err := c.GetWithBackground(listKey, &petLogs)
	if err != nil {
		// Cache miss, get from database
		logs, err := db.StoreDB.GetPetLogsByPetID(ctxRedis, db.GetPetLogsByPetIDParams{
			Petid: petID,
		})
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, log := range logs {
			petLogs = append(petLogs, PetLogInfo{
				LogID:    log.LogID,
				PetID:    log.Petid,
				Title:    log.Title.String,
				Notes:    log.Notes.String,
				Datetime: log.Datetime.Time,
				PetName:  log.Name.String,
				PetType:  log.Type.String,
				PetBreed: log.Breed.String,
			})
		}

		// Cache for 1 hour
		err = c.SetWithBackground(listKey, &petLogs, time.Hour*1)
		if err != nil {
			fmt.Printf("Error when set cache for key %s: %v", listKey, err)
		}
	}

	return petLogs, nil
}

// Load all pet logs for a user from cache or DB
func (c *ClientType) PetLogsListByUsernameLoadCache(username string) ([]PetLogInfo, error) {
	listKey := fmt.Sprintf("%s:user:%s", PET_LOG_KEY, username)

	var petLogs []PetLogInfo
	err := c.GetWithBackground(listKey, &petLogs)
	if err != nil {
		// Cache miss, get from database
		logs, err := db.StoreDB.GetAllPetLogsByUsername(ctxRedis, db.GetAllPetLogsByUsernameParams{
			Username: username,
		})
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, log := range logs {
			petLogs = append(petLogs, PetLogInfo{
				LogID:    log.LogID,
				PetID:    log.Petid,
				Title:    log.Title.String,
				Notes:    log.Notes.String,
				Datetime: log.Datetime.Time,
				PetName:  log.PetName,
				PetType:  log.PetType,
				PetBreed: log.PetBreed.String,
			})
		}

		// Cache for 1 hour
		err = c.SetWithBackground(listKey, &petLogs, time.Hour*1)
		if err != nil {
			fmt.Printf("Error when set cache for key %s: %v", listKey, err)
		}
	}

	return petLogs, nil
}

// Remove a single pet log from cache
func (client *ClientType) RemovePetLogCache(logID int64) {
	logKey := fmt.Sprintf("%s:%d", PET_LOG_KEY, logID)
	client.RemoveCacheByKey(logKey)
}

// Clear all pet logs for a pet
func (client *ClientType) ClearPetLogsByPetCache(petID int64) {
	petKey := fmt.Sprintf("%s:pet:%d", PET_LOG_KEY, petID)
	client.RemoveCacheByKey(petKey)
}

// Clear all pet logs for a user
func (client *ClientType) ClearPetLogsByUserCache(username string) {
	userKey := fmt.Sprintf("%s:user:%s", PET_LOG_KEY, username)
	client.RemoveCacheByKey(userKey)
}

// Clear all pet logs cache
func (client *ClientType) ClearPetLogsCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", PET_LOG_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			fmt.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

// PetScheduleInfo represents cached pet schedule data
type PetScheduleInfo struct {
	ID               int64     `json:"id"`
	PetID            int64     `json:"pet_id"`
	Title            string    `json:"title"`
	ReminderDateTime time.Time `json:"reminder_datetime"`
	EventRepeat      string    `json:"event_repeat"`
	EndType          bool      `json:"end_type"`
	EndDate          time.Time `json:"end_date"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	IsActive         bool      `json:"is_active"`
}

// Load a single pet schedule from cache or DB
func (c *ClientType) PetScheduleLoadCache(scheduleID int64) (*PetScheduleInfo, error) {
	scheduleKey := fmt.Sprintf("%s:%d", PET_SCHEDULE_KEY, scheduleID)
	scheduleInfo := PetScheduleInfo{}
	err := c.GetWithBackground(scheduleKey, &scheduleInfo)
	if err != nil {
		// Cache miss, get from database
		schedule, err := db.StoreDB.GetPetScheduleById(ctxRedis, scheduleID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}

		// Convert database result to cache structure
		reminderTime, _ := time.Parse(time.RFC3339, schedule.ReminderDatetime.Time.Format(time.RFC3339))
		endDate, _ := time.Parse(time.RFC3339, schedule.EndDate.Time.Format(time.RFC3339))
		createdAt, _ := time.Parse(time.RFC3339, schedule.CreatedAt.Time.Format(time.RFC3339))

		scheduleInfo := PetScheduleInfo{
			ID:               schedule.ID,
			PetID:            schedule.PetID.Int64,
			Title:            schedule.Title.String,
			ReminderDateTime: reminderTime,
			EventRepeat:      schedule.EventRepeat.String,
			EndType:          schedule.EndType.Bool,
			EndDate:          endDate,
			Notes:            schedule.Notes.String,
			CreatedAt:        createdAt,
			IsActive:         schedule.IsActive.Bool,
		}

		// Cache for 6 hours
		err = c.SetWithBackground(scheduleKey, &scheduleInfo, time.Hour*6)
		if err != nil {
			fmt.Printf("Error when set cache for key %s: %v", scheduleKey, err)
		}
		return &scheduleInfo, nil
	}
	return &scheduleInfo, nil
}

// Load all schedules for a pet from cache or DB
func (c *ClientType) PetSchedulesByPetIDLoadCache(petID int64) ([]PetScheduleInfo, error) {
	listKey := fmt.Sprintf("%s:pet:%d", PET_SCHEDULE_KEY, petID)

	var schedules []PetScheduleInfo
	err := c.GetWithBackground(listKey, &schedules)
	if err != nil {
		// Cache miss, get from database
		dbSchedules, err := db.StoreDB.ListPetSchedulesByPetID(ctxRedis, db.ListPetSchedulesByPetIDParams{
			PetID: pgtype.Int8{Int64: petID, Valid: true},
		})
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, schedule := range dbSchedules {
			// Parse dates
			reminderTime, _ := time.Parse(time.RFC3339, schedule.ReminderDatetime.Time.Format(time.RFC3339))
			endDate, _ := time.Parse(time.RFC3339, schedule.EndDate.Time.Format(time.RFC3339))
			createdAt, _ := time.Parse(time.RFC3339, schedule.CreatedAt.Time.Format(time.RFC3339))

			schedules = append(schedules, PetScheduleInfo{
				ID:               schedule.ID,
				PetID:            schedule.PetID.Int64,
				Title:            schedule.Title.String,
				ReminderDateTime: reminderTime,
				EventRepeat:      schedule.EventRepeat.String,
				EndType:          schedule.EndType.Bool,
				EndDate:          endDate,
				Notes:            schedule.Notes.String,
				CreatedAt:        createdAt,
				IsActive:         schedule.IsActive.Bool,
			})
		}

		// Cache for 1 hour
		err = c.SetWithBackground(listKey, &schedules, time.Hour*1)
		if err != nil {
			fmt.Printf("Error when set cache for key %s: %v", listKey, err)
		}
	}

	return schedules, nil
}

// Remove a single pet schedule from cache
func (client *ClientType) RemovePetScheduleCache(scheduleID int64) {
	scheduleKey := fmt.Sprintf("%s:%d", PET_SCHEDULE_KEY, scheduleID)
	client.RemoveCacheByKey(scheduleKey)
}

// Clear all schedules for a pet
func (client *ClientType) ClearPetSchedulesByPetCache(petID int64) {
	petKey := fmt.Sprintf("%s:pet:%d", PET_SCHEDULE_KEY, petID)
	client.RemoveCacheByKey(petKey)
}

// Clear all pet schedules cache
func (client *ClientType) ClearPetSchedulesCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", PET_SCHEDULE_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			fmt.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

// GetPetsByUsernameCache loads a list of pets by username from the cache or from the database
func (c *ClientType) GetPetsByUsernameCache(username string) ([]PetInfo, error) {
	cacheKey := fmt.Sprintf("%s:by-user:%s", PET_INFO_KEY, username)
	var petsList []PetInfo

	err := c.GetWithBackground(cacheKey, &petsList)
	if err != nil {
		// Cache miss, get from database
		pets, err := db.StoreDB.ListPetsByUsername(ctxRedis, db.ListPetsByUsernameParams{
			Username: username,
			Limit:    100,
			Offset:   0,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting pets from database: %w", err)
		}

		petsList = make([]PetInfo, 0, len(pets))
		for _, pet := range pets {
			petInfo := PetInfo{
				Petid:           pet.Petid,
				Username:        pet.Username,
				Name:            pet.Name,
				Type:            pet.Type,
				Breed:           pet.Breed.String,
				Gender:          pet.Gender.String,
				Healthnotes:     pet.Healthnotes.String,
				MicrochipNumber: pet.MicrochipNumber.String,
				BOD:             pet.BirthDate.Time.Format("2006-01-02"),
				Age:             int16(pet.Age.Int32),
				Weight:          pet.Weight.Float64,
				DataImage:       pet.DataImage,
				OriginalImage:   pet.OriginalImage.String,
			}
			petsList = append(petsList, petInfo)
		}

		// Store in cache for future requests (30 minute cache)
		err = c.SetWithBackground(cacheKey, petsList, time.Minute*30)
		if err != nil {
			log.Printf("Error caching pets list for user %s: %v", username, err)
		}
	}

	return petsList, nil
}

// ClearUserPetsCache clears the cache for a user's pets
func (client *ClientType) ClearUserPetsCache(username string) {
	userPetsKey := fmt.Sprintf("%s:by-user:%s", PET_INFO_KEY, username)
	client.RemoveCacheByKey(userPetsKey)
}

// PetSchedulesByUsernameLoadCache loads all pet schedules for a username
func (c *ClientType) PetSchedulesByUsernameLoadCache(username string) (map[int64][]PetScheduleInfo, error) {
	cacheKey := fmt.Sprintf("%s:by-user:%s", PET_SCHEDULE_KEY, username)
	var schedulesByPet map[int64][]PetScheduleInfo

	err := c.GetWithBackground(cacheKey, &schedulesByPet)
	if err != nil {
		// Cache miss, get from database
		// First get all pets for the user
		pets, err := db.StoreDB.ListPetsByUsername(ctxRedis, db.ListPetsByUsernameParams{
			Username: username,
			Limit:    100,
			Offset:   0,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting pets from database: %w", err)
		}

		schedulesByPet = make(map[int64][]PetScheduleInfo)

		// For each pet, get its schedules
		for _, pet := range pets {
			petID := pet.Petid
			schedules, err := db.StoreDB.ListPetSchedulesByPetID(ctxRedis, db.ListPetSchedulesByPetIDParams{
				PetID:  pgtype.Int8{Int64: petID, Valid: true},
				Limit:  100,
				Offset: 0,
			})

			if err != nil {
				log.Printf("Error getting schedules for pet %d: %v", petID, err)
				continue
			}

			petSchedules := make([]PetScheduleInfo, 0, len(schedules))
			for _, s := range schedules {
				scheduleInfo := PetScheduleInfo{
					ID:               s.ID,
					PetID:            s.PetID.Int64,
					Title:            s.Title.String,
					ReminderDateTime: s.ReminderDatetime.Time,
					EventRepeat:      s.EventRepeat.String,
					EndType:          s.EndType.Bool,
					EndDate:          s.EndDate.Time,
					Notes:            s.Notes.String,
					CreatedAt:        s.CreatedAt.Time,
					IsActive:         s.IsActive.Bool,
				}
				petSchedules = append(petSchedules, scheduleInfo)
			}

			schedulesByPet[petID] = petSchedules
		}

		// Store in cache for future requests (15 minute cache)
		err = c.SetWithBackground(cacheKey, schedulesByPet, time.Minute*15)
		if err != nil {
			log.Printf("Error caching schedules for user %s: %v", username, err)
		}
	}

	return schedulesByPet, nil
}

// ClearUserSchedulesCache clears all schedule caches for a user
func (client *ClientType) ClearUserSchedulesCache(username string) {
	userSchedulesKey := fmt.Sprintf("%s:by-user:%s", PET_SCHEDULE_KEY, username)
	client.RemoveCacheByKey(userSchedulesKey)
}

// PetLogSummaryByPetIDLoadCache loads a summary of recent pet logs
func (c *ClientType) PetLogSummaryByPetIDLoadCache(petID int64, limit int32) ([]PetLogInfo, error) {
	cacheKey := fmt.Sprintf("%s:summary:%d:%d", PET_LOG_KEY, petID, limit)
	var logSummary []PetLogInfo

	err := c.GetWithBackground(cacheKey, &logSummary)
	if err != nil {
		// Cache miss, get from database
		logs, err := db.StoreDB.GetPetLogsByPetID(ctxRedis, db.GetPetLogsByPetIDParams{
			Petid:  petID,
			Limit:  limit,
			Offset: 0,
		})

		if err != nil {
			return nil, fmt.Errorf("error getting log summary from database: %w", err)
		}

		logSummary = make([]PetLogInfo, 0, len(logs))
		for _, log := range logs {
			logInfo := PetLogInfo{
				LogID:    log.LogID,
				PetID:    log.Petid,
				Title:    log.Title.String,
				Notes:    log.Notes.String,
				Datetime: log.Datetime.Time,
				PetName:  log.Name.String,
				PetType:  log.Type.String,
				PetBreed: log.Breed.String,
			}
			logSummary = append(logSummary, logInfo)
		}

		// Store in cache for future requests (10 minute cache)
		err = c.SetWithBackground(cacheKey, logSummary, time.Minute*10)
		if err != nil {
			log.Printf("Error caching log summary for pet %d: %v", petID, err)
		}
	}

	return logSummary, nil
}

// ClearPetLogSummaryCache clears pet log summary cache
func (client *ClientType) ClearPetLogSummaryCache(petID int64) {
	// Clear all summary caches with different limits
	pattern := fmt.Sprintf("%s:summary:%d:*", PET_LOG_KEY, petID)
	client.RemoveCacheBySubString(pattern)
}

// DoctorInfo represents cached doctor data
type DoctorInfo struct {
	DoctorID       int64  `json:"doctor_id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	DoctorName     string `json:"doctor_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Specialization string `json:"specialization"`
	YearsOfExp     int32  `json:"years_of_exp"`
	Education      string `json:"education"`
	Certificate    string `json:"certificate"`
	Bio            string `json:"bio"`
	DataImage      []byte `json:"data_image"`
}

// ShiftInfo represents cached doctor shift data
type ShiftInfo struct {
	ID         int64     `json:"id"`
	DoctorID   int64     `json:"doctor_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Date       string    `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	DoctorName string    `json:"doctor_name,omitempty"`
}

// Load a single doctor from cache or DB
func (c *ClientType) DoctorLoadCache(doctorID int64) (*DoctorInfo, error) {
	doctorKey := fmt.Sprintf("%s:%d", DOCTOR_INFO_KEY, doctorID)
	doctorInfo := DoctorInfo{}
	err := c.GetWithBackground(doctorKey, &doctorInfo)
	if err != nil {
		// Cache miss, get from database
		doctor, err := db.StoreDB.GetDoctor(ctxRedis, doctorID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}

		// Get user details for the doctor
		user, err := db.StoreDB.GetUserByID(ctxRedis, doctor.UserID)
		if err != nil {
			return nil, err
		}

		// Convert database result to cache structure
		doctorData := DoctorInfo{
			DoctorID:       doctor.ID,
			Username:       user.Username,
			FullName:       user.FullName,
			DoctorName:     user.FullName,
			Email:          user.Email,
			Role:           user.Role.String,
			Address:        user.Address.String,
			PhoneNumber:    user.PhoneNumber.String,
			Specialization: doctor.Specialization.String,
			YearsOfExp:     doctor.YearsOfExperience.Int32,
			Education:      doctor.Education.String,
			Certificate:    doctor.CertificateNumber.String,
			Bio:            doctor.Bio.String,
			DataImage:      []byte(user.DataImage),
		}

		// Cache for 6 hours
		err = c.SetWithBackground(doctorKey, &doctorData, time.Hour*6)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", doctorKey, err)
		}
		return &doctorData, nil
	}
	return &doctorInfo, nil
}

// Load doctor by username from cache or DB
func (c *ClientType) DoctorByUsernameLoadCache(username string) (*DoctorInfo, error) {
	usernameKey := fmt.Sprintf("%s:username:%s", DOCTOR_INFO_KEY, username)
	doctorInfo := DoctorInfo{}
	err := c.GetWithBackground(usernameKey, &doctorInfo)
	if err != nil {
		// Cache miss, get from database
		user, err := db.StoreDB.GetUser(ctxRedis, username)
		if err != nil {
			return nil, err
		}

		doctor, err := db.StoreDB.GetDoctorByUserId(ctxRedis, user.ID)
		if err != nil {
			return nil, err
		}

		// Convert database result to cache structure
		doctorData := DoctorInfo{
			DoctorID:       doctor.ID,
			Username:       user.Username,
			FullName:       user.FullName,
			DoctorName:     user.FullName,
			Email:          user.Email,
			Role:           user.Role.String,
			Address:        user.Address.String,
			PhoneNumber:    user.PhoneNumber.String,
			Specialization: doctor.Specialization.String,
			YearsOfExp:     doctor.YearsOfExperience.Int32,
			Education:      doctor.Education.String,
			Certificate:    doctor.CertificateNumber.String,
			Bio:            doctor.Bio.String,
			DataImage:      []byte(user.DataImage),
		}

		// Cache for 6 hours
		err = c.SetWithBackground(usernameKey, &doctorData, time.Hour*6)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", usernameKey, err)
		}
		return &doctorData, nil
	}
	return &doctorInfo, nil
}

// Get all doctors with cache
func (c *ClientType) DoctorsListLoadCache() ([]DoctorInfo, error) {
	listKey := fmt.Sprintf("%s:list", DOCTOR_INFO_KEY)

	var doctorsList []DoctorInfo
	err := c.GetWithBackground(listKey, &doctorsList)
	if err != nil {
		// Cache miss, get from database
		doctors, err := db.StoreDB.ListDoctors(ctxRedis)
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, doc := range doctors {
			doctorsList = append(doctorsList, DoctorInfo{
				DoctorID:       doc.DoctorID,
				Username:       doc.Username,
				FullName:       doc.FullName,
				DoctorName:     doc.FullName,
				Email:          doc.Email,
				Role:           doc.Role.String,
				Specialization: doc.Specialization.String,
				YearsOfExp:     doc.YearsOfExperience.Int32,
				Education:      doc.Education.String,
				Certificate:    doc.CertificateNumber.String,
				Bio:            doc.Bio.String,
				DataImage:      doc.DataImage,
			})
		}

		// Cache for 1 hour
		err = c.SetWithBackground(listKey, &doctorsList, time.Hour*1)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", listKey, err)
		}
	}

	return doctorsList, nil
}

// Get shifts for a doctor with cache
func (c *ClientType) DoctorShiftsLoadCache(doctorID int64) ([]ShiftInfo, error) {
	shiftsKey := fmt.Sprintf("%s:%d:shifts", DOCTOR_INFO_KEY, doctorID)

	var shifts []ShiftInfo
	err := c.GetWithBackground(shiftsKey, &shifts)
	if err != nil {
		// Cache miss, get from database
		dbShifts, err := db.StoreDB.GetShiftByDoctorId(ctxRedis, doctorID)
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, shift := range dbShifts {
			// Get doctor name for each shift
			doctor, err := db.StoreDB.GetDoctor(ctxRedis, shift.DoctorID)
			if err != nil {
				// Continue even if we can't get doctor name
				// Type assertion for StartTime and EndTime
				var startTime, endTime time.Time

				// Without doctor name, just set default start/end times
				startTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
				endTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

				shifts = append(shifts, ShiftInfo{
					ID:        shift.ID,
					DoctorID:  shift.DoctorID,
					StartTime: startTime,
					EndTime:   endTime,
					Date:      shift.Date.Time.Format("2006-01-02"),
					CreatedAt: shift.CreatedAt.Time,
				})
				continue
			}

			// Declare time variables
			var startTime, endTime time.Time

			// Parse time from pgtype.Time
			// Convert microseconds from pgtype.Time to hours, minutes, seconds
			if shift.StartTime.Valid {
				startTime = time.Date(2000, 1, 1,
					int(shift.StartTime.Microseconds/1000000/3600),            // Hour
					int((shift.StartTime.Microseconds/1000000%3600)/60),       // Minute
					int(shift.StartTime.Microseconds/1000000%60), 0, time.UTC) // Second
			}

			if shift.EndTime.Valid {
				endTime = time.Date(2000, 1, 1,
					int(shift.EndTime.Microseconds/1000000/3600),            // Hour
					int((shift.EndTime.Microseconds/1000000%3600)/60),       // Minute
					int(shift.EndTime.Microseconds/1000000%60), 0, time.UTC) // Second
			}

			shifts = append(shifts, ShiftInfo{
				ID:         shift.ID,
				DoctorID:   shift.DoctorID,
				StartTime:  startTime,
				EndTime:    endTime,
				Date:       shift.Date.Time.Format("2006-01-02"),
				CreatedAt:  shift.CreatedAt.Time,
				DoctorName: doctor.Name,
			})
		}

		// Cache for 30 minutes
		err = c.SetWithBackground(shiftsKey, &shifts, time.Minute*30)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", shiftsKey, err)
		}
	}

	return shifts, nil
}

// Get all shifts with cache
func (c *ClientType) AllShiftsLoadCache() ([]ShiftInfo, error) {
	shiftsKey := fmt.Sprintf("%s:all-shifts", DOCTOR_INFO_KEY)

	var shifts []ShiftInfo
	err := c.GetWithBackground(shiftsKey, &shifts)
	if err != nil {
		// Cache miss, get from database
		dbShifts, err := db.StoreDB.GetShifts(ctxRedis)
		if err != nil {
			return nil, err
		}

		// Transform to our cache format
		for _, shift := range dbShifts {
			// Try to get doctor name for the shift
			doctor, err := db.StoreDB.GetDoctor(ctxRedis, shift.DoctorID)
			if err != nil {
				// Continue even if we can't get doctor name
				// Since GetShiftsRow doesn't include time info, we need to get time slots
				timeSlots, err := db.StoreDB.GetTimeSlotsByShiftID(ctxRedis, shift.ID)
				var startTime, endTime time.Time

				// Use default times if we can't get time slots
				if err == nil && len(timeSlots) > 0 {
					// Use the first and last time slot
					firstSlot := timeSlots[0]
					lastSlot := timeSlots[len(timeSlots)-1]

					// Convert from microseconds to datetime
					startTime = time.Date(2000, 1, 1,
						int(firstSlot.StartTime.Microseconds/1000000/3600),
						int((firstSlot.StartTime.Microseconds/1000000%3600)/60),
						int(firstSlot.StartTime.Microseconds/1000000%60), 0, time.UTC)

					endTime = time.Date(2000, 1, 1,
						int(lastSlot.EndTime.Microseconds/1000000/3600),
						int((lastSlot.EndTime.Microseconds/1000000%3600)/60),
						int(lastSlot.EndTime.Microseconds/1000000%60), 0, time.UTC)
				}

				shifts = append(shifts, ShiftInfo{
					ID:        shift.ID,
					DoctorID:  shift.DoctorID,
					StartTime: startTime,
					EndTime:   endTime,
					Date:      shift.Date.Time.Format("2006-01-02"),
					CreatedAt: shift.CreatedAt.Time,
				})
				continue
			}

			// Get time slots for this shift
			timeSlots, err := db.StoreDB.GetTimeSlotsByShiftID(ctxRedis, shift.ID)
			var startTime, endTime time.Time

			// Use default times if we can't get time slots
			if err == nil && len(timeSlots) > 0 {
				// Use the first and last time slot
				firstSlot := timeSlots[0]
				lastSlot := timeSlots[len(timeSlots)-1]

				// Convert from microseconds to datetime
				startTime = time.Date(2000, 1, 1,
					int(firstSlot.StartTime.Microseconds/1000000/3600),
					int((firstSlot.StartTime.Microseconds/1000000%3600)/60),
					int(firstSlot.StartTime.Microseconds/1000000%60), 0, time.UTC)

				endTime = time.Date(2000, 1, 1,
					int(lastSlot.EndTime.Microseconds/1000000/3600),
					int((lastSlot.EndTime.Microseconds/1000000%3600)/60),
					int(lastSlot.EndTime.Microseconds/1000000%60), 0, time.UTC)
			}

			shifts = append(shifts, ShiftInfo{
				ID:         shift.ID,
				DoctorID:   shift.DoctorID,
				StartTime:  startTime,
				EndTime:    endTime,
				Date:       shift.Date.Time.Format("2006-01-02"),
				CreatedAt:  shift.CreatedAt.Time,
				DoctorName: doctor.Name,
			})
		}

		// Cache for 30 minutes
		err = c.SetWithBackground(shiftsKey, &shifts, time.Minute*30)
		if err != nil {
			log.Printf("Error when set cache for key %s: %v", shiftsKey, err)
		}
	}

	return shifts, nil
}

// Remove a doctor from cache
func (client *ClientType) RemoveDoctorCache(doctorID int64) {
	doctorKey := fmt.Sprintf("%s:%d", DOCTOR_INFO_KEY, doctorID)
	client.RemoveCacheByKey(doctorKey)
}

// Remove a doctor by username from cache
func (client *ClientType) RemoveDoctorByUsernameCache(username string) {
	usernameKey := fmt.Sprintf("%s:username:%s", DOCTOR_INFO_KEY, username)
	client.RemoveCacheByKey(usernameKey)
}

// Clear all doctor shift caches
func (client *ClientType) ClearDoctorShiftsCache(doctorID int64) {
	shiftsKey := fmt.Sprintf("%s:%d:shifts", DOCTOR_INFO_KEY, doctorID)
	client.RemoveCacheByKey(shiftsKey)

	// Also clear the all shifts cache
	allShiftsKey := fmt.Sprintf("%s:all-shifts", DOCTOR_INFO_KEY)
	client.RemoveCacheByKey(allShiftsKey)
}

// Clear all doctor caches
func (client *ClientType) ClearAllDoctorCache() {
	iter := client.RedisClient.Scan(ctxRedis, 0, fmt.Sprintf("%s*", DOCTOR_INFO_KEY), 0).Iterator()
	for iter.Next(ctxRedis) {
		er := client.RemoveCacheByKey(iter.Val())
		if er != nil {
			log.Printf("Error when remove cache for key %s: %v", iter.Val(), er)
		}
	}
}

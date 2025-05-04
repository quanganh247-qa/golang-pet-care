package pet

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetServiceInterface interface {
	CreatePet(ctx *gin.Context, username string, req createPetRequest) (*CreatePetResponse, error)
	GetPetByID(ctx *gin.Context, petid int64) (*CreatePetResponse, error)
	ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) (*util.PaginationResponse[CreatePetResponse], error)
	UpdatePet(ctx *gin.Context, petid int64, req updatePetRequest) error
	DeletePet(ctx context.Context, petid int64) error
	ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]CreatePetResponse, error)
	SetPetInactive(ctx context.Context, petid int64) error
	GetPetLogsByPetIDService(ctx *gin.Context, pet_id int64, pagination *util.Pagination) ([]PetLogWithPetInfo, error)
	InsertPetLogService(ctx context.Context, req PetLogWithPetInfo) error
	DeletePetLogService(ctx context.Context, logID int64) error
	UpdatePetLogService(ctx context.Context, req UpdatePetLogRequeststruct, log_id int64) error
	UpdatePetAvatar(ctx *gin.Context, petid int64, req updatePetAvatarRequest) error
	GetAllPetLogsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) (*util.PaginationResponse[PetLogWithPetInfo], error)
	GetPetOwnerByPetID(ctx *gin.Context, petID int64) (*user.UserResponse, error)
	GetDetailsPetLogService(ctx context.Context, log_id int64) (*PetLogWithPetInfo, error)
}

func (s *PetService) CreatePet(ctx *gin.Context, username string, req createPetRequest) (*CreatePetResponse, error) {
	var pet CreatePetResponse

	bod, err := time.Parse("2006-01-02", req.BOD)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BOD: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.CreatePet(ctx, db.CreatePetParams{
			Username:        username,
			Name:            req.Name,
			Type:            req.Type,
			Breed:           pgtype.Text{String: req.Breed, Valid: true},
			Age:             pgtype.Int4{Int32: int32(req.Age), Valid: true},
			Weight:          pgtype.Float8{Float64: req.Weight, Valid: true},
			Gender:          pgtype.Text{String: req.Gender, Valid: true},
			Healthnotes:     pgtype.Text{String: req.Healthnotes, Valid: true},
			DataImage:       req.DataImage,
			OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
			BirthDate:       pgtype.Date{Time: bod, Valid: true},
			MicrochipNumber: pgtype.Text{String: req.MicrochipNumber, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create pet: %w", err)
		}

		if _, err := q.CreateMedicalRecord(ctx, pgtype.Int8{Int64: res.Petid, Valid: true}); err != nil {
			return fmt.Errorf("failed to create medical record: %w", err)
		}

		pet = CreatePetResponse{
			Petid:           res.Petid,
			Username:        res.Username,
			Name:            res.Name,
			Type:            res.Type,
			Breed:           res.Breed.String,
			DataImage:       pet.DataImage,
			OriginalImage:   res.OriginalImage.String,
			MicrochipNumber: res.MicrochipNumber.String,
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &pet, nil
}

func (s *PetService) GetPetByID(ctx *gin.Context, petid int64) (*CreatePetResponse, error) {

	// s.redis.ClearPetInfoCache()

	// Try to get from cache
	pet, err := s.redis.PetInfoLoadCache(petid)
	if err != nil {
		// Log cache miss to the context for logging middleware
		ctx.Set("cache_status", "MISS")
		ctx.Set("cache_source", "db")
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	// Log cache hit to the context for logging middleware
	ctx.Set("cache_status", "HIT")
	ctx.Set("cache_source", "redis")

	return &CreatePetResponse{
		Petid:           pet.Petid,
		Username:        pet.Username,
		Name:            pet.Name,
		Type:            pet.Type,
		Breed:           pet.Breed,
		Gender:          pet.Gender,
		Healthnotes:     pet.Healthnotes,
		Age:             pet.Age,
		BOD:             pet.BOD,
		Weight:          pet.Weight,
		DataImage:       pet.DataImage,
		OriginalImage:   pet.OriginalImage,
		MicrochipNumber: pet.MicrochipNumber,
	}, nil
}

func (s *PetService) ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) (*util.PaginationResponse[CreatePetResponse], error) {
	var pets []CreatePetResponse
	offset := (pagination.Page - 1) * pagination.PageSize
	var totalCount int64

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Get total count first
		count, err := q.CountPets(ctx)
		if err != nil {
			return fmt.Errorf("failed to count pets: %w", err)
		}
		totalCount = count

		listParams := db.ListPetsParams{
			Limit:  int32(pagination.PageSize),
			Offset: int32(offset),
		}

		res, err := q.ListPets(ctx, listParams)
		if err != nil {
			return fmt.Errorf("failed to list pets: %w", err)
		}

		for _, r := range res {
			pets = append(pets, CreatePetResponse{
				Petid:         r.Petid,
				Username:      r.Username,
				Name:          r.Name,
				Type:          r.Type,
				Breed:         r.Breed.String,
				Gender:        r.Gender.String,
				Healthnotes:   r.Healthnotes.String,
				Age:           int16(r.Age.Int32),
				Weight:        r.Weight.Float64,
				DataImage:     r.DataImage,
				OriginalImage: r.OriginalImage.String,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	response := &util.PaginationResponse[CreatePetResponse]{
		Count: totalCount,
		Rows:  &pets,
	}
	return response.Build(), nil
}

func (s *PetService) UpdatePet(ctx *gin.Context, petid int64, req updatePetRequest) error {

	var params db.UpdatePetParams

	pet, err := s.storeDB.GetPetByID(ctx, petid)
	if err != nil {
		return fmt.Errorf("failed to get pet: %w", err)
	}

	if req.BOD != "" {
		bod, err := time.Parse("2006-01-02", req.BOD)
		if err != nil {
			return fmt.Errorf("failed to parse BOD: %w", err)
		}
		params.BirthDate = pgtype.Date{Time: bod, Valid: true}

		age := time.Since(bod).Hours() / 24 / 365
		params.Age = pgtype.Int4{Int32: int32(age), Valid: true}
	} else {
		params.BirthDate = pet.BirthDate
		age := time.Since(pet.BirthDate.Time).Hours() / 24 / 365
		params.Age = pgtype.Int4{Int32: int32(age), Valid: true}
	}
	if req.Name != "" {
		params.Name = req.Name
	} else {
		params.Name = pet.Name
	}
	if req.Gender != "" {
		params.Gender = pgtype.Text{String: req.Gender, Valid: true}
	} else {
		params.Name = pet.Name
	}

	if req.Type != "" {
		params.Type = req.Type
	} else {
		params.Type = pet.Type
	}

	if req.Breed != "" {
		params.Breed = pgtype.Text{String: req.Breed, Valid: true}
	} else {
		params.Breed = pet.Breed
	}

	// if req.Age != 0 {
	// 	params.Age = pgtype.Int4{Int32: int32(req.Age), Valid: true}
	// } else {
	// 	params.Age = pet.Age
	// }

	if req.Weight != 0 {
		params.Weight = pgtype.Float8{Float64: req.Weight, Valid: true}
	} else {
		params.Weight = pet.Weight
	}

	if req.Healthnotes != "" {
		params.Healthnotes = pgtype.Text{String: req.Healthnotes, Valid: true}
	} else {
		params.Healthnotes = pet.Healthnotes
	}

	params.Petid = petid

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdatePet(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to update pet: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	// Clear the pet info cache
	if s.redis != nil {
		s.redis.RemovePetInfoCache(petid)
		// Also clear the user's pet list cache
		s.redis.ClearUserPetsCache(pet.Username)
	}

	return nil
}

func (s *PetService) UpdatePetAvatar(ctx *gin.Context, petid int64, req updatePetAvatarRequest) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdatePetAvatar(ctx, db.UpdatePetAvatarParams{
			Petid:     petid,
			DataImage: req.DataImage,
		})
		if err != nil {
			return fmt.Errorf("failed to update pet avatar: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	go s.redis.RemovePetInfoCache(petid)
	return nil
}

func (s *PetService) DeletePet(ctx context.Context, petid int64) error {
	return s.storeDB.DeletePet(ctx, petid)
}

func (s *PetService) ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]CreatePetResponse, error) {
	var pets []CreatePetResponse

	// If we're requesting the first page with a reasonable page size, try cache first
	if pagination.Page == 1 && pagination.PageSize <= 100 {
		cachedPets, err := s.redis.GetPetsByUsernameCache(username)
		if err == nil {
			// We have the pets cached, convert them to our response format
			for _, pet := range cachedPets {
				pets = append(pets, CreatePetResponse{
					Petid:           pet.Petid,
					Username:        pet.Username,
					Name:            pet.Name,
					Type:            pet.Type,
					Breed:           pet.Breed,
					Gender:          pet.Gender,
					Healthnotes:     pet.Healthnotes,
					Age:             pet.Age,
					BOD:             pet.BOD,
					Weight:          pet.Weight,
					DataImage:       pet.DataImage,
					OriginalImage:   pet.OriginalImage,
					MicrochipNumber: pet.MicrochipNumber,
				})
			}
			// Log cache hit
			ctx.Set("cache_status", "HIT")
			ctx.Set("cache_source", "redis:user_pets")
			return pets, nil
		}
	}

	// Cache miss or pagination beyond what's cached, query DB
	// Log cache miss
	ctx.Set("cache_status", "MISS")
	ctx.Set("cache_source", "db")

	offset := (pagination.Page - 1) * pagination.PageSize

	listParams := db.ListPetsByUsernameParams{
		Username: username,
		Limit:    int32(pagination.PageSize),
		Offset:   int32(offset),
	}

	res, err := s.storeDB.ListPetsByUsername(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list pets for username %s: %w", username, err)
	}

	for _, r := range res {
		pets = append(pets, CreatePetResponse{
			Petid:           r.Petid,
			Username:        r.Username,
			Name:            r.Name,
			Type:            r.Type,
			Breed:           r.Breed.String,
			Gender:          r.Gender.String,
			Healthnotes:     r.Healthnotes.String,
			Age:             int16(r.Age.Int32),
			Weight:          r.Weight.Float64,
			DataImage:       r.DataImage,
			OriginalImage:   r.OriginalImage.String,
			MicrochipNumber: r.MicrochipNumber.String,
		})
	}

	return pets, nil
}

func (s *PetService) SetPetInactive(ctx context.Context, petid int64) error {
	return s.storeDB.SetPetInactive(ctx, petid)
}

func (s *PetService) GetPetLogsByPetIDService(ctx *gin.Context, pet_id int64, pagination *util.Pagination) ([]PetLogWithPetInfo, error) {
	var pets []PetLogWithPetInfo

	// If first page and reasonable page size, try to use cache
	if pagination.Page == 1 && pagination.PageSize <= 20 && s.redis != nil {
		cachedLogs, err := s.redis.PetLogSummaryByPetIDLoadCache(pet_id, int32(pagination.PageSize))
		if err == nil {
			// Cache hit, convert to response format
			for _, log := range cachedLogs {
				pets = append(pets, PetLogWithPetInfo{
					PetID:    log.PetID,
					LogID:    log.LogID,
					DateTime: log.Datetime.Format(time.RFC3339),
					Title:    log.Title,
					Notes:    log.Notes,
					PetName:  log.PetName,
					PetType:  log.PetType,
					PetBreed: log.PetBreed,
				})
			}
			return pets, nil
		}
	}

	// Cache miss or pagination beyond what's cached
	offset := (pagination.Page - 1) * pagination.PageSize

	listParams := db.GetPetLogsByPetIDParams{
		Petid:  pet_id,
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	}

	res, err := s.storeDB.GetPetLogsByPetID(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list pets  %s: %w", pet_id, err)
	}

	for _, r := range res {
		pets = append(pets, PetLogWithPetInfo{
			PetID:    r.Petid,
			LogID:    r.LogID,
			DateTime: r.Datetime.Time.Format(time.RFC3339),
			Title:    r.Title.String,
			Notes:    r.Notes.String,
			PetName:  r.Name.String,
			PetType:  r.Type.String,
			PetBreed: r.Breed.String,
		})
	}

	return pets, nil
}

// Add log for pet
func (s *PetService) InsertPetLogService(ctx context.Context, req PetLogWithPetInfo) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		pet, err := q.GetPetByID(ctx, req.PetID)
		if err != nil {
			return fmt.Errorf("failed to get pet: %w", err)
		}

		dateTime, err := time.Parse(time.RFC3339, req.DateTime)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}

		var insertParams = db.CreatePetLogParams{
			Petid:    pet.Petid,
			Title:    pgtype.Text{String: req.Title, Valid: true},
			Notes:    pgtype.Text{String: req.Notes, Valid: true},
			Datetime: pgtype.Timestamp{Time: dateTime, Valid: true},
		}
		_, err = q.CreatePetLog(ctx, insertParams)
		if err != nil {
			return fmt.Errorf("failed to insert pet log: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Invalidate cache for the pet logs
	s.redis.ClearPetLogsByPetCache(req.PetID)
	middleware.InvalidateCache("pet_logs")

	return nil
}

// DeletePetLogService delete log for pet
func (s *PetService) DeletePetLogService(ctx context.Context, logID int64) error {
	// Get the log first to know which pet it belongs to
	var petID int64
	if s.redis != nil {

		if log, err := s.storeDB.GetPetLogByID(ctx, logID); err == nil {
			petID = log.Petid
		}
	}

	// Delete the pet log
	err := s.storeDB.DeletePetLog(ctx, logID)
	if err != nil {
		return fmt.Errorf("failed to delete pet log: %w", err)
	}

	// Invalidate cache
	middleware.InvalidateCache("pet_logs")

	// Also clear the specific caches if Redis is available
	if s.redis != nil && petID != 0 {
		s.redis.ClearPetLogsByPetCache(petID)
		s.redis.ClearPetLogSummaryCache(petID)
	}

	return nil
}

// UpdatePetLogService update log for pet
func (s *PetService) UpdatePetLogService(ctx context.Context, req UpdatePetLogRequeststruct, log_id int64) error {

	log, err := s.storeDB.GetPetLogByID(ctx, log_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("pet log not found: %d", log_id)
		}
		return fmt.Errorf("failed to get pet log: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		// Check if the pet exists
		_, err := q.GetPetByID(ctx, log.Petid)
		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("pet not found: %d", log.Petid)
			}
			return fmt.Errorf("failed to get pet: %w", err)
		}

		// Finally, update the log
		updateParams := db.UpdatePetLogParams{
			LogID: log_id,
		}

		if req.Title != "" {
			updateParams.Title = pgtype.Text{String: req.Title, Valid: true}
		} else {
			updateParams.Title = log.Title
		}
		if req.Notes != "" {
			updateParams.Notes = pgtype.Text{String: req.Notes, Valid: true}
		} else {
			updateParams.Notes = log.Notes
		}
		if req.DateTime != "" {
			dateTime, err := time.Parse(time.RFC3339, req.DateTime)
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
			updateParams.Datetime = pgtype.Timestamp{Time: dateTime, Valid: true}
		} else {
			updateParams.Datetime = log.Datetime

		}

		err = q.UpdatePetLog(ctx, updateParams)
		if err != nil {
			return fmt.Errorf("failed to update pet log: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	// Invalidate cache for the pet logs
	s.redis.ClearPetLogsByPetCache(log.Petid)
	middleware.InvalidateCache("pet_logs")

	return nil
}

// GetDetailsPetLogService get log for pet
func (s *PetService) GetDetailsPetLogService(ctx context.Context, log_id int64) (*PetLogWithPetInfo, error) {
	log, err := s.storeDB.GetPetLogByID(ctx, log_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet log: %w", err)
	}

	pet, err := s.storeDB.GetPetByID(ctx, log.Petid)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	return &PetLogWithPetInfo{
		PetID:    pet.Petid,
		PetName:  pet.Name,
		PetType:  pet.Type,
		PetBreed: pet.Breed.String,
		LogID:    log.LogID,
		DateTime: log.Datetime.Time.Format(time.RFC3339),
		Title:    log.Title.String,
		Notes:    log.Notes.String,
	}, nil
}

// Add this method to your PetService implementation
func (s *PetService) GetAllPetLogsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) (*util.PaginationResponse[PetLogWithPetInfo], error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	// Get total count
	count, err := s.storeDB.CountAllPetLogsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to count pet logs: %w", err)
	}

	// Get logs with pagination
	logs, err := s.storeDB.GetAllPetLogsByUsername(ctx, db.GetAllPetLogsByUsernameParams{
		Username: username,
		Limit:    int32(pagination.PageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pet logs: %w", err)
	}

	// Transform to response model
	var petLogs []PetLogWithPetInfo
	for _, log := range logs {
		petLogs = append(petLogs, PetLogWithPetInfo{
			PetID:    log.Petid,
			PetName:  log.PetName,
			PetType:  log.PetType,
			PetBreed: log.PetBreed.String,
			LogID:    log.LogID,
			DateTime: log.Datetime.Time.Format(time.RFC3339),
			Title:    log.Title.String,
			Notes:    log.Notes.String,
		})
	}

	return &util.PaginationResponse[PetLogWithPetInfo]{
		Count: count,
		Rows:  &petLogs,
	}, nil
}

// Add this new method to your PetService
func (s *PetService) GetPetOwnerByPetID(ctx *gin.Context, petID int64) (*user.UserResponse, error) {
	// First get the pet to retrieve owner username
	pet, err := s.storeDB.GetPetByID(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	// Then get the owner details using the username from pet record
	owner, err := s.storeDB.GetUser(ctx, pet.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	return &user.UserResponse{
		Username:    owner.Username,
		Email:       owner.Email,
		Role:        owner.Role.String,
		PhoneNumber: owner.PhoneNumber.String,
		FullName:    owner.FullName,
		Address:     owner.Address.String,
	}, nil
}

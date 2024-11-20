package pet

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"strings"
=======
>>>>>>> 67140c6 (updated create pet)
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetServiceInterface interface {
<<<<<<< HEAD
	CreatePet(ctx *gin.Context, username string, req createPetRequest) (*CreatePetResponse, error)
	GetPetByID(ctx *gin.Context, petid int64) (*CreatePetResponse, error)
	ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) ([]CreatePetResponse, error)
	UpdatePet(ctx *gin.Context, petid int64, req updatePetRequest) error
	DeletePet(ctx context.Context, petid int64) error
	ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]CreatePetResponse, error)
=======
	CreatePet(ctx *gin.Context, username string, req createPetRequest) (*createPetResponse, error)
	GetPetByID(ctx *gin.Context, petid int64) (*createPetResponse, error)
	ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) ([]createPetResponse, error)
	UpdatePet(ctx *gin.Context, petid int64, req createPetRequest) error
	DeletePet(ctx context.Context, petid int64) error
	ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]createPetResponse, error)
>>>>>>> c73e2dc (pagination function)
	SetPetInactive(ctx context.Context, petid int64) error
	GetPetLogsByPetIDService(ctx *gin.Context, pet_id int64, pagination *util.Pagination) ([]PetLog, error)
	InsertPetLogService(ctx context.Context, req PetLog) error
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	DeletePetLogService(ctx context.Context, logID int64) error
	UpdatePetLogService(ctx context.Context, req PetLog, log_id int64) error
	UpdatePetAvatar(ctx *gin.Context, petid int64, req updatePetAvatarRequest) error
=======
>>>>>>> 7e616af (add pet log schema)
=======
	DeletePetLogService(ctx context.Context, petID int64, logID int64) error
=======
	DeletePetLogService(ctx context.Context, logID int64) error
>>>>>>> 884b92e (update pet logs api)
	UpdatePetLogService(ctx context.Context, req PetLog, log_id int64) error
>>>>>>> 3835eb4 (update pet_schedule api)
}

<<<<<<< HEAD
func (s *PetService) CreatePet(ctx *gin.Context, username string, req createPetRequest) (*CreatePetResponse, error) {
	var pet CreatePetResponse

	bod, err := time.Parse("2006-01-02", req.BOD)
=======
func (s *PetService) CreatePet(ctx *gin.Context, username string, req createPetRequest) (*createPetResponse, error) {
	var pet createPetResponse
<<<<<<< HEAD
	bod, _, err := util.ParseStringToTime(req.BOD, "")
>>>>>>> 9d28896 (image pet)
=======

	bod, err := time.Parse("2006-01-02", req.BOD)
>>>>>>> 67140c6 (updated create pet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BOD: %w", err)
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.CreatePet(ctx, db.CreatePetParams{
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9d28896 (image pet)
			Username:        username,
			Name:            req.Name,
			Type:            req.Type,
			Breed:           pgtype.Text{String: req.Breed, Valid: true},
			Age:             pgtype.Int4{Int32: int32(req.Age), Valid: true},
			Weight:          pgtype.Float8{Float64: req.Weight, Valid: true},
			Gender:          pgtype.Text{String: req.Gender, Valid: true},
			Healthnotes:     pgtype.Text{String: req.Healthnotes, Valid: true},
			DataImage:       req.DataImage,
<<<<<<< HEAD
<<<<<<< HEAD
			OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
			BirthDate:       pgtype.Date{Time: bod, Valid: true},
			MicrochipNumber: pgtype.Text{String: req.MicrochipNumber, Valid: true},
<<<<<<< HEAD
=======
			Username:    username,
			Name:        req.Name,
			Type:        req.Type,
			Breed:       pgtype.Text{String: req.Breed, Valid: true},
			Age:         pgtype.Int4{Int32: int32(req.Age), Valid: true},
			Weight:      pgtype.Float8{Float64: req.Weight, Valid: true},
			Gender:      pgtype.Text{String: req.Gender, Valid: true},
			Healthnotes: pgtype.Text{String: req.Healthnotes, Valid: true},
>>>>>>> 0fb3f30 (user images)
=======
			OriginalImage:   req.OriginalImage,
=======
			OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
>>>>>>> 272832d (redis cache)
			BirthDate:       pgtype.Date{Time: bod, Valid: true},
			MicrochipNumber: pgtype.Text{String: req.MicrophoneNumber, Valid: true},
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 67140c6 (updated create pet)
		})
		if err != nil {
			return fmt.Errorf("failed to create pet: %w", err)
		}
<<<<<<< HEAD

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
=======
		pet = createPetResponse{
<<<<<<< HEAD
			Petid:     res.Petid,
			Username:  res.Username,
			Name:      res.Name,
			Type:      res.Type,
			Breed:     res.Breed.String,
			Age:       int16(res.Age.Int32),
			Weight:    res.Weight.Float64,
			DataImage: pet.DataImage,
>>>>>>> 9d28896 (image pet)
=======
			Petid:           res.Petid,
			Username:        res.Username,
			Name:            res.Name,
			Type:            res.Type,
			Breed:           res.Breed.String,
			DataImage:       pet.DataImage,
			OriginalImage:   res.OriginalImage.String,
			MicrochipNumber: res.MicrochipNumber.String,
>>>>>>> 67140c6 (updated create pet)
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return &pet, nil
}

<<<<<<< HEAD
func (s *PetService) GetPetByID(ctx *gin.Context, petid int64) (*CreatePetResponse, error) {

	pet, err := s.redis.PetInfoLoadCache(petid)
=======
func (s *PetService) GetPetByID(ctx *gin.Context, petid int64) (*createPetResponse, error) {
	var pet createPetResponse
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.GetPetByID(ctx, petid)
		if err != nil {
			return fmt.Errorf("failed to get pet: %w", err)
		}
		pet = createPetResponse{
			Petid:         res.Petid,
			Username:      res.Username,
			Name:          res.Name,
			Type:          res.Type,
			Breed:         res.Breed.String,
			BOD:           res.BirthDate.Time.Format("2006-01-02"),
			Age:           int16(res.Age.Int32),
			Weight:        res.Weight.Float64,
			DataImage:     res.DataImage,
			OriginalImage: res.OriginalImage.String,
		}
		return nil
	})
>>>>>>> 9d28896 (image pet)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}
	return &CreatePetResponse{
		Petid:           pet.Petid,
		Username:        pet.Username,
		Name:            pet.Name,
		Type:            pet.Type,
		Breed:           pet.Breed,
		Age:             pet.Age,
		BOD:             pet.BOD,
		Weight:          pet.Weight,
		DataImage:       pet.DataImage,
		OriginalImage:   pet.OriginalImage,
		MicrochipNumber: pet.MicrochipNumber,
	}, nil
}

<<<<<<< HEAD
func (s *PetService) ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) ([]CreatePetResponse, error) {
	var pets []CreatePetResponse
=======
func (s *PetService) ListPets(ctx *gin.Context, req listPetsRequest, pagination *util.Pagination) ([]createPetResponse, error) {
	var pets []createPetResponse
>>>>>>> c73e2dc (pagination function)
	offset := (pagination.Page - 1) * pagination.PageSize

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		listParams := db.ListPetsParams{
			Limit:  int32(pagination.PageSize),
			Offset: int32(offset),
		}

		res, err := q.ListPets(ctx, listParams)
		if err != nil {
			return fmt.Errorf("failed to list pets: %w", err)
		}

		for _, r := range res {
<<<<<<< HEAD
			pets = append(pets, CreatePetResponse{
=======
			pets = append(pets, createPetResponse{
>>>>>>> 9d28896 (image pet)
				Petid:         r.Petid,
				Username:      r.Username,
				Name:          r.Name,
				Type:          r.Type,
				Breed:         r.Breed.String,
				Age:           int16(r.Age.Int32),
				Weight:        r.Weight.Float64,
				DataImage:     r.DataImage,
<<<<<<< HEAD
<<<<<<< HEAD
				OriginalImage: r.OriginalImage.String,
=======
				OriginalImage: r.OriginalImage,
>>>>>>> 9d28896 (image pet)
=======
				OriginalImage: r.OriginalImage.String,
>>>>>>> 272832d (redis cache)
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return pets, nil
}

<<<<<<< HEAD
func (s *PetService) UpdatePet(ctx *gin.Context, petid int64, req updatePetRequest) error {

	var params db.UpdatePetParams

	pet, err := s.storeDB.GetPetByID(ctx, petid)
	if err != nil {
		return fmt.Errorf("failed to get pet: %w", err)
=======
func (s *PetService) UpdatePet(ctx *gin.Context, petid int64, req createPetRequest) error {
	params := db.UpdatePetParams{
		Petid:       petid,
		Name:        req.Name,
		Type:        req.Type,
		Breed:       pgtype.Text{String: req.Breed, Valid: true},
		Age:         pgtype.Int4{Int32: int32(req.Age), Valid: true},
		Weight:      pgtype.Float8{Float64: req.Weight, Valid: true},
		Gender:      pgtype.Text{String: req.Gender, Valid: true},
		Healthnotes: pgtype.Text{String: req.Healthnotes, Valid: true},
>>>>>>> 0fb3f30 (user images)
	}

	if req.BOD != "" {
		bod, err := time.Parse("2006-01-02", req.BOD)
		if err != nil {
			return fmt.Errorf("failed to parse BOD: %w", err)
		}
		params.BirthDate = pgtype.Date{Time: bod, Valid: true}
	} else {
		params.BirthDate = pet.BirthDate

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

	if req.Age != 0 {
		params.Age = pgtype.Int4{Int32: int32(req.Age), Valid: true}
	} else {
		params.Age = pet.Age
	}

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
	go s.redis.RemovePetInfoCache(petid)
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

<<<<<<< HEAD
func (s *PetService) ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]CreatePetResponse, error) {
	var pets []CreatePetResponse
=======
func (s *PetService) ListPetsByUsername(ctx *gin.Context, username string, pagination *util.Pagination) ([]createPetResponse, error) {
	var pets []createPetResponse
>>>>>>> c73e2dc (pagination function)
	offset := (pagination.Page - 1) * pagination.PageSize

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		listParams := db.ListPetsByUsernameParams{
			Username: username,
			Limit:    int32(pagination.PageSize),
			Offset:   int32(offset),
		}

		res, err := q.ListPetsByUsername(ctx, listParams)
		if err != nil {
			return fmt.Errorf("failed to list pets for user %s: %w", username, err)
		}

		for _, r := range res {
<<<<<<< HEAD
			pets = append(pets, CreatePetResponse{
=======
			pets = append(pets, createPetResponse{
>>>>>>> 9d28896 (image pet)
				Petid:         r.Petid,
				Username:      r.Username,
				Name:          r.Name,
				Type:          r.Type,
				Breed:         r.Breed.String,
				Age:           int16(r.Age.Int32),
				Weight:        r.Weight.Float64,
				DataImage:     r.DataImage,
<<<<<<< HEAD
<<<<<<< HEAD
				OriginalImage: r.OriginalImage.String,
=======
				OriginalImage: r.OriginalImage,
>>>>>>> 9d28896 (image pet)
=======
				OriginalImage: r.OriginalImage.String,
>>>>>>> 272832d (redis cache)
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return pets, nil
}

func (s *PetService) SetPetInactive(ctx context.Context, petid int64) error {
	return s.storeDB.SetPetInactive(ctx, petid)
}

func (s *PetService) GetPetLogsByPetIDService(ctx *gin.Context, pet_id int64, pagination *util.Pagination) ([]PetLog, error) {
	var pets []PetLog
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
		pets = append(pets, PetLog{
			PetID:    r.Petid,
			LogID:    r.LogID,
			DateTime: r.Datetime.Time.Format(time.RFC3339),
			Title:    r.Title.String,
			Notes:    r.Notes.String,
		})
	}

	return pets, nil
}

// Add log for pet
func (s *PetService) InsertPetLogService(ctx context.Context, req PetLog) error {

	log := db.CreatePetLogParams{
		Petid: req.PetID,
		Title: pgtype.Text{String: req.Title, Valid: true},
		Notes: pgtype.Text{String: req.Notes, Valid: true},
	}
	if req.DateTime == "" {
		log.Datetime = pgtype.Timestamp{Time: time.Now(), Valid: true}
	} else {
		datetime, err := time.Parse(time.RFC3339, req.DateTime)
		if err != nil {
			return fmt.Errorf("failed to parse DateTime: %w", err)
		}
		log.Datetime = pgtype.Timestamp{Time: datetime, Valid: true}
	}

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.CreatePetLog(ctx, log)
		if err != nil {
			return fmt.Errorf("failed to insert pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// DeletePetLogService delete log for pet
func (s *PetService) DeletePetLogService(ctx context.Context, logID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeletePetLog(ctx, logID)
		if err != nil {
			return fmt.Errorf("failed to delete pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction delete log failed: %w", err)
	}
	return nil
}

// UpdatePetLogService update log for pet
func (s *PetService) UpdatePetLogService(ctx context.Context, req PetLog, log_id int64) error {

	pet_log, err := s.storeDB.GetPetLogByID(ctx, db.GetPetLogByIDParams{
		Petid: req.PetID,
		LogID: log_id,
	})
	if err != nil {
		return fmt.Errorf("failed to get pet log: %w", err)
	}
	// check input
	if req.DateTime == "" {
		req.DateTime = pet_log.Datetime.Time.Format(time.RFC3339)
	}
	if req.Title == "" {
		req.Title = pet_log.Title.String
	}
	if req.Notes == "" {
		req.Notes = pet_log.Notes.String
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdatePetLog(ctx, db.UpdatePetLogParams{
			LogID: log_id,
			Title: pgtype.Text{String: req.Title, Valid: true},
			Notes: pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction update log failed: %w", err)
	}
	return nil
}

func formatTreatments(treatments []db.GetTreatmentsByPetRow) string {
	var result strings.Builder
	for _, t := range treatments {
		result.WriteString(fmt.Sprintf("- Condition: %s\n  Status: %s\n  Period: %s to %s\n",
			t.Disease,
			t.Status.String,
			t.StartDate.Time.Format("2006-01-02"),
			t.EndDate.Time.Format("2006-01-02"),
		))
	}
	return result.String()
}

func formatVaccinations(vaccinations []db.Vaccination) string {
	var result strings.Builder
	for _, v := range vaccinations {
		result.WriteString(fmt.Sprintf("- Vaccine: %s\n  Administered: %s\n  Next Due: %s\n  Provider: %s\n",
			v.Vaccinename,
			v.Dateadministered.Time.Format("2006-01-02"),
			v.Nextduedate.Time.Format("2006-01-02"),
			v.Vaccineprovider.String))
	}
	return result.String()
}

// Log for pet

func (s *PetService) GetPetLogsByPetIDService(ctx *gin.Context, pet_id int64, pagination *util.Pagination) ([]PetLog, error) {
	var pets []PetLog
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
		pets = append(pets, PetLog{
			PetID:    r.Petid,
			LogID:    r.LogID,
			DateTime: r.Datetime.Time.Format(time.RFC3339),
			Title:    r.Title.String,
			Notes:    r.Notes.String,
		})
	}

	return pets, nil
}

// Add log for pet
func (s *PetService) InsertPetLogService(ctx context.Context, req PetLog) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.InsertPetLog(ctx, db.InsertPetLogParams{
			Petid:    req.PetID,
			Datetime: pgtype.Timestamp{Time: time.Now(), Valid: true},
			Title:    pgtype.Text{String: req.Title, Valid: true},
			Notes:    pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insert pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// DeletePetLogService delete log for pet
func (s *PetService) DeletePetLogService(ctx context.Context, logID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeletePetLog(ctx, logID)
		if err != nil {
			return fmt.Errorf("failed to delete pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction delete log failed: %w", err)
	}
	return nil
}

// UpdatePetLogService update log for pet
func (s *PetService) UpdatePetLogService(ctx context.Context, req PetLog, log_id int64) error {

	pet_log, err := s.storeDB.GetPetLogByID(ctx, db.GetPetLogByIDParams{
		Petid: req.PetID,
		LogID: log_id,
	})
	if err != nil {
		return fmt.Errorf("failed to get pet log: %w", err)
	}
	// check input
	if req.DateTime == "" {
		req.DateTime = pet_log.Datetime.Time.Format(time.RFC3339)
	}
	if req.Title == "" {
		req.Title = pet_log.Title.String
	}
	if req.Notes == "" {
		req.Notes = pet_log.Notes.String
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdatePetLog(ctx, db.UpdatePetLogParams{
			LogID: log_id,
			Title: pgtype.Text{String: req.Title, Valid: true},
			Notes: pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update pet log: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction update log failed: %w", err)
	}
	return nil
}

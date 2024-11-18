package petschedule

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e01abc5 (pet schedule api)
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
=======
>>>>>>> e01abc5 (pet schedule api)
=======
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
>>>>>>> 6610455 (feat: redis queue)
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleControllerInterface interface {
	createPetSchedule(ctx *gin.Context)
	getAllSchedulesByPet(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
	listPetSchedulesByUsername(ctx *gin.Context)
	activePetSchedule(ctx *gin.Context)
	deletePetSchedule(ctx *gin.Context)
	updatePetScheduleService(ctx *gin.Context)
	generateScheduleSuggestion(ctx *gin.Context)
=======
>>>>>>> e01abc5 (pet schedule api)
=======
	listPetSchedulesByUsername(ctx *gin.Context)
>>>>>>> 6610455 (feat: redis queue)
}

func (c *PetScheduleController) createPetSchedule(ctx *gin.Context) {

<<<<<<< HEAD
<<<<<<< HEAD
	petIDStr := ctx.Param("petid")
=======
	petIDStr := ctx.Query("pet_id")
>>>>>>> e01abc5 (pet schedule api)
=======
	petIDStr := ctx.Param("petid")
>>>>>>> 6610455 (feat: redis queue)
	if petIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missed value pet id"})
		return
	}
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid convert pet id"})
		return
	}

	var req PetScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	err = c.service.CreatePetScheduleService(ctx, req, petID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Pet schedule created successfully"})
}

func (s *PetScheduleController) getAllSchedulesByPet(ctx *gin.Context) {
<<<<<<< HEAD
<<<<<<< HEAD
	petIDStr := ctx.Param("petid")
=======
	petIDStr := ctx.Query("pet_id")
>>>>>>> e01abc5 (pet schedule api)
=======
	petIDStr := ctx.Param("petid")
>>>>>>> 6610455 (feat: redis queue)
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		// Handle the error, e.g., by sending an error response or logging it
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet_id"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	schedules, err := s.service.GetAllSchedulesByPetService(ctx, petID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Schedules", schedules))
<<<<<<< HEAD
}

func (s *PetScheduleController) listPetSchedulesByUsername(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	schedules, err := s.service.ListPetSchedulesByUsernameService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Schedules", schedules))
}

func (s *PetScheduleController) activePetSchedule(ctx *gin.Context) {

	var req ActiceRemider
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	scheduleIDStr := ctx.Param("schedule_id")
	scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id"})
		return
	}

	err = s.service.ActivePetScheduleService(ctx, scheduleID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Active reminder", "Success"))
}

func (s *PetScheduleController) deletePetSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("schedule_id")
	scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id"})
		return
	}

	err = s.service.DeletePetScheduleService(ctx, scheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Delete reminder", "Success"))
}

func (s *PetScheduleController) updatePetScheduleService(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("schedule_id")
	scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id"})
		return
	}

	var req PetScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	err = s.service.UpdatePetScheduleService(ctx, scheduleID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Update reminder", "Success"))

}

func (s *PetScheduleController) generateScheduleSuggestion(ctx *gin.Context) {
	var req ScheduleSuggestion
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	response, err := s.service.ProcessSuggestionGemini(ctx, req.Voice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Schedule suggestion", response))
=======
type PetScheduleControllerInterface interface {
>>>>>>> 272832d (redis cache)
=======
>>>>>>> e01abc5 (pet schedule api)
}

func (s *PetScheduleController) listPetSchedulesByUsername(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	schedules, err := s.service.ListPetSchedulesByUsernameService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorValidator(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Schedules", schedules))
}

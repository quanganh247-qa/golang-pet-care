package reminder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ReminderControllerInterface interface {
	CreateReminderController(ctx *gin.Context)
}

func (controller *ReminderController) CreateReminderController(ctx *gin.Context) {
	var req CreateReminderRequest
	fmt.Println("CreateReminderController", req.Petid, req.Userid, req.Remindertype, req.Reminderdate, req.Description.String, req.Status.String)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := controller.reminderService.CreateReminderService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("created successful", res))
}

//api -> controller -> service -> model -> db

package rooms

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type RoomControllerInterface interface {
	CreateRoom(ctx *gin.Context)
	GetRoomByID(ctx *gin.Context)
	ListRooms(ctx *gin.Context)
	UpdateRoom(ctx *gin.Context)
	DeleteRoom(ctx *gin.Context)
}

func (c *RoomController) CreateRoom(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	res, err := c.service.CreateRoom(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Room created successfully", res))
}

func (c *RoomController) GetRoomByID(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	roomIDInt, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetRoomByID(ctx, roomIDInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Room retrieved successfully", res))
}

func (c *RoomController) ListRooms(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	rooms, err := c.service.ListRooms(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Rooms retrieved successfully", rooms))
}

func (c *RoomController) UpdateRoom(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	roomIDInt, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	err = c.service.UpdateRoom(ctx, authPayload.Username, roomIDInt, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Room updated successfully", nil))
}

func (c *RoomController) DeleteRoom(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	roomIDInt, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	err = c.service.DeleteRoom(ctx, authPayload.Username, roomIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Room deleted successfully", nil))
}

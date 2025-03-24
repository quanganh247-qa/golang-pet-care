package pet

import (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"encoding/json"
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
	"encoding/json"
>>>>>>> 67140c6 (updated create pet)
=======
	"encoding/json"
>>>>>>> 67140c6 (updated create pet)
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> c73e2dc (pagination function)
=======
	"io/ioutil"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 473cd1d (uplaod image method)
=======
	"fmt"
>>>>>>> c73e2dc (pagination function)
=======
	"io/ioutil"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 473cd1d (uplaod image method)
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetControllerInterface interface {
	CreatePet(ctx *gin.Context)
	GetPetByID(ctx *gin.Context)
	ListPets(ctx *gin.Context)
	ListPetsByUsername(ctx *gin.Context)
	UpdatePet(ctx *gin.Context)
	DeletePet(ctx *gin.Context)
	GetPetLogsByPetID(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	InsertPetLog(ctx *gin.Context)
	DeletePetLog(ctx *gin.Context)
	UpdatePetLog(ctx *gin.Context)
	UpdatePetAvatar(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 7e616af (add pet log schema)
=======
	InsertPetLog(ctx *gin.Context)
	DeletePetLog(ctx *gin.Context)
	UpdatePetLog(ctx *gin.Context)
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
	GetPetProfileSummary(ctx *gin.Context)
>>>>>>> ffc9071 (AI suggestion)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
>>>>>>> 7e616af (add pet log schema)
=======
	InsertPetLog(ctx *gin.Context)
	DeletePetLog(ctx *gin.Context)
	UpdatePetLog(ctx *gin.Context)
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
	GetPetProfileSummary(ctx *gin.Context)
>>>>>>> ffc9071 (AI suggestion)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
}

func (c *PetController) CreatePet(ctx *gin.Context) {
	var req createPetRequest

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 67140c6 (updated create pet)
=======
>>>>>>> 67140c6 (updated create pet)
	// Parse the JSON data from the "data" form field
	jsonData := ctx.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
<<<<<<< HEAD
<<<<<<< HEAD

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
=======
	name := ctx.PostForm("name")
	t := ctx.PostForm("type")
	breed := ctx.PostForm("breed")
	age := ctx.PostForm("age")
	weight := ctx.PostForm("weight")
	gender := ctx.PostForm("gender")
	healthnotes := ctx.PostForm("healthnotes")
	microchip := ctx.PostForm("microchip_number")
	bod := ctx.PostForm("birth_date")
=======
>>>>>>> 67140c6 (updated create pet)

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

<<<<<<< HEAD
=======
	name := ctx.PostForm("name")
	t := ctx.PostForm("type")
	breed := ctx.PostForm("breed")
	age := ctx.PostForm("age")
	weight := ctx.PostForm("weight")
	gender := ctx.PostForm("gender")
	healthnotes := ctx.PostForm("healthnotes")
	microchip := ctx.PostForm("microchip_number")
	bod := ctx.PostForm("birth_date")
=======
>>>>>>> 67140c6 (updated create pet)

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

<<<<<<< HEAD
>>>>>>> 9d28896 (image pet)
	// Handle image file
	file, header, err := ctx.Request.FormFile("image")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("image is required")))
		return
	}
	defer file.Close()

	// Read the file content into a byte array
	dataImage, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data image"})
		return
	}
	// get original image
<<<<<<< HEAD
>>>>>>> 9d28896 (image pet)

=======
>>>>>>> 473cd1d (uplaod image method)
=======

>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 473cd1d (uplaod image method)
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

	req.OriginalImage = originalImageName
	req.DataImage = dataImage
=======
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
<<<<<<< HEAD
>>>>>>> c73e2dc (pagination function)

	req.Name = name
	req.Type = t
	req.Breed = breed
	req.Healthnotes = healthnotes
	req.Gender = gender
=======
>>>>>>> 67140c6 (updated create pet)
	req.OriginalImage = header.Filename
	req.DataImage = dataImage
=======
>>>>>>> 473cd1d (uplaod image method)

	req.OriginalImage = originalImageName
	req.DataImage = dataImage
=======
>>>>>>> c73e2dc (pagination function)

	req.Name = name
	req.Type = t
	req.Breed = breed
	req.Healthnotes = healthnotes
	req.Gender = gender
=======
>>>>>>> 67140c6 (updated create pet)
	req.OriginalImage = header.Filename
	req.DataImage = dataImage
=======
>>>>>>> 473cd1d (uplaod image method)

	req.OriginalImage = originalImageName
	req.DataImage = dataImage

	res, err := c.service.CreatePet(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) GetPetByID(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}

	res, err := c.service.GetPetByID(ctx, petid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) ListPets(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	req := listPetsRequest{
		Type:  ctx.Query("type"),
		Breed: ctx.Query("breed"),
		// Age:    int(limit),
		// Weight: float64(offset),
	}

	pets, err := c.service.ListPets(ctx, req, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pets)
}

func (c *PetController) UpdatePet(ctx *gin.Context) {
	petid, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	var req updatePetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdatePet(ctx, petid, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Pet updated successfully"})
}

func (c *PetController) DeletePet(ctx *gin.Context) {
	petid, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	err = c.service.SetPetInactive(ctx, petid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Pet set to inactive successfully"})
}

func (c *PetController) ListPetsByUsername(ctx *gin.Context) {
	// username := ctx.Param("username")
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	fmt.Println(authPayload.Username)
>>>>>>> c73e2dc (pagination function)
=======
>>>>>>> 0637caf (upadted get lít)
=======
	fmt.Println(authPayload.Username)
>>>>>>> c73e2dc (pagination function)
=======
>>>>>>> 0637caf (upadted get lít)
	pets, err := c.service.ListPetsByUsername(ctx, authPayload.Username, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pets)
}

func (c *PetController) GetPetLogsByPetID(ctx *gin.Context) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	petidStr := ctx.Param("pet_id")
=======
	petidStr := ctx.Param("petid")
>>>>>>> 7e616af (add pet log schema)
=======
	petidStr := ctx.Param("pet_id")
>>>>>>> 2fe5baf (treatment phase)
=======
	petidStr := ctx.Param("petid")
>>>>>>> 7e616af (add pet log schema)
=======
	petidStr := ctx.Param("pet_id")
>>>>>>> 2fe5baf (treatment phase)
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.GetPetLogsByPetIDService(ctx, petid, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> 3835eb4 (update pet_schedule api)

func (c *PetController) InsertPetLog(ctx *gin.Context) {
	var req PetLog
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.service.InsertPetLogService(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Insert pet log successfully"})
}

func (c *PetController) DeletePetLog(ctx *gin.Context) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	logidStr := ctx.Param("log_id")
=======
=======
>>>>>>> 3835eb4 (update pet_schedule api)
	petidStr := ctx.Param("petid")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}
<<<<<<< HEAD
=======
=======
>>>>>>> 884b92e (update pet logs api)
	// petidStr := ctx.Param("petid")
	// petid, err := strconv.ParseInt(petidStr, 10, 64)
	// if err != nil {
	// 	ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
	// 	return
	// }
<<<<<<< HEAD
>>>>>>> 884b92e (update pet logs api)

=======
>>>>>>> eb8d761 (updated pet schedule)
	logidStr := ctx.Param("logid")
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	logidStr := ctx.Param("log_id")
>>>>>>> 2fe5baf (treatment phase)
=======
=======
>>>>>>> 884b92e (update pet logs api)

=======
>>>>>>> eb8d761 (updated pet schedule)
	logidStr := ctx.Param("logid")
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	logidStr := ctx.Param("log_id")
>>>>>>> 2fe5baf (treatment phase)
	logid, err := strconv.ParseInt(logidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	err = c.service.DeletePetLogService(ctx, logid)
=======
	err = c.service.DeletePetLogService(ctx, petid, logid)
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	err = c.service.DeletePetLogService(ctx, logid)
>>>>>>> 884b92e (update pet logs api)
=======
	err = c.service.DeletePetLogService(ctx, petid, logid)
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	err = c.service.DeletePetLogService(ctx, logid)
>>>>>>> 884b92e (update pet logs api)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Delete pet log successfully"})
}

func (c *PetController) UpdatePetLog(ctx *gin.Context) {
	var req PetLog
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	logidStr := ctx.Param("log_id")
=======
	logidStr := ctx.Param("logid")
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	logidStr := ctx.Param("log_id")
>>>>>>> 2fe5baf (treatment phase)
=======
	logidStr := ctx.Param("logid")
>>>>>>> 3835eb4 (update pet_schedule api)
=======
	logidStr := ctx.Param("log_id")
>>>>>>> 2fe5baf (treatment phase)
	logid, err := strconv.ParseInt(logidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

	err = c.service.UpdatePetLogService(ctx, req, logid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Update pet log successfully"})
}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

func (c *PetController) UpdatePetAvatar(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
=======

func (c *PetController) UpdatePetAvatar(ctx *gin.Context) {
<<<<<<< HEAD
	petidStr := ctx.Param("petid")
>>>>>>> e30b070 (Get list appoinment by user)
=======
	petidStr := ctx.Param("pet_id")
>>>>>>> 2fe5baf (treatment phase)
=======

func (c *PetController) UpdatePetAvatar(ctx *gin.Context) {
<<<<<<< HEAD
	petidStr := ctx.Param("petid")
>>>>>>> e30b070 (Get list appoinment by user)
=======
	petidStr := ctx.Param("pet_id")
>>>>>>> 2fe5baf (treatment phase)
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(400, util.ErrorResponse(err))
		return
	}

	err = c.service.UpdatePetAvatar(ctx, petid, updatePetAvatarRequest{
		DataImage:     dataImage,
		OriginalImage: originalImageName,
	})
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Update pet avatar successfully"})
}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 7e616af (add pet log schema)
=======
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
=======
>>>>>>> ffc9071 (AI suggestion)

func (c *PetController) GetPetProfileSummary(ctx *gin.Context) {
	petIDStr := ctx.Param("pet_id")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	summary, err := c.service.GeneratePetProfileSummary(ctx, petID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
<<<<<<< HEAD
>>>>>>> ffc9071 (AI suggestion)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
>>>>>>> 7e616af (add pet log schema)
=======
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
>>>>>>> ffc9071 (AI suggestion)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)

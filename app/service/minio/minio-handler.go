package minio

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

type FileResponse struct {
	ID   int64  `json:"id"`
	URL  string `json:"url"`
	Path string `json:"path"`
}

type MinioHandler struct {
	storeDB db.Store
}

func NewMinioHandler(router *gin.Engine) *MinioHandler {
	// Initialize with the global StoreDB instance
	return &MinioHandler{
		storeDB: db.StoreDB,
	}
}

// sanitizeBucketName ensures the bucket name follows S3 naming conventions
// - Contains only lowercase letters, numbers, dots (.), and hyphens (-)
// - Must not contain underscores or uppercase letters
// - Must be between 3 and 63 characters long
// - Must not be an IP address
// - Must not start or end with dot or hyphen
// - Must not contain two adjacent periods
// - Must not use a period adjacent to a hyphen
func sanitizeBucketName(name string) string {
	// If empty, use a default
	if name == "" {
		return "pet-files"
	}

	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any characters that aren't allowed
	reg := regexp.MustCompile("[^a-z0-9.-]")
	name = reg.ReplaceAllString(name, "")

	// Replace multiple dots with a single dot
	for strings.Contains(name, "..") {
		name = strings.ReplaceAll(name, "..", ".")
	}

	// Replace dot-hyphen or hyphen-dot with just a hyphen
	name = strings.ReplaceAll(name, ".-", "-")
	name = strings.ReplaceAll(name, "-.", "-")

	// Trim dots and hyphens from start and end
	name = strings.Trim(name, ".-")

	// Ensure it's at least 3 characters
	if len(name) < 3 {
		name = "pet-" + name
	}

	// If still too short, use a default
	if len(name) < 3 {
		name = "pet-files"
	}

	// Truncate if longer than 63 characters
	if len(name) > 63 {
		name = name[:63]
		// Ensure we don't end with dot or hyphen after truncation
		name = strings.TrimRight(name, ".-")
	}

	return name
}

func (s *MinioHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pet_id := c.PostForm("pet_id")
	pet_id_int, err := strconv.ParseInt(pet_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pet, err := s.storeDB.GetPetByID(c, pet_id_int)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sanitize bucket name
	bucketName := sanitizeBucketName(pet.Name)
	fileName := file.Filename
	objectType := file.Header.Get("Content-Type")
	objectSize := file.Size

	url, err := HandleFileUpload(c, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := s.storeDB.CreateFile(c, db.CreateFileParams{
		FileName: fileName,
		FileType: objectType,
		FileSize: objectSize,
		FilePath: bucketName + "/" + fileName,
		PetID:    pgtype.Int8{Int64: pet_id_int, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url, "file_id": res.ID})
}

func (s *MinioHandler) GetFile(c *gin.Context) {
	fileID := c.Param("file_id")
	fileIDInt, err := strconv.ParseInt(fileID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, err := s.storeDB.GetFileByID(c, fileIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url, err := FetchPreSignedURL(c, fileIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileResponse := FileResponse{
		ID:   file.ID,
		URL:  url,
		Path: file.FilePath,
	}

	c.JSON(http.StatusOK, gin.H{"file": fileResponse})
}

func (s *MinioHandler) GetFiles(c *gin.Context) {
	pet_id := c.Query("pet_id")
	pet_id_int, err := strconv.ParseInt(pet_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files, err := s.storeDB.GetFiles(c, pgtype.Int8{Int64: pet_id_int, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var fileResponses []FileResponse
	for _, file := range files {
		url, err := FetchPreSignedURL(c, file.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fileResponses = append(fileResponses, FileResponse{
			ID:   file.ID,
			URL:  url,
			Path: file.FilePath,
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": fileResponses})
}

func Routes(routerGroup middleware.RouterGroup, minioHandler *MinioHandler) {
	minio := routerGroup.RouterDefault.Group("/")

	// Authentication required routes
	authRoute := routerGroup.RouterAuth(minio)
	{
		authRoute.POST("upload", minioHandler.UploadFile)
		authRoute.GET("file/:file_id", minioHandler.GetFile)
		authRoute.GET("files", minioHandler.GetFiles)
	}
}

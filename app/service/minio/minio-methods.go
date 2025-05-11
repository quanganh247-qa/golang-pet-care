package minio

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// FetchPreSignedURLByProjectFileID fetches the pre-signed URL for a file based on the given ProjectFileID.
// It returns the pre-signed URL or an error if not found or failed to generate the URL.
func FetchPreSignedURL(c *gin.Context, id int64) (string, error) {
	// Step 1: Fetch the ProjectFile record from the database based on the given file ID
	mcclient, err := GetMinIOClient()
	if err != nil {
		return "", fmt.Errorf("error getting MinIO client: %v", err)
	}

	file, err := db.StoreDB.GetFileByID(c, id)
	if err != nil {
		return "", fmt.Errorf("error getting file: %v", err)
	}

	// Step 2: Ensure that the file path (S3 key) is present in the ProjectFile
	if file.FilePath == "" {
		// Log missing file path error and return error
		log.Println(color.RedString("File path is empty for file ID: %d", id))
		return "", errors.New("File path is empty")
	}

	expiry := time.Hour * 3 // URL validity: 3 hours

	// // Split the FilePath to get the bucket name (first part) and the object name (rest)
	// parts := strings.Split(file.FilePath, "/")
	// if len(parts) < 1 {
	// 	log.Println(color.RedString("Invalid file path format: %s", file.FilePath))
	// 	return "", errors.New("Invalid file path format")
	// }

	// // The first part is the bucket name
	// bucketName := parts[0]

	pet, err := db.StoreDB.GetPetByID(c, file.PetID.Int64)
	if err != nil {
		return "", fmt.Errorf("error getting pet: %v", err)
	}

	bucketName := sanitizeBucketName(pet.Name)

	// Generate pre-signed URL
	url, err := mcclient.GetPresignedURL(c, bucketName, file.FileName, expiry)
	if err != nil {
		// Log failure to generate pre-signed URL and return error
		log.Println(color.RedString("Failed to generate pre-signed URL for file path: %s", file.FilePath))
		return "", errors.New("Failed to generate pre-signed URL")
	}

	// Return the pre-signed URL
	return url, nil
}

// HandlePetAvatarFileUpload handles file uploads and returns the URL of the uploaded file to MinIO
func HandlePetAvatarFileUpload(c *gin.Context, tempFilePath, email, username string) (string, int64, error) {
	mcclient, err := GetMinIOClient()
	if err != nil {
		return "", 0, fmt.Errorf("error getting MinIO client: %v", err)
	}

	// Create a MinIO bucket for the project if it doesn't exist
	err = mcclient.CreateBucket(c, username)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create bucket %s: %v", username, err)
	}

	// Open the temporary file
	file, err := os.Open(tempFilePath)
	if err != nil {
		log.Println(color.RedString("Failed to open temporary file: %s", tempFilePath))
		return "", 0, fmt.Errorf("failed to open temporary file %s: %v", tempFilePath, err)
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(color.RedString("Failed to get file stats: %s", tempFilePath))
		return "", 0, fmt.Errorf("failed to get file stats for %s: %v", tempFilePath, err)
	}

	// Read file contents into a []byte
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Println(color.RedString("Failed to read file content: %s", tempFilePath))
		return "", 0, fmt.Errorf("failed to read file content %s: %v", tempFilePath, err)
	}

	fileName := "avatar"
	fileSize := fileInfo.Size()
	fileType := filepath.Ext(tempFilePath)

	// Upload the file to MinIO
	err = mcclient.UploadFile(c, username, fileName, fileContent)
	if err != nil {
		log.Println(color.RedString("Failed to upload file to MinIO: %s", fileName))
		return "", 0, fmt.Errorf("failed to upload file %s to MinIO: %v", fileName, err)
	}

	// Get the presigned URL to access the file
	imageURL, err := mcclient.GetPresignedURL(c, username, fileName, time.Duration(24)*time.Hour)
	if err != nil {
		log.Println(color.RedString("Failed to get presigned URL: %s", fileName))
		return "", 0, fmt.Errorf("failed to get presigned URL for file %s: %v", fileName, err)
	}

	// Create a ProjectFile record
	newFile, err := db.StoreDB.CreateFile(c, db.CreateFileParams{
		FileName: fileName,
		FilePath: username,
		FileSize: fileSize,
		FileType: fileType,
	})
	if err != nil {
		log.Println(color.RedString("Failed to create pet avatar: %v", err))
		return "", 0, fmt.Errorf("failed to create pet avatar")
	}

	// Return the presigned URL and the project file ID
	return imageURL, newFile.ID, nil
}

// UpdateCoverFileUpload handles file uploads and returns the URL of the uploaded file
func UpdateCoverFileUpload(c *gin.Context, email, username string, coverID int64) (string, int64, error) {
	mcclient, err := GetMinIOClient()
	if err != nil {
		return "", 0, fmt.Errorf("error getting MinIO client: %v", err)
	}

	mcclient.CreateBucket(c, username)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create bucket %s: %v", username, err)
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		return "", 0, fmt.Errorf("failed to get the file: %v", err)
	}
	defer file.Close()

	fileName := "cover"
	// Create ProjectFile record
	fileSize := fileHeader.Size                       // Get the file size
	fileType := fileHeader.Header.Get("Content-Type") // Get the file type

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read file content for %s: %v", fileName, err)
	}

	err = mcclient.UploadFile(c, username, fileName, fileContent)
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload file %s to MinIO: %v", fileName, err)
	}

	imageURL, err := mcclient.GetPresignedURL(c, username, fileName, time.Duration(24)*time.Hour)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get presigned URL for file %s: %v", fileName, err)
	}
	var projectFile db.File
	if coverID == 0 {
		projectFile = db.File{
			FileName: fileName,
			FilePath: username,
			FileSize: fileSize,
			FileType: fileType,
		}
		// Save the project to the database
		_, err := db.StoreDB.CreateFile(c, db.CreateFileParams{
			FileName: fileName,
			FilePath: username,
			FileSize: fileSize,
			FileType: fileType,
		})
		if err != nil {
			log.Println(color.RedString("Failed to create project cover page: %v", err))
			return "", 0, fmt.Errorf("failed to create project file")
		}
	} else {
		projectFile = db.File{
			ID:       coverID,
			FileName: fileName,
			FilePath: username,
			FileSize: fileSize,
			FileType: fileType,
		}
		_, err := db.StoreDB.UpdateFile(c, db.UpdateFileParams{
			ID:       coverID,
			FileName: fileName,
			FilePath: username,
			FileSize: fileSize,
			FileType: fileType,
		})
		if err != nil {
			log.Println(color.RedString("Failed to update project cover page: %v", err))
			return "", 0, fmt.Errorf("failed to update project file")
		}
	}

	return imageURL, projectFile.ID, nil
}

// HandleFileUpload handles file uploads and returns the URL of the uploaded file
func HandleFileUpload(c *gin.Context, bucketName string) (string, error) {
	mcclient, err := GetMinIOClient()
	if err != nil {
		return "", fmt.Errorf("error getting MinIO client: %v", err)
	}

	//check if bucket exists
	exists, err := mcclient.Client.BucketExists(c, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket %s: %v", bucketName, err)
	}
	if !exists {
		mcclient.CreateBucket(c, bucketName)
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("failed to get the file: %v", err)
	}
	defer file.Close()

	fileName := fileHeader.Filename
	// fileSize := fileHeader.Size
	// fileType := fileHeader.Header.Get("Content-Type")

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content for %s: %v", fileName, err)
	}

	err = mcclient.UploadFile(c, bucketName, fileName, fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to upload file %s to MinIO: %v", fileName, err)
	}

	url, err := mcclient.GetPresignedURL(c, bucketName, fileName, time.Duration(24)*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL for file %s: %v", fileName, err)
	}

	return url, nil
}

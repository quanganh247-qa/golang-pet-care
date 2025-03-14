package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	once sync.Once
	mc   *Client
)

// Client is a wrapper around the MinIO client
type Client struct {
	Client   *minio.Client
	Endpoint string
	Region   string
}

// NewMinIOClient initializes and returns the MinIO client (singleton)
func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*Client, error) {
	// Check for empty accessKey and secretKey before using once.Do()
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("accessKey and secretKey must not be empty")
	}

	once.Do(func() {
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Println(color.RedString("failed to create MinIO client: %v", err))
			return
		}

		mc = &Client{
			Client:   minioClient,
			Endpoint: endpoint,
		}
	})

	// Return the client or an error if initialization failed
	if mc == nil {
		return nil, fmt.Errorf("MinIO client initialization failed")
	}

	return mc, nil
}

// CheckConnection checks if the MinIO client is connected and operational
func (mc *Client) CheckConnection(ctx context.Context) error {
	_, err := mc.Client.ListBuckets(ctx)
	if err != nil {
		log.Println(color.RedString("failed to check connection: %v", err))
		return fmt.Errorf("failed to check connection: %v", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	log.Println(color.GreenString("MinIO client is connected successfully"))
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> ada3717 (Docker file)
=======
	log.Println(color.GreenString("MinIO client is connected successfully"))
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> ada3717 (Docker file)
	return nil
}

// CreateBucket creates a new bucket in MinIO and logs errors before exiting if any
func (mc *Client) CreateBucket(ctx context.Context, bucketName string) error {
	// Check if the bucket exists
	found, err := mc.Client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Println(color.RedString("Failed to check if bucket exists: %v", err))
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	// If bucket doesn't exist, create it
	if !found {
		err = mc.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: mc.Region, ObjectLocking: true})
		if err != nil {
			log.Println(color.RedString("Failed to create bucket: %v", err))
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	return nil
}

// DeleteBucket deletes a bucket from MinIO and logs errors before exiting if any
func (mc *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	// Check if the bucket exists
	found, err := mc.Client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Println(color.RedString("Failed to check if bucket exists: %v", err))
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	// If bucket doesn't exist, log info and skip deletion
	if !found {
		log.Println(color.YellowString("Bucket %s not found, skipping deletion", bucketName))
		return nil
	}

	// List objects in the bucket and delete them
	objectsCh := mc.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive:    true,
		WithVersions: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			log.Println(color.RedString("Error listing object: %v", object.Err))
			continue
		}

		err := mc.Client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{
			VersionID:        object.VersionID,
			GovernanceBypass: true,
		})
		if err != nil {
			log.Println(color.RedString("Error deleting object %s: %v", object.Key, err))
		}
	}

	// Attempt to remove the bucket
	err = mc.Client.RemoveBucket(ctx, bucketName)
	if err != nil {
		log.Println(color.RedString("Failed to delete bucket: %v", err))
		return fmt.Errorf("failed to delete bucket: %v", err)
	}

	log.Println(color.GreenString("Bucket %s deleted successfully", bucketName))
	return nil
}

// UploadFile uploads content as a file to a specified bucket and path and logs errors.
func (mc *Client) UploadFile(ctx context.Context, bucketName, objectName string, fileContent []byte) error {
	// Create a reader from the byte slice
	file := bytes.NewReader(fileContent)

	// Get the file size
	fileSize := int64(len(fileContent))

	// Log the upload initiation
	log.Println(color.GreenString("Uploading file %s to bucket %s", objectName, bucketName))

	// Upload the file to MinIO
	_, err := mc.Client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		// Log the error
		log.Println(color.RedString("Failed to upload file %s: %v", objectName, err))
		return fmt.Errorf("failed to upload file: %v", err)
	}

	// Log success
	log.Println(color.GreenString("Successfully uploaded file %s to bucket %s", objectName, bucketName))

	return nil
}

// DownloadFile downloads a file from a bucket and logs errors.
func (mc *Client) DownloadFile(ctx context.Context, bucketName, objectName, downloadPath string) error {
	// Ensure the directory for the download path exists
	err := os.MkdirAll(filepath.Dir(downloadPath), 0755)
	if err != nil {
		// Log the error
		log.Println(color.RedString("Failed to create directory %s: %v", filepath.Dir(downloadPath), err))
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Get the object from the bucket
	object, err := mc.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		// Log the error
		log.Println(color.RedString("Failed to get object %s from bucket %s: %v", objectName, bucketName, err))
		return fmt.Errorf("failed to get object: %v", err)
	}
	defer object.Close()

	// Create the file to save the downloaded object
	file, err := os.Create(downloadPath)
	if err != nil {
		// Log the error
		log.Println(color.RedString("Failed to create file %s: %v", downloadPath, err))
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy the object content to the file
	_, err = io.Copy(file, object)
	if err != nil {
		// Log the error
		log.Println(color.RedString("Failed to download file %s: %v", objectName, err))
		return fmt.Errorf("failed to download file: %v", err)
	}

	// Log the success of the download
	log.Println(color.GreenString("Successfully downloaded file %s from bucket %s to %s", objectName, bucketName, downloadPath))

	return nil
}

// GetPresignedURL generates a presigned URL for an object in the specified bucket
func (mc *Client) GetPresignedURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	presignedURL, err := mc.Client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		log.Println(color.RedString("failed to generate presigned URL for object %s: %v", objectName, err))
		return "", fmt.Errorf("failed to generate presigned URL for object %s: %v", objectName, err)
	}
	return presignedURL.String(), nil
}

// GetMinIOClient returns the singleton MinIO client
func GetMinIOClient() (*Client, error) {
	if mc == nil {
		return nil, fmt.Errorf("MinIO client is not initialized")
	}
	return mc, nil
}

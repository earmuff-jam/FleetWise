package bucket

import (
	"os"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/minio/minio-go"
)

// InitializeStorageAndBucket...
//
// Allows for the creation of storage container and bucket
func InitializeStorageAndBucket() {
	client, err := initializeStorage()
	if err != nil {
		config.Log("unable to initialize minio client storage", err)
		return
	}
	initializeBucket(client)
}

// initializeStorage ...
//
// Initializes MinIO bucket storage
func initializeStorage() (*minio.Client, error) {

	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	endpoint := os.Getenv("MINIO_APP_LOCALHOST_URL")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")

	config.Log("setting up bucket storage for user %s", nil, accessKeyID)
	useSSL := false // true for HTTPS

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		config.Log("Failed to initialize MinIO client", err)
		return nil, err
	}

	config.Log("Connected to MinIO bucket storage", nil)
	return minioClient, nil
}

// initializeBucket ...
//
// creates new bucket if the bucket does not exist.
// if the bucket already exists, we do not create the new bucket
func initializeBucket(minioClient *minio.Client) {
	bucketName := os.Getenv("MINIO_APP_BUCKET_NAME")
	location := os.Getenv("MINIO_APP_BUCKET_LOCATION")

	err := minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check if the bucket already exists
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			config.Log("Selected bucket %s already exists.", nil, bucketName)
		} else {
			config.Log("Failed to create bucket", err)
		}
	}
}

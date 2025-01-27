package bucket

import (
	"errors"
	"io"
	"os"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/minio/minio-go"
)

// UploadDocumentInBucket ...
//
// uploads document in the bucket. function attempts to communicate with the bucket
func UploadDocumentInBucket(objectName string, filePath string, contentType string) error {

	client, err := initializeStorage()
	if err != nil {
		config.Log("unable to initialize minio client storage", err)
		return err
	}
	bucketName := os.Getenv("MINIO_APP_BUCKET_NAME")

	_, err = client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		config.Log("unable to add object to the selected bucket", err)
		return err
	}
	config.Log("upload successful", nil)
	return nil
}

// RetrieveDocumentFromBucket ...
//
// Retrieves the selected document from the bucket storage
func RetrieveDocumentFromBucket(documentID string) ([]byte, string, string, error) {
	client, err := initializeStorage()
	if err != nil {
		config.Log("unable to initialize minio client storage", err)
		return nil, "", "", err
	}
	bucketName := os.Getenv("MINIO_APP_BUCKET_NAME")

	// Fetch the object from the bucket
	object, err := client.GetObject(bucketName, documentID, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			config.Log("Object not found: %s", err, documentID)
			return nil, "", "", nil // Gracefully return empty data if object doesn't exist
		}
		config.Log("unable to retrieve object from the bucket", err)
		return nil, "", "", err
	}
	defer object.Close()

	// Get object information
	objectStat, err := object.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			config.Log("Object metadata not found: %s", err, documentID)
			return nil, "", "", errors.New("NoSuchKey") // Catch the error code and return the error code
		}
		config.Log("unable to retrieve object metadata", err)
		return nil, "", "", err
	}

	// Read the content
	content, err := io.ReadAll(object)
	if err != nil {
		config.Log("Error reading object content", err)
		return nil, "", "", err
	}

	return content, objectStat.ContentType, objectStat.Key, nil
}

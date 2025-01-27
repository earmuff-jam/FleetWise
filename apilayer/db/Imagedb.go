package db

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/earmuff-jam/fleetwise/bucket"
	"github.com/earmuff-jam/fleetwise/config"
)

// UploadImage ...
func UploadImage(file multipart.File, header *multipart.FileHeader, userID string) error {

	tempFilePath := "/tmp/" + header.Filename
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		config.Log("Unable to create temporary file", err)
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		config.Log("Unable to copy file content", err)
		return err
	}

	contentType := header.Header.Get("Content-Type")
	err = bucket.UploadDocumentInBucket(userID, tempFilePath, contentType)
	if err != nil {
		config.Log("Unable to upload document", err)
		return err
	}
	// cleanup temp file
	defer os.Remove(tempFilePath)
	return nil
}

// FetchImage ...
func FetchImage(userID string) ([]byte, string, string, error) {

	content, contentType, fileName, err := bucket.RetrieveDocumentFromBucket(userID)
	if err != nil {
		if err.Error() == "NoSuchKey" {
			config.Log("cannot find the selected document", err)
			return nil, "", "", err
		}
		config.Log("unable to retrieve the selected document", err)
		return nil, "", "", err
	}
	return content, contentType, fileName, nil
}

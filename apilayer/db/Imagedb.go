package db

import (
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/earmuff-jam/fleetwise/bucket"
)

// UploadImage ...
func UploadImage(file multipart.File, header *multipart.FileHeader, userID string) error {

	tempFilePath := "/tmp/" + header.Filename
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Unable to create temporary file. Error: %+v", err)
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Printf("Unable to copy file content. Error: %+v", err)
		return err
	}

	contentType := header.Header.Get("Content-Type")
	err = bucket.UploadDocumentInBucket(userID, tempFilePath, contentType)
	if err != nil {
		log.Printf("Unable to upload document. Error: %+v", err)
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
			log.Printf("cannot find the selected document. error: %+v", err)
			return nil, "", "", err
		}
		log.Printf("unable to retrieve the selected document. error: %+v", err)
		return nil, "", "", err
	}
	return content, contentType, fileName, nil
}

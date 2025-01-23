package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/earmuff-jam/fleetwise/db"
	"github.com/gorilla/mux"
)

// UploadImage ...
// swagger:route POST /api/v1/{id}/upload Images UploadImage
//
// # Uploads an image with the selected id as filename. Used to persist user profile pictures, category or maintenance plan images
//
// Parameters:
//   - +name: id
//     in: path
//     description: The unique ID identifying the selected image object
//     type: string
//     required: true
//   - +name: FileHeader
//     in: body
//     description: The full file details of the selected image object
//     type: string
//     format: byte
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UploadImage(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	ID := vars["id"]

	if len(ID) <= 0 {
		log.Printf("Unable to upload image with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	file, header, err := r.FormFile("imageSrc")
	if err != nil {
		log.Printf("Unable to retrieve file. Error: %+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	defer file.Close()

	err = db.UploadImage(file, header, ID)
	if err != nil {
		log.Printf("Unable to upload image. error: %+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("File uploaded successfully"))
}

// FetchImage ...
// swagger:route GET /api/v1/{id}/fetchImage Images FetchImage
//
// # Retrives the current image that matches the passed in ID. If the image does not exist, the application will still return a HTTP ok response to ensure that the UI will not fail due to missing image
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected image
//     type: string
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func FetchImage(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	ID := vars["id"]

	if len(ID) <= 0 {
		log.Printf("Unable to retrieve image without a valid ID")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	content, contentType, fileName, err := db.FetchImage(ID)
	if err != nil {
		if err.Error() == "NoSuchKey" {
			log.Printf("cannot find the selected document. error: %+v", err)
			// Write response headers and content only
			rw.Header().Set("Content-Type", contentType)
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("NoSuchKey"))
			return
		}
		log.Printf("Failed to retrieve image. error: %+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Write response headers and content
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Content-Disposition", "inline; filename="+fileName)
	rw.WriteHeader(http.StatusOK)
	rw.Write(content)
}

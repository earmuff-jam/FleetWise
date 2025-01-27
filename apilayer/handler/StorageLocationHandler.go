package handler

import (
	"encoding/json"
	"net/http"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
)

// GetAllStorageLocations ...
// swagger:route GET /api/v1/locations StorageLocations getAllStorageLocations
//
// # Retrieves the list of locations that users can use to store an asset.
//
// Responses:
// 200: []StorageLocation
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAllStorageLocations(rw http.ResponseWriter, r *http.Request, user string) {

	resp, err := db.RetrieveAllStorageLocation(user)
	if err != nil {
		config.Log("Unable to retrieve storage locations", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

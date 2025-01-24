package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/earmuff-jam/fleetwise/db"
)

// GetAssetsAndSummary ...
// swagger:route GET /api/v1/locations Summary GetAssetsAndSummary
//
// # Retrieves list of categories and its associated asset names, list of maintenance plans and its associated assets and
// list of assets that are past due. For the asset list that is past due, if there is no returnable flag set, then it does not
// return that values. Also returns an array of existing assets.
// All results are created by the selected user and in order of type and date last updated at.
//
// Parameters:
//   - +name: id
//     in: query
//     description: The id of the selected user
//     type: string
//     required: true
//
// Responses:
// 200: []AssetSummary
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAssetsAndSummary(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")

	if userID == "" {
		log.Printf("Unable to retrieve asset summary with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	resp, err := db.RetrieveAssetsAndSummary(user, userID)
	if err != nil {
		log.Printf("Unable to retrieve asset summary. error: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

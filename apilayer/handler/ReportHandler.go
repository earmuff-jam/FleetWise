package handler

import (
	"encoding/json"
	"net/http"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetReports ...
// swagger:route GET /api/reports/{id} Reports getReports
//
// # Retrieves the reports on the assets created by the selected user based on the query parameters.
//
// Parameters:
//   - +name: id
//     in: path
//     description: The userID of the selected user
//     type: string
//     required: true
//   - +name: since
//     in: query
//     description: The start date for the result to begin from
//     type: string
//     required: true
//   - +name: includeOverdue
//     in: query
//     description: The boolean flag to determine if query response should consider overdue items or not.
//     type: boolean
//     required: true
//
// Responses:
// 200: []Report
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetReports(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	sinceDateTime := r.URL.Query().Get("since")
	includeOverdueAssets := r.URL.Query().Get("includeOverdue")

	if !ok || len(id) <= 0 {
		config.Log("Unable to retrieve details without an id.", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		config.Log("Unable to retrieve details with provided id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if sinceDateTime == "" {
		config.Log("unable to retrieve details without a start date time.", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if len(includeOverdueAssets) == 0 {
		includeOverdueAssets = "false"
	}

	resp, err := db.RetrieveReports(user, parsedUUID, sinceDateTime, includeOverdueAssets)
	if err != nil {
		config.Log("Unable to retrieve report details", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

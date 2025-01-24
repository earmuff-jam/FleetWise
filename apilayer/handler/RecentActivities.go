package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/earmuff-jam/fleetwise/db"
	"github.com/gorilla/mux"
)

// GetRecentActivities ...
// swagger:route GET /api/v1/profile/recent-activities GetRecentActivities getRecentActivities
//
// # Retrieves user actions for categories, maintenance plans, and assets. List all created, updated, or deleted.
//
// Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: limit
//     in: query
//     description: The limit of recent activities.
//     required: true
//     type: integer
//     format: int32
//   - +name: until
//     in: query
//     description: The timestamp with time zone until the data should be retrieved for
//     required: false
//     type: string
//     format: date-time
//
// Responses:
//
//	200: []RecentActivity
//	400: MessageResponse
//	404: MessageResponse
//	500: MessageResponse
func GetRecentActivities(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	userID := vars["id"]
	limit := r.URL.Query().Get("limit")
	untilDate := r.URL.Query().Get("until")

	if userID == "" {
		log.Printf("Unable to retrieve recent activities with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("Unable to retrieve recent activities without limit. error :%+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveRecentActivities(user, userID, limitInt, untilDate)
	if err != nil {
		log.Printf("Unable to retrieve recent activities. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

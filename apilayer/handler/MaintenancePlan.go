package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/gorilla/mux"
)

// GetAllMaintenancePlans ...
// swagger:route GET /api/v1/maintenance-plans MaintenancePlans getAllMaintenancePlans
//
// # Retrieves the list of maintenance plans that each asset can be associated with.
// Each user can have thier own set of maintenance plans. All plans are specific to the selected user
//
// Users can assign asset to multiple plans.
//
// // Parameters:
//   - +name: id
//     in: query
//     description: The userID of the logged in user
//     required: true
//     type: string
//   - +name: limit
//     in: query
//     description: The limit of maintenance plans
//     required: true
//     type: integer
//     format: int32
//
// Responses:
//
// 200: []MaintenancePlan
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAllMaintenancePlans(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	limit := r.URL.Query().Get("limit")

	if userID == "" {
		log.Printf("Unable to retrieve maintenance plans with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	resp, err := db.RetrieveAllMaintenancePlans(user, userID, limitInt)
	if err != nil {
		log.Printf("Unable to retrieve maintenance plans. error: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// GetMaintenancePlan ...
// swagger:route GET /api/v1/plan MaintenancePlans GetMaintenancePlan
//
// # Retrieve a selected maintenance plan
//
// // Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: mID
//     in: query
//     description: The maintenance id of the selected plan
//     required: true
//     type: string
//
// Responses:
// 200: MaintenancePlan
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	maintenanceID := r.URL.Query().Get("mID")

	if userID == "" {
		log.Printf("Unable to retrieve selected maintenance plan with empty user id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if maintenanceID == "" {
		log.Printf("Unable to retrieve selected maintenance plan with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveMaintenancePlan(user, userID, maintenanceID)
	if err != nil {
		log.Printf("unable to create new maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// GetAllMaintenancePlanItems ...
// swagger:route GET /api/v1/plans/items MaintenancePlans getAllMaintenancePlanItems
//
// # Retrieves the list of assets for a specific maintenance plans.
//
// // Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: limit
//     in: query
//     description: The limit of maintenance plans
//     required: true
//     type: integer
//     format: int32
//   - +name: mID
//     in: query
//     description: The maintenance id of the selected plan
//     required: true
//     type: string
//
// Responses:
//
// 200: []MaintenanceItemResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAllMaintenancePlanItems(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	limit := r.URL.Query().Get("limit")
	maintenancePlanID := r.URL.Query().Get("mID")

	if userID == "" {
		log.Printf("Unable to retrieve associated item for selected maintenance item with empty user id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if maintenancePlanID == "" {
		log.Printf("Unable to retrieve associated items with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	resp, err := db.RetrieveAllMaintenancePlanItems(user, userID, maintenancePlanID, limitInt)
	if err != nil {
		log.Printf("Unable to retrieve associated items. error: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// AddItemsInMaintenancePlan ...
// swagger:route POST /api/v1/category/items MaintenancePlans addItemsInMaintenancePlan
//
// # Add selected items in a specific maintenance plan
//
// Parameters:
//   - +name: MaintenanceItemRequest
//     in: body
//     description:The object containing the array of assets to be removed from the association for maintenance plans
//     type: MaintenanceItemRequest
//     required: true
//
// Responses:
// 200: []MaintenanceItemResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func AddItemsInMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	draftMaintenancePlan := &model.MaintenanceItemRequest{}
	err := json.NewDecoder(r.Body).Decode(draftMaintenancePlan)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.AddAssetToMaintenancePlan(user, draftMaintenancePlan)
	if err != nil {
		log.Printf("Unable to add assets to existing maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveAssociationFromMaintenancePlan ...
// swagger:route POST /api/v1/plan/remove/items MaintenancePlans RemoveAssociationFromMaintenancePlan
//
// # Removes association between the maintenance plan and asset item
//
// Parameters:
//   - +name: MaintenanceItemRequest
//     in: body
//     description: The object containing the array of assets to be removed
//     type: MaintenanceItemRequest
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveAssociationFromMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	draftMaintenancePlanItemRequest := &model.MaintenanceItemRequest{}
	err := json.NewDecoder(r.Body).Decode(draftMaintenancePlanItemRequest)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	err = db.RemoveAssetAssociationFromMaintenancePlan(user, draftMaintenancePlanItemRequest)
	if err != nil {
		log.Printf("Unable to remove assets from selected maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}

// CreateMaintenancePlan ...
// swagger:route POST /api/v1/plan MaintenancePlans createMaintenancePlan
//
// # Create maintenance plan
//
// Parameters:
//   - +name: MaintenancePlan
//     in: body
//     description: The object containing details of the maintenance plan
//     type: MaintenancePlan
//     required: true
//
// Responses:
// 200: MaintenancePlan
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func CreateMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	draftMaintenancePlan := &model.MaintenancePlan{}
	err := json.NewDecoder(r.Body).Decode(draftMaintenancePlan)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.CreateMaintenancePlan(user, draftMaintenancePlan)
	if err != nil {
		log.Printf("unable to create new maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// UpdateMaintenancePlan ...
// swagger:route PUT /api/v1/plan/{id} MaintenancePlans UpdateMaintenancePlan
//
// # Update maintenance plan function updates the selected maintenance plan with new values
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected maintenance plan to update
//     type: string
//     required: true
//   - +name: MaintenancePlan
//     in: body
//     description: The object containing details of the maintenance plan
//     type: MaintenancePlan
//     required: true
//
// Responses:
// 200: MaintenancePlan
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UpdateMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	maintenancePlanID := vars["id"]

	if len(maintenancePlanID) <= 0 {
		log.Printf("Unable to update maintenance plan with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	draftMaintenancePlan := &model.MaintenancePlan{}
	err := json.NewDecoder(r.Body).Decode(draftMaintenancePlan)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.UpdateMaintenancePlan(user, draftMaintenancePlan)
	if err != nil {
		log.Printf("Unable to update new maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveMaintenancePlan ...
// swagger:route DELETE /api/v1/plan/{id} MaintenancePlans removeMaintenancePlan
//
// # Removes a selected maintenance plan based on the id
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the maintenance plan to delete
//     type: string
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveMaintenancePlan(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	maintenancePlanID := vars["id"]

	if len(maintenancePlanID) <= 0 {
		log.Printf("Unable to delete maintenance plan with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	err := db.RemoveMaintenancePlan(user, maintenancePlanID)
	if err != nil {
		log.Printf("Unable to remove maintenance plan. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(maintenancePlanID)
}

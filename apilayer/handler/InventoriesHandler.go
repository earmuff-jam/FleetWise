package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/gorilla/mux"
)

const defaultHiddenStatus = "HIDDEN"

// GetAllInventories ...
// swagger:route GET /api/v1/profile/{id}/inventories Assets getAllInventories
//
// # Retrieves the list of assets that belong to the selected user
//
// // Parameters:
//   - +name: id
//     in: path
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: until
//     in: query
//     description: The timestamp with time zone until the data should be retrieved for. Returns all inventories if not passed in.
//     required: false
//     type: string
//     format: date-time
//
// Responses:
// 200: []Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAllInventories(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	userID := vars["id"]
	sinceDateTime := r.URL.Query().Get("since")

	if len(userID) <= 0 {
		config.Log("Unable to retrieve assets with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveAllInventoriesForUser(user, userID, sinceDateTime)
	if err != nil {
		config.Log("Unable to retrieve all existing assets", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// GetInventoryByID ...
// swagger:route GET /api/v1/profile/{id}/inventories/{invID} Assets getInventoryByID
//
// # Retrieves the selected inventory that matches the passed in ID created by the user.
//
// // Parameters:
//   - +name: id
//     in: path
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: invID
//     in: path
//     description: The inventoryID of the inventory details.
//     required: true
//     type: string
//
// Responses:
// 200: Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetInventoryByID(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	userID := vars["id"]
	invID := vars["invID"]

	if len(userID) <= 0 {
		config.Log("Unable to retrieve inventories with empty profile id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	if len(invID) <= 0 {
		config.Log("Unable to retrieve inventories with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveSelectedInv(user, userID, invID)
	if err != nil {
		config.Log("Unable to retrieve inventory with selected id", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// UpdateAssetColumn ...
// swagger:route GET /api/v1/profile/{id}/inventories/{assetID} Assets updateAssetColumn
//
// # Updates the selected quantity or price column field with passed in data field. All colaborators can update selected field.
//
// // Parameters:
//   - +name: id
//     in: path
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: assetID
//     in: path
//     description: The asssetID of the asset details.
//     required: true
//     type: string
//   - +name: UpdateAssetColumn
//     in: body
//     description: The object containing the details to update the selected asset with.
//     required: true
//     type: UpdateAssetColumn
//
// Responses:
// 200: Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UpdateAssetColumn(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	userID := vars["id"]
	assetID := vars["asssetID"]

	if len(userID) <= 0 {
		config.Log("Unable to retrieve assets with empty profile id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	if len(assetID) <= 0 {
		config.Log("Unable to retrieve assets with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var updateAssetCol model.UpdateAssetColumn
	if err := json.NewDecoder(r.Body).Decode(&updateAssetCol); err != nil {
		config.Log("unable to decode selected data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	isValidColumnName := false
	for _, v := range []string{"quantity", "price"} {
		if updateAssetCol.ColumnName == v {
			isValidColumnName = true
		}
	}

	if !isValidColumnName {
		config.Log("unable to update column name", nil)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := db.UpdateAsset(user, userID, updateAssetCol)
	if err != nil {
		config.Log("Unable to update asset", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// AddInventoryInBulk ...
// swagger:route POST /api/profile/{id}/inventories/bulk Assets addInventoryInBulk
//
// # Bulk upload inventory list. returns the list of uploaded items.
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected inventory
//     type: string
//     required: true
//   - +name: InventoryListRequest
//     in: body
//     description: The list of inventories to add into the db to support bulk upload
//     type: InventoryListRequest
//     required: true
//
// Responses:
//
// 200: []Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func AddInventoryInBulk(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to add new item with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var inventoryMap map[string]model.RawInventory
	if err := json.NewDecoder(r.Body).Decode(&inventoryMap); err != nil {
		config.Log("unable to decode data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var inventoryList []model.Inventory
	for _, v := range inventoryMap {

		draftInventory := model.Inventory{
			Name:           v.Name,
			Description:    v.Description,
			Price:          v.Price,
			Quantity:       int(v.Quantity),
			Location:       v.StorageLocation,
			Color:          v.Color,
			SKU:            v.SKU,
			Barcode:        v.Barcode,
			BoughtAt:       v.PurchaseLocation,
			MaxWeight:      v.MaximumWeight,
			MinWeight:      v.MinimumWeight,
			MaxHeight:      v.MaximumHeight,
			MinHeight:      v.MinimumHeight,
			Status:         defaultHiddenStatus,
			CreatedAt:      time.Now(),
			CreatedBy:      userID,
			UpdatedAt:      time.Now(),
			UpdatedBy:      userID,
			SharableGroups: []string{userID},
		}
		inventoryList = append(inventoryList, draftInventory)
	}

	inventoryListRequest := model.InventoryListRequest{
		InventoryList: inventoryList,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
	}

	// Assuming db.AddInventoryInBulk returns an error
	resp, err := db.AddInventoryInBulk(user, userID, inventoryListRequest)
	if err != nil {
		config.Log("unable to add new item during bulk insert", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// AddNewInventory ...
// swagger:route POST /api/profile/{id}/inventories Assets addNewInventory
//
// # Add new personal item to the user inventory
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected user
//     type: string
//     required: true
//   - +name: Inventory
//     in: body
//     description: The inventory object to add into the db
//     type: object
//     required: true
//
// Responses:
// 200: Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func AddNewInventory(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to add new item with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var inventory model.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
		config.Log("unable to decode selected data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := db.AddInventory(user, userID, inventory)
	if err != nil {
		config.Log("Unable to add new item", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// UpdateSelectedInventory ...
// swagger:route PUT /api/profile/{id}/inventories Assets updateSelectedInventory
//
// # Update selected inventory with details.
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected user
//     type: string
//     required: true
//   - +name: Inventory
//     in: body
//     description: The inventory object to add into the db
//     type: Inventory
//     required: true
//
// Responses:
// 200: Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UpdateSelectedInventory(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to update selected inventory without id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var inventory model.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
		config.Log("Error decoding data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := db.UpdateInventory(user, userID, inventory)
	if err != nil {
		config.Log("Unable to update selected inventory", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveSelectedInventory ...
// swagger:route POST /api/profile/{id}/inventories Assets removeSelectedInventory
//
// # Remove selected inventories.
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected user
//     type: string
//     required: true
//   - +name: pruneInventoryIDMap
//     in: body
//     description: The inventory object to add into the db
//     type: Inventory
//     required: true
//
// Responses:
// 200: Inventory
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveSelectedInventory(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to remove selected inventory without id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var pruneInventoriesIDMap map[string]string
	if err := json.NewDecoder(r.Body).Decode(&pruneInventoriesIDMap); err != nil {
		config.Log("Error decoding data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var pruneInventoriesIDs []string
	for _, v := range pruneInventoriesIDMap {
		pruneInventoriesIDs = append(pruneInventoriesIDs, v)
	}

	resp, err := db.DeleteInventory(user, userID, pruneInventoriesIDs)
	if err != nil {
		config.Log("Unable to remove selected inventory", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

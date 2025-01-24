package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/gorilla/mux"
)

// GetAllCategories ...
// swagger:route GET /api/v1/categories Categories getAllCategories
//
// Retrieves the list of categories that each asset can be associated with. Each user can have their own set of categories. All categories are specific to the selected user.
//
// Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: limit
//     in: query
//     description: The limit of categories
//     required: true
//     type: integer
//     format: int32
//
// Responses:
//
//	200: []Category
//	400: MessageResponse
//	404: MessageResponse
//	500: MessageResponse
func GetAllCategories(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	limit := r.URL.Query().Get("limit")

	if userID == "" {
		log.Printf("Unable to retrieve categories with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	resp, err := db.RetrieveAllCategories(user, userID, limitInt)
	if err != nil {
		log.Printf("Unable to retrieve categories. error: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// GetCategory ...
// swagger:route GET /api/v1/category Categories GetCategory
//
// # Retrieves the selected category
//
// // Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: catID
//     in: query
//     description: The category id
//     required: true
//     type: string
//
// Responses:
//
// 200: Category
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetCategory(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	categoryID := r.URL.Query().Get("catID")

	if userID == "" {
		log.Printf("Unable to retrieve associated item for selected category with empty user id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if categoryID == "" {
		log.Printf("Unable to retrieve associated items with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveCategory(user, userID, categoryID)
	if err != nil {
		// if there are no rows, we still want to return blank category
		if err == sql.ErrNoRows {
			log.Printf("unable to retrieve selected category. error: %+v", err)
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(resp)
			return
		} else {
			log.Printf("Unable to retrieve categories. error: %+v", err)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(err)
			return
		}

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// GetAllCategoryItems ...
// swagger:route GET /api/v1/category/items Categories getAllCategoryItems
//
// # Retrieves the list of assets for a specific category.
//
// // Parameters:
//   - +name: id
//     in: query
//     description: The userID of the selected user
//     required: true
//     type: string
//   - +name: limit
//     in: query
//     description: The limit of categories
//     required: true
//     type: integer
//     format: int32
//   - +name: catID
//     in: query
//     description: The category id of the selected plan
//     required: true
//     type: string
//
// Responses:
//
// 200: []CategoryItemResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetAllCategoryItems(rw http.ResponseWriter, r *http.Request, user string) {

	userID := r.URL.Query().Get("id")
	limit := r.URL.Query().Get("limit")
	categoryID := r.URL.Query().Get("catID")

	if userID == "" {
		log.Printf("Unable to retrieve associated item for selected category with empty user id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if categoryID == "" {
		log.Printf("Unable to retrieve associated items with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	resp, err := db.RetrieveAllCategoryItems(user, userID, categoryID, limitInt)
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

// CreateCategory ...
// swagger:route POST /api/v1/category Categories createCategory
//
// # Create category
//
// Parameters:
//   - +name: Category
//     in: body
//     description: The category object to save
//     type: Category
//     required: true
//
// Responses:
// 200: Category
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func CreateCategory(rw http.ResponseWriter, r *http.Request, user string) {

	draftCategory := &model.Category{}
	err := json.NewDecoder(r.Body).Decode(draftCategory)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.CreateCategory(user, draftCategory)
	if err != nil {
		log.Printf("Unable to create new category. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// AddItemsInCategory ...
// swagger:route POST /api/v1/category/items Categories addItemsInCategory
//
// # Add selected items in a specific category
//
// Parameters:
//   - +name: CategoryItemRequest
//     in: body
//     description: The object containing the array of assets to be removed from the association for categories
//     type: CategoryItemRequest
//     required: true
//
// Responses:
// 200: []CategoryItem
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func AddItemsInCategory(rw http.ResponseWriter, r *http.Request, user string) {

	draftCategory := &model.CategoryItemRequest{}
	err := json.NewDecoder(r.Body).Decode(draftCategory)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.AddAssetToCategory(user, draftCategory)
	if err != nil {
		log.Printf("Unable to add assets to existing category. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveAssociationFromCategory ...
// swagger:route POST /api/v1/category/remove/items Categories RemoveAssociationFromCategory
//
// # Removes association between the category and asset item
//
// Parameters:
//   - +name: CategoryItemRequest
//     in: body
//     description: The object containing the array of assets to update with the userID
//     type: CategoryItemRequest
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveAssociationFromCategory(rw http.ResponseWriter, r *http.Request, user string) {

	draftCategoryItemRequest := &model.CategoryItemRequest{}
	err := json.NewDecoder(r.Body).Decode(draftCategoryItemRequest)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	err = db.RemoveAssetAssociationFromCategory(user, draftCategoryItemRequest)
	if err != nil {
		log.Printf("Unable to remove assets from selected category. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}

// UpdateCategory ...
// swagger:route PUT /api/v1/category/{id} Categories updateCategory
//
// # Update category function updates the selected category with new values
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected category to update for
//     type: string
//     required: true
//   - +name: Category
//     in: body
//     description: The category object to update details for
//     type: Category
//     required: true
//
// Responses:
// 200: Category
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UpdateCategory(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	categoryID := vars["id"]

	if len(categoryID) <= 0 {
		log.Printf("Unable to update category with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	draftCategory := &model.Category{}
	err := json.NewDecoder(r.Body).Decode(draftCategory)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	resp, err := db.UpdateCategory(user, draftCategory)
	if err != nil {
		log.Printf("Unable to update new category. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveCategory ...
// swagger:route DELETE /api/v1/category/{id} Categories removeCategory
//
// # Removes a selected category based on the id
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the category to delete
//     type: string
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveCategory(rw http.ResponseWriter, r *http.Request, user string) {

	vars := mux.Vars(r)
	categoryID := vars["id"]

	if len(categoryID) <= 0 {
		log.Printf("Unable to delete category with empty id")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	err := db.RemoveCategory(user, categoryID)
	if err != nil {
		log.Printf("Unable to remove category. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return

	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(categoryID)
}

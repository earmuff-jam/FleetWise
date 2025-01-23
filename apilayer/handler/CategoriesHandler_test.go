package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllCategories(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllCategories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var categories []model.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(categories), 0)
}

func Test_GetAllCategories_IncorrectCategoryID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllCategories(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllCategories_InvalidDBUser(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllCategories(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetCategory(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllCategories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var categories []model.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(categories), 0)

	selectedCategory := categories[0]

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/category?id=%s&catID=%s", prevUser.ID.String(), selectedCategory.ID), nil)
	w = httptest.NewRecorder()
	GetAllCategoryItems(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var associatedCategoryItems []model.Category
	err = json.Unmarshal(data, &associatedCategoryItems)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(associatedCategoryItems), 0)
}

func Test_GetCategory_IncorrectCategoryID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetCategory_InvalidDBUser(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllCategoryItems(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllCategories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var categories []model.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(categories), 0)

	selectedCategory := categories[0]

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d&catID=%s", prevUser.ID.String(), 5, selectedCategory.ID), nil)
	w = httptest.NewRecorder()
	GetAllCategoryItems(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var associatedCategoryItems []model.CategoryItem
	err = json.Unmarshal(data, &associatedCategoryItems)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(associatedCategoryItems), 0) // does not test if items are associated with selected category
}

func Test_GetAllCategoryItems_IncorrectCategoryID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllCategoryItems(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllCategoryItems_InvalidDBUser(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllCategoryItems(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_CreateCategory(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	draftCategory := model.Category{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "general",
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateCategory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedCategory model.Category
	err = json.Unmarshal(data, &selectedCategory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedCategory.Name, "Kitty litter box")
	assert.Equal(t, selectedCategory.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedCategory.StatusName, "general")

	// cleanup
	db.RemoveCategory(config.CTO_USER, selectedCategory.ID)
}

func Test_CreateCategory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	CreateCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_CreateCategory_InvalidDBUser(t *testing.T) {

	draftCategory := model.Category{}
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddItemsInCategory(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	draftCategory := model.Category{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "general",
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateCategory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedCategory model.Category
	err = json.Unmarshal(data, &selectedCategory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedCategory.Name, "Kitty litter box")
	assert.Equal(t, selectedCategory.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedCategory.StatusName, "general")
	assert.Equal(t, len(selectedCategory.SharableGroups), 1)
	assert.Equal(t, selectedCategory.SharableGroups[0], prevUser.ID.String())

	// retrieve inventory list
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	GetAllInventories(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var inventories []model.Inventory
	err = json.Unmarshal(data, &inventories)
	if err != nil {
		t.Errorf("expected error to be nil but got %+v", err)
	}

	selectedInventory := inventories[0]

	draftAssets := model.CategoryItemRequest{
		ID:            selectedCategory.ID,
		UserID:        prevUser.ID.String(),
		AssetIDs:      []string{selectedInventory.ID},
		Collaborators: []string{prevUser.ID.String(), prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(draftAssets)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	// retrieve inventory list
	req = httptest.NewRequest(http.MethodPost, "/api/v1/category/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	AddItemsInCategory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var categoryItemResponse []model.CategoryItemResponse
	err = json.Unmarshal(data, &categoryItemResponse)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, categoryItemResponse[0].CategoryID, selectedCategory.ID)
	assert.Equal(t, categoryItemResponse[0].ItemID, selectedInventory.ID)
	assert.Equal(t, len(categoryItemResponse[0].SharableGroups), 2)
	assert.Equal(t, categoryItemResponse[0].SharableGroups[0], prevUser.ID.String())
	assert.Equal(t, categoryItemResponse[0].SharableGroups[1], prevUser.ID.String())

	// cleanup
	db.RemoveCategory(config.CTO_USER, categoryItemResponse[0].CategoryID)
}

func Test_AddItemsInCategory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddItemsInCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddItemsInCategory_InvalidDBUser(t *testing.T) {

	draftCategory := model.Category{}
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	AddItemsInCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveAssociationFromCategory(t *testing.T) {

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	draftCategory := model.Category{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "general",
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateCategory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedCategory model.Category
	err = json.Unmarshal(data, &selectedCategory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedCategory.Name, "Kitty litter box")
	assert.Equal(t, selectedCategory.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedCategory.StatusName, "general")
	assert.Equal(t, len(selectedCategory.SharableGroups), 1)
	assert.Equal(t, selectedCategory.SharableGroups[0], prevUser.ID.String())

	// retrieve inventory list
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	GetAllInventories(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var inventories []model.Inventory
	err = json.Unmarshal(data, &inventories)
	if err != nil {
		t.Errorf("expected error to be nil but got %+v", err)
	}

	selectedInventory := inventories[0]

	draftAssets := model.CategoryItemRequest{
		ID:            selectedCategory.ID,
		UserID:        prevUser.ID.String(),
		AssetIDs:      []string{selectedInventory.ID},
		Collaborators: []string{prevUser.ID.String(), prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(draftAssets)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	// retrieve inventory list
	req = httptest.NewRequest(http.MethodPost, "/api/v1/category/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	AddItemsInCategory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var categoryItemResponse []model.CategoryItemResponse
	err = json.Unmarshal(data, &categoryItemResponse)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, categoryItemResponse[0].CategoryID, selectedCategory.ID)
	assert.Equal(t, categoryItemResponse[0].ItemID, selectedInventory.ID)
	assert.Equal(t, len(categoryItemResponse[0].SharableGroups), 2)
	assert.Equal(t, categoryItemResponse[0].SharableGroups[0], prevUser.ID.String())
	assert.Equal(t, categoryItemResponse[0].SharableGroups[1], prevUser.ID.String())

	// Remove items from category

	draftAssetToRemove := model.CategoryItemRequest{
		ID:            selectedCategory.ID,
		UserID:        prevUser.ID.String(),
		AssetIDs:      []string{selectedInventory.ID},
		Collaborators: []string{prevUser.ID.String(), prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(draftAssetToRemove)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/category/remove/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	RemoveAssociationFromCategory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	// cleanup
	db.RemoveCategory(config.CTO_USER, categoryItemResponse[0].CategoryID)
}

func Test_RemoveAssociationFromCategory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveAssociationFromCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveAssociationFromCategory_InvalidDBUser(t *testing.T) {

	draftCategory := model.Category{}
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	RemoveAssociationFromCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateCategory(t *testing.T) {

	// profile are automatically derieved from the auth table
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	config.PreloadAllTestVariables()
	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllCategories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var categories []model.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	selectedCategory := categories[0]
	assert.Greater(t, len(categories), 1)

	selectedCategory.Name = "Updated Title"
	selectedCategory.Description = "Update title description"
	selectedCategory.Status = "completed" // change status so id is not passed to service

	// Marshal the draftUpdateEvent into JSON bytes
	updateCategoryRequest, err := json.Marshal(selectedCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/categories/%s", selectedCategory.ID), bytes.NewBuffer(updateCategoryRequest))
	req = mux.SetURLVars(req, map[string]string{"id": selectedCategory.ID})
	w = httptest.NewRecorder()
	UpdateCategory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var updatedCategory model.Category
	err = json.Unmarshal(data, &updatedCategory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, "Updated Title", updatedCategory.Name)

	// cleanup
	db.RemoveCategory(config.CTO_USER, updatedCategory.ID)
}

func Test_UpdateCategory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories/%s", ""), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateCategory_InvalidDBUser(t *testing.T) {

	draftCategory := model.Category{}
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	UpdateCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveCategory(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/category/0802c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "200 OK", res.Status)
}

func Test_RemoveCategory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveCategory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveCategory_InvalidDBUser(t *testing.T) {

	draftCategory := model.Category{}
	requestBody, err := json.Marshal(draftCategory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	RemoveCategory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

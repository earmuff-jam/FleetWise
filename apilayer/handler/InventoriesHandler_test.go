package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllInventories(t *testing.T) {

	// profile are automatically derieved from the auth table. due to this, we attempt to create a new user
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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetAllInventories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
}

func Test_GetAllInventories_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllInventories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var foundInventories []model.Inventory
	err = json.Unmarshal(data, &foundInventories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 0, len(foundInventories)) // empty inventory list
}

func Test_GetAllInventories_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllInventories(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllInventories_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllInventories(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllInventories_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllInventories(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetInventoryByID(t *testing.T) {

	// profile are automatically derieved from the auth table. due to this, we attempt to create a new user
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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetAllInventories(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var inventories []model.Inventory
	err = json.Unmarshal(data, &inventories)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.GreaterOrEqual(t, len(inventories), 1)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String(), "invID": inventories[0].ID})
	w = httptest.NewRecorder()
	GetInventoryByID(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
}

func Test_GetInventoryByID_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetInventoryByID(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetInventoryByID_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetInventoryByID(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetInventoryByID_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetInventoryByID(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddInventoryInBulk(t *testing.T) {

	// profile are automatically derieved from the auth table. due to this, we attempt to create a new user
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

	draftInventory := model.RawInventory{
		Name:        "Alexandro Kitteyy Litter",
		Description: "Kitty litter for a pro name game",
		Price:       23.99,
		Barcode:     "1231231231",
		SKU:         "1231231231",
		Quantity:    12,
		Location:    "Broom Closet",
		BoughtAt:    "Walmart",
	}

	inventoryMap := make(map[string]model.RawInventory)
	inventoryMap["0"] = draftInventory

	requestBody, err := json.Marshal(inventoryMap)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories/bulk", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddInventoryInBulk(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory []model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// cleanup
	removeInventory := []string{selectedInventory[0].ID}
	db.DeleteInventory(config.CTO_USER, prevUser.ID.String(), removeInventory)
}

func Test_AddInventoryInBulk_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddInventoryInBulk(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddInventoryInBulk_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddInventoryInBulk(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddNewInventory(t *testing.T) {

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

	draftInventory := model.Inventory{
		Name:           "Alexandro Kitteyy Litter",
		Description:    "Kitty litter for a pro name game",
		Price:          23.99,
		Status:         "HIDDEN",
		Barcode:        "1231231231",
		SKU:            "1231231231",
		Quantity:       12,
		IsReturnable:   true,
		ReturnLocation: "Target",
		MaxWeight:      "12",
		MinWeight:      "120",
		MaxHeight:      "24",
		MinHeight:      "12",
		Location:       "Broom Closet",
		CreatedAt:      time.Now(),
		CreatedBy:      prevUser.ID.String(),
		BoughtAt:       "Walmart",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)

	// cleanup
	removeInventory := []string{selectedInventory.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventory.ID, removeInventory)
}

func Test_AddNewInventory_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddNewInventory_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddNewInventory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddNewInventory_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddNewInventory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddNewInventory_Return_Notes(t *testing.T) {

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

	draftInventoryWithReturnNotes := model.Inventory{
		Name:           "Alexandro Kitteyy Litter",
		Description:    "Kitty litter for a pro name game",
		Price:          23.99,
		Status:         "HIDDEN",
		Barcode:        "1231231231",
		SKU:            "1231231231",
		Quantity:       12,
		IsReturnable:   true,
		ReturnLocation: "Target",
		ReturnNotes:    "print out the return label",
		MaxWeight:      "12",
		MinWeight:      "120",
		MaxHeight:      "24",
		MinHeight:      "12",
		Location:       "Broom Closet",
		CreatedAt:      time.Now(),
		CreatedBy:      prevUser.ID.String(),
		BoughtAt:       "Walmart",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftInventoryWithReturnNotes)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventoryWithReturnNotes model.Inventory
	err = json.Unmarshal(data, &selectedInventoryWithReturnNotes)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventoryWithReturnNotes.Name)
	assert.Equal(t, "Broom Closet", selectedInventoryWithReturnNotes.Location)

	// cleanup
	removeInventory := []string{selectedInventoryWithReturnNotes.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventoryWithReturnNotes.ID, removeInventory)
}

func Test_UpdateSelectedInventory(t *testing.T) {
	// profile are automatically derieved from the auth table. due to this, we attempt to create a new user
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

	draftInventory := model.Inventory{
		Name:        "Alexandro Kitteyy Litter",
		Description: "Kitty litter for a pro name game",
		Price:       23.99,
		Status:      "HIDDEN",
		Barcode:     "1231231231",
		SKU:         "1231231231",
		Quantity:    12,
		Location:    "Broom Closet",
		CreatedAt:   time.Now(),
		CreatedBy:   prevUser.ID.String(),
		BoughtAt:    "Walmart",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)

	// test to perform update but with consistent storage location name
	// TODO: write another test that checks two different inserts into storage locations as
	// this only tests updating the found name.

	selectedInventory.Name = "Alexandro Kitty Litter"

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(selectedInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/profile/%s/inventories", prevUser.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	UpdateSelectedInventory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var updatedInventory model.Inventory
	err = json.Unmarshal(data, &updatedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitty Litter", updatedInventory.Name)
	assert.Equal(t, "Broom Closet", updatedInventory.Location)

	// cleanup fn
	removeInventory := []string{selectedInventory.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventory.ID, removeInventory)
}

func Test_UpdateSelectedInventory_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateSelectedInventory_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateSelectedInventory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateSelectedInventory_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateSelectedInventory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateAssetColumn_Quantity(t *testing.T) {

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

	draftInventory := model.Inventory{
		Name:           "Alexandro Kitteyy Litter",
		Description:    "Kitty litter for a pro name game",
		Price:          23.99,
		Status:         "HIDDEN",
		Barcode:        "1231231231",
		SKU:            "1231231231",
		Quantity:       12,
		IsReturnable:   true,
		ReturnLocation: "Target",
		MaxWeight:      "12",
		MinWeight:      "120",
		MaxHeight:      "24",
		MinHeight:      "12",
		Location:       "Broom Closet",
		CreatedAt:      time.Now(),
		CreatedBy:      prevUser.ID.String(),
		BoughtAt:       "Walmart",
	}

	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)
	assert.Equal(t, 12, selectedInventory.Quantity)

	draftUpdateInv := model.UpdateAssetColumn{
		AssetID:     selectedInventory.ID,
		ColumnName:  "quantity",
		InputColumn: "1",
	}

	requestBody, err = json.Marshal(draftUpdateInv)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/profile/%s/inventories/%s", draftUserCredentials.ID.String(), selectedInventory.ID), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String(), "asssetID": selectedInventory.ID})
	w = httptest.NewRecorder()
	UpdateAssetColumn(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var updatedSelectedInv model.Inventory
	err = json.Unmarshal(data, &updatedSelectedInv)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", updatedSelectedInv.Name)
	assert.Equal(t, "Broom Closet", updatedSelectedInv.Location)
	assert.Equal(t, 1, updatedSelectedInv.Quantity)

	// cleanup
	removeInventory := []string{selectedInventory.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventory.ID, removeInventory)
}

func Test_UpdateAssetColumn_Cost(t *testing.T) {

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

	draftInventory := model.Inventory{
		Name:           "Alexandro Kitteyy Litter",
		Description:    "Kitty litter for a pro name game",
		Price:          23.99,
		Status:         "HIDDEN",
		Barcode:        "1231231231",
		SKU:            "1231231231",
		Quantity:       12,
		IsReturnable:   true,
		ReturnLocation: "Target",
		MaxWeight:      "12",
		MinWeight:      "120",
		MaxHeight:      "24",
		MinHeight:      "12",
		Location:       "Broom Closet",
		CreatedAt:      time.Now(),
		CreatedBy:      prevUser.ID.String(),
		BoughtAt:       "Walmart",
	}

	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)
	assert.Equal(t, 12, selectedInventory.Quantity)

	draftUpdateInv := model.UpdateAssetColumn{
		AssetID:     selectedInventory.ID,
		ColumnName:  "price",
		InputColumn: "99.81",
	}

	requestBody, err = json.Marshal(draftUpdateInv)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/profile/%s/inventories/%s", draftUserCredentials.ID.String(), selectedInventory.ID), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String(), "asssetID": selectedInventory.ID})
	w = httptest.NewRecorder()
	UpdateAssetColumn(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var updatedSelectedInv model.Inventory
	err = json.Unmarshal(data, &updatedSelectedInv)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", updatedSelectedInv.Name)
	assert.Equal(t, "Broom Closet", updatedSelectedInv.Location)
	assert.Equal(t, 99.81, updatedSelectedInv.Price)

	// cleanup
	removeInventory := []string{selectedInventory.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventory.ID, removeInventory)
}

func Test_UpdateAssetColumn_InvalidColumn(t *testing.T) {

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

	draftInventory := model.Inventory{
		Name:           "Alexandro Kitteyy Litter",
		Description:    "Kitty litter for a pro name game",
		Price:          23.99,
		Status:         "HIDDEN",
		Barcode:        "1231231231",
		SKU:            "1231231231",
		Quantity:       12,
		IsReturnable:   true,
		ReturnLocation: "Target",
		MaxWeight:      "12",
		MinWeight:      "120",
		MaxHeight:      "24",
		MinHeight:      "12",
		Location:       "Broom Closet",
		CreatedAt:      time.Now(),
		CreatedBy:      prevUser.ID.String(),
		BoughtAt:       "Walmart",
	}

	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)
	assert.Equal(t, 12, selectedInventory.Quantity)

	draftUpdateInv := model.UpdateAssetColumn{
		AssetID:     selectedInventory.ID,
		ColumnName:  "name",
		InputColumn: "new name;",
	}

	requestBody, err = json.Marshal(draftUpdateInv)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/profile/%s/inventories/%s", draftUserCredentials.ID.String(), selectedInventory.ID), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String(), "asssetID": selectedInventory.ID})
	w = httptest.NewRecorder()
	UpdateAssetColumn(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)

	// cleanup
	removeInventory := []string{selectedInventory.ID}
	db.DeleteInventory(config.CTO_USER, selectedInventory.ID, removeInventory)
}

func Test_UpdateAssetColumn_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/0802c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":      "0802c692-b8e2-4824-a870-e52f4a0cccf8",
		"assetID": "0802c692-b8e2-4824-a870-e52f4a0cccf8"},
	)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateAssetColumn(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateAssetColumn_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/0802c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":      "",
		"assetID": ""},
	)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateAssetColumn(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveSelectedInventory(t *testing.T) {
	// profile are automatically derieved from the auth table. due to this, we attempt to create a new user
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

	draftInventory := model.Inventory{
		Name:        "Alexandro Kitteyy Litter",
		Description: "Kitty litter for a pro name game",
		Price:       23.99,
		Status:      "HIDDEN",
		Barcode:     "1231231231",
		SKU:         "1231231231",
		Quantity:    12,
		Location:    "Broom Closet",
		CreatedAt:   time.Now(),
		CreatedBy:   prevUser.ID.String(),
		BoughtAt:    "Walmart",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w := httptest.NewRecorder()
	AddNewInventory(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedInventory model.Inventory
	err = json.Unmarshal(data, &selectedInventory)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "Alexandro Kitteyy Litter", selectedInventory.Name)
	assert.Equal(t, "Broom Closet", selectedInventory.Location)

	removeInventory := map[string]string{"0": selectedInventory.ID}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(removeInventory)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/inventories/prune", draftUserCredentials.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": draftUserCredentials.ID.String()})
	w = httptest.NewRecorder()
	RemoveSelectedInventory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
}

func Test_RemoveSelectedInventory_WrongUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/prune", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveSelectedInventory_NoUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/prune", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveSelectedInventory_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/prune", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveSelectedInventory(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveSelectedInventory_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/0802c692-b8e2-4824-a870-e52f4a0cccf8/inventories/prune", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveSelectedInventory(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

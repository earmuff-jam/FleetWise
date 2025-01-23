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

func Test_GetAllMaintenancePlans(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllMaintenancePlans(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var maintenancePlan []model.MaintenancePlan
	err = json.Unmarshal(data, &maintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(maintenancePlan), 0)
}

func Test_GetAllMaintenancePlans_IncorrectPlanID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", "", 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllMaintenancePlans(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllMaintenancePlans_InvalidDBUser(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllMaintenancePlans(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetMaintenancePlan(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllMaintenancePlans(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var maintenancePlan []model.MaintenancePlan
	err = json.Unmarshal(data, &maintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(maintenancePlan), 0)

	selectedMaintenancePlan := maintenancePlan[0]
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/plan?id=%s&mID=%s", prevUser.ID.String(), selectedMaintenancePlan.ID), nil)
	w = httptest.NewRecorder()
	GetMaintenancePlan(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var selected model.MaintenancePlan
	err = json.Unmarshal(data, &selected)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, selectedMaintenancePlan.ID, selected.ID)
}

func Test_GetMaintenancePlan_IncorrectPlanID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/plan?id=%s&mID=%s", "", ""), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetMaintenancePlan_InvalidDBUser(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/plan?id=%s&mID=%s", prevUser.ID.String(), ""), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_CreateMaintenancePlan(t *testing.T) {

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

	draftMaintenancePlan := model.MaintenancePlan{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "completed",
		PlanType:       "annual",
		PlanDue:        time.Now().AddDate(1, 0, 0),
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedMaintenancePlan model.MaintenancePlan
	err = json.Unmarshal(data, &selectedMaintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedMaintenancePlan.Name, "Kitty litter box")
	assert.Equal(t, selectedMaintenancePlan.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedMaintenancePlan.StatusName, "completed")
	assert.Equal(t, selectedMaintenancePlan.PlanType, "annual")

	// cleanup
	db.RemoveMaintenancePlan(config.CTO_USER, selectedMaintenancePlan.ID)
}

func Test_CreateMaintenancePlan_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	CreateMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_CreateMaintenancePlan_InvalidDBUser(t *testing.T) {

	draftMaintenancePlan := model.MaintenancePlan{}
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateMaintenancePlan(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllMaintenancePlans(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var maintenancePlans []model.MaintenancePlan
	err = json.Unmarshal(data, &maintenancePlans)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	selectedMaintenancePlan := maintenancePlans[0]
	assert.Greater(t, len(maintenancePlans), 1)

	selectedMaintenancePlan.Name = "Updated Title"
	selectedMaintenancePlan.Description = "Update title description"
	selectedMaintenancePlan.Status = "completed" // change status so id is not passed to service

	// Marshal the draftUpdateEvent into JSON bytes
	updateMaintenancePlanRequest, err := json.Marshal(selectedMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/plan/%s", selectedMaintenancePlan.ID), bytes.NewBuffer(updateMaintenancePlanRequest))
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	UpdateMaintenancePlan(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, res.StatusCode, 200)
	assert.Greater(t, len(data), 0)

	var updatedMaintenancePlan model.MaintenancePlan
	err = json.Unmarshal(data, &updatedMaintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, updatedMaintenancePlan.Name, "Updated Title")
	assert.Equal(t, updatedMaintenancePlan.StatusName, "completed")

	// cleanup
	db.RemoveMaintenancePlan(config.CTO_USER, updatedMaintenancePlan.ID)
}

func Test_UpdateMaintenancePlan_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	UpdateMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateMaintenancePlan_InvalidDBUser(t *testing.T) {

	draftMaintenancePlan := model.MaintenancePlan{}
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	UpdateMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveMaintenancePlan(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/plan/0802c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "200 OK", res.Status)
}

func Test_RemoveMaintenancePlan_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan/0802c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveMaintenancePlan_InvalidDBUser(t *testing.T) {

	draftMaintenancePlan := model.MaintenancePlan{}
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	RemoveMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllMaintenancePlanItems(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/maintenance-plans?id=%s&limit=%d", prevUser.ID.String(), 5), nil)
	w := httptest.NewRecorder()
	GetAllMaintenancePlans(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var maintenancePlan []model.MaintenancePlan
	err = json.Unmarshal(data, &maintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(maintenancePlan), 0)

	selectedMaintenancePlan := maintenancePlan[0]
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/plans/items?id=%s&limit=%d&mID=%s", prevUser.ID.String(), 5, selectedMaintenancePlan.ID), nil)
	w = httptest.NewRecorder()
	GetAllMaintenancePlanItems(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var selected model.MaintenanceItemResponse
	err = json.Unmarshal(data, &selected)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
}

func Test_GetAllMaintenancePlanItems_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/plans/items?id=%s&limit=%d&mID=%s", "", 5, ""), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllMaintenancePlanItems(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAllMaintenancePlanItems_InvalidDBUser(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/plans/items?id=%s&limit=%d&mID=%s", "", 5, ""), nil)
	w := httptest.NewRecorder()
	GetAllMaintenancePlanItems(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddItemsInMaintenancePlan(t *testing.T) {

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

	draftMaintenancePlan := model.MaintenancePlan{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "completed",
		PlanType:       "annual",
		PlanDue:        time.Now().AddDate(1, 0, 0),
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedMaintenancePlan model.MaintenancePlan
	err = json.Unmarshal(data, &selectedMaintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedMaintenancePlan.Name, "Kitty litter box")
	assert.Equal(t, selectedMaintenancePlan.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedMaintenancePlan.StatusName, "completed")
	assert.Equal(t, selectedMaintenancePlan.PlanType, "annual")
	assert.Equal(t, len(selectedMaintenancePlan.SharableGroups), 1)
	assert.Equal(t, selectedMaintenancePlan.SharableGroups[0], prevUser.ID.String())

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

	draftAssets := model.MaintenanceItemRequest{
		ID:            selectedMaintenancePlan.ID,
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
	req = httptest.NewRequest(http.MethodPost, "/api/v1/plans/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	AddItemsInMaintenancePlan(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var maintenanceItemResponse []model.MaintenanceItemResponse
	err = json.Unmarshal(data, &maintenanceItemResponse)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, maintenanceItemResponse[0].MaintenancePlanID, selectedMaintenancePlan.ID)
	assert.Equal(t, maintenanceItemResponse[0].ItemID, selectedInventory.ID)
	assert.Equal(t, len(maintenanceItemResponse[0].SharableGroups), 2)
	assert.Equal(t, maintenanceItemResponse[0].SharableGroups[0], prevUser.ID.String())
	assert.Equal(t, maintenanceItemResponse[0].SharableGroups[1], prevUser.ID.String())

	// cleanup
	db.RemoveMaintenancePlan(config.CTO_USER, maintenanceItemResponse[0].MaintenancePlanID)
}

func Test_AddItemsInMaintenancePlan_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/items", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	AddItemsInMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_AddItemsInMaintenancePlan_InvalidDBUser(t *testing.T) {

	draftMaintenanceItemReq := model.MaintenanceItemRequest{}
	requestBody, err := json.Marshal(draftMaintenanceItemReq)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/items", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	AddItemsInMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveAssociationFromMaintenancePlan(t *testing.T) {

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

	draftMaintenancePlan := model.MaintenancePlan{
		Name:           "Kitty litter box",
		Description:    "Palette storage of kitty litter",
		Color:          "#f7f7f7",
		Status:         "completed",
		PlanType:       "annual",
		PlanDue:        time.Now().AddDate(1, 0, 0),
		CreatedBy:      prevUser.ID.String(),
		UpdatedBy:      prevUser.ID.String(),
		SharableGroups: []string{prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftMaintenancePlan)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	CreateMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)

	var selectedMaintenancePlan model.MaintenancePlan
	err = json.Unmarshal(data, &selectedMaintenancePlan)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, selectedMaintenancePlan.Name, "Kitty litter box")
	assert.Equal(t, selectedMaintenancePlan.Description, "Palette storage of kitty litter")
	assert.Equal(t, selectedMaintenancePlan.StatusName, "completed")
	assert.Equal(t, selectedMaintenancePlan.PlanType, "annual")
	assert.Equal(t, len(selectedMaintenancePlan.SharableGroups), 1)
	assert.Equal(t, selectedMaintenancePlan.SharableGroups[0], prevUser.ID.String())

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

	draftAssets := model.MaintenanceItemRequest{
		ID:            selectedMaintenancePlan.ID,
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
	req = httptest.NewRequest(http.MethodPost, "/api/v1/plans/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	AddItemsInMaintenancePlan(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	var maintenanceItemResponse []model.MaintenanceItemResponse
	err = json.Unmarshal(data, &maintenanceItemResponse)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, maintenanceItemResponse[0].MaintenancePlanID, selectedMaintenancePlan.ID)
	assert.Equal(t, maintenanceItemResponse[0].ItemID, selectedInventory.ID)
	assert.Equal(t, len(maintenanceItemResponse[0].SharableGroups), 2)
	assert.Equal(t, maintenanceItemResponse[0].SharableGroups[0], prevUser.ID.String())
	assert.Equal(t, maintenanceItemResponse[0].SharableGroups[1], prevUser.ID.String())

	// Remove items from category

	draftAssetToRemove := model.MaintenanceItemRequest{
		ID:            selectedMaintenancePlan.ID,
		UserID:        prevUser.ID.String(),
		AssetIDs:      []string{selectedInventory.ID},
		Collaborators: []string{prevUser.ID.String(), prevUser.ID.String()},
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(draftAssetToRemove)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/plan/remove/items", bytes.NewBuffer(requestBody))
	w = httptest.NewRecorder()
	RemoveAssociationFromCategory(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	// cleanup
	db.RemoveMaintenancePlan(config.CTO_USER, maintenanceItemResponse[0].MaintenancePlanID)
}

func Test_RemoveAssociationFromMaintenancePlan_IncorrectUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan/remove/items", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "request"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	RemoveAssociationFromMaintenancePlan(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_RemoveAssociationFromMaintenancePlan_InvalidDBUser(t *testing.T) {

	draftMaintenanceItemRequest := model.MaintenanceItemRequest{}
	requestBody, err := json.Marshal(draftMaintenanceItemRequest)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plan/remove/items", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	RemoveAssociationFromMaintenancePlan(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

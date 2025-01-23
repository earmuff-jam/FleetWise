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

func Test_GetAllUserProfiles(t *testing.T) {

	config.PreloadAllTestVariables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/list", nil)
	w := httptest.NewRecorder()
	GetAllUserProfiles(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var currentUsers []model.Profile
	err = json.Unmarshal(data, &currentUsers)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(currentUsers), 1)
	assert.Equal(t, "john", currentUsers[0].Username)
}

func Test_GetAllUserProfiles_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/event/0902c692-b8e2-4824-a870-e52f4a0cccf8/shared", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllUserProfiles(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetProfileStats(t *testing.T) {

	config.PreloadAllTestVariables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0902c692-b8e2-4824-a870-e52f4a0cccf8/stats", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0902c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	GetProfileStats(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var profileStats model.ProfileStats
	err = json.Unmarshal(data, &profileStats)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
}

func Test_GetProfileStats_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0902c692-b8e2-4824-a870-e52f4a0cccf8/stats", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetProfileStats(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetNotifications_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/0902c692-b8e2-4824-a870-e52f4a0cccf8/notifications", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetNotifications(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetProfileApi(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetProfile(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	t.Logf("response = %+v", string(data))

	w = httptest.NewRecorder()
	GetProfile(w, req, config.CEO_USER)
	res = w.Result()
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetUsernameApi(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/username", prevUser.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetUsername(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var userName string
	err = json.Unmarshal(data, &userName)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Equal(t, "john", userName)
	t.Logf("response = %+v", string(data))

	w = httptest.NewRecorder()
	GetUsername(w, req, config.CEO_USER)
	res = w.Result()
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_UpdateProfileApi(t *testing.T) {

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

	draftUserProfile := model.Profile{
		Username: "Unit test user",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserProfile)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/profile/%s", prevUser.ID), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	UpdateProfile(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var newProfile model.Profile
	err = json.Unmarshal(data, &newProfile)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.True(t, len(string(data)) > 0)
	assert.Equal(t, newProfile.Username, draftUserProfile.Username)

	// cleanup
	cleanUpProfile := model.Profile{
		Username:     "john",
		FullName:     "John Doe",
		PhoneNumber:  "1234567890",
		AboutMe:      "I like to climb trees and hike with my friends",
		EmailAddress: "admin@gmail.com",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err = json.Marshal(cleanUpProfile)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s", prevUser.ID), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	UpdateProfile(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var prevProfile model.Profile
	err = json.Unmarshal(data, &prevProfile)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, "john", prevProfile.Username)

	w = httptest.NewRecorder()
	UpdateProfile(w, req, config.CEO_USER)
	res = w.Result()
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetFavouriteItems(t *testing.T) {
	config.PreloadAllTestVariables()

	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	prevUser, err := db.RetrieveUser(config.CTO_USER, &draftUserCredentials)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/profile/%s/fav", prevUser.ID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetFavouriteItems(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var favItems []model.FavouriteItem
	err = json.Unmarshal(data, &favItems)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
}

func Test_SaveFovouriteItems_Category(t *testing.T) {
	config.PreloadAllTestVariables()

	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

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

	draftFavouriteItems := model.FavouriteItem{
		CategoryID: selectedCategory.ID,
		CreatedBy:  prevUser.ID.String(),
		UpdatedBy:  prevUser.ID.String(),
	}

	requestBody, err = json.Marshal(draftFavouriteItems)
	if err != nil {
		t.Errorf("failed to marshall request body. error: %+v", err)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/profile/%s/fav", prevUser.ID.String()), bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w = httptest.NewRecorder()
	SaveFavItem(w, req, config.CTO_USER)
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var favItems []model.FavouriteItem
	err = json.Unmarshal(data, &favItems)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(favItems), 1)

	// cleanup
	db.RemoveCategory(config.CTO_USER, selectedCategory.ID)

}

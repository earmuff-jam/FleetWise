package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohit2530/communityCare/config"
	"github.com/mohit2530/communityCare/db"
	"github.com/mohit2530/communityCare/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetAssetsAndSummary(t *testing.T) {

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

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/summary?id=%s", prevUser.ID.String()), nil)
	w := httptest.NewRecorder()
	GetAssetsAndSummary(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var asResponse model.AssetsAndSummaryResponse
	err = json.Unmarshal(data, &asResponse)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
	assert.Greater(t, len(asResponse.AssetList), 0)
	assert.Greater(t, len(asResponse.AssetSummaryList), 0)
}

func Test_GetAssetsAndSummary_IncorrectCategoryID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories?id=%s&limit=%d", "", 5), nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAssetsAndSummary(w, req, config.CTO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetAssetsAndSummary_InvalidDBUser(t *testing.T) {

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
	GetAssetsAndSummary(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

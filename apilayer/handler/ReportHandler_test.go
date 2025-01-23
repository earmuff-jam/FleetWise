package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mohit2530/communityCare/config"
	"github.com/mohit2530/communityCare/db"
	"github.com/mohit2530/communityCare/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetReports(t *testing.T) {
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
	sinceTimeRange := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/reports/%s?since=%s&includeOverdue=%s", prevUser.ID, sinceTimeRange, "false"), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetReports(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var reports []model.Report
	err = json.Unmarshal(data, &reports)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	parsedTime, _ := time.Parse(time.RFC3339, sinceTimeRange)

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(reports), 1)
	assert.Equal(t, 108.00, reports[0].ItemValuation)
	assert.Equal(t, 0.00, reports[0].TotalCategoryItemsCost)
	assert.Equal(t, parsedTime, reports[0].SelectedTimeRange)
}

func Test_GetReports_UntilDate_IncludeOverdue(t *testing.T) {
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
	sinceTimeRange := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/reports/%s?since=%s&includeOverdue=%s", prevUser.ID, sinceTimeRange, "true"), nil)
	req = mux.SetURLVars(req, map[string]string{"id": prevUser.ID.String()})
	w := httptest.NewRecorder()
	GetReports(w, req, config.CTO_USER)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var reports []model.Report
	err = json.Unmarshal(data, &reports)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	parsedTime, _ := time.Parse(time.RFC3339, sinceTimeRange)

	assert.Equal(t, 200, res.StatusCode)
	assert.GreaterOrEqual(t, len(reports), 1)
	assert.Equal(t, 108.00, reports[0].ItemValuation)
	assert.Equal(t, 0.00, reports[0].TotalCategoryItemsCost)
	assert.Equal(t, parsedTime, reports[0].SelectedTimeRange)
}

func Test_GetReports_InvalidDBUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/0902c692-b8e2-4824-a870-e52f4a0cccf8", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0802c692-b8e2-4824-a870-e52f4a0cccf8"})
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetReports(w, req, config.CEO_USER)
	res := w.Result()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

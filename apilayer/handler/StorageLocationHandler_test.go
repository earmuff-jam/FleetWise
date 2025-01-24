package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllStorageLocations(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/locations", nil)
	w := httptest.NewRecorder()
	config.PreloadAllTestVariables()
	GetAllStorageLocations(w, req, config.CTO_USER)
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
	GetAllStorageLocations(w, req, config.CEO_USER)
	res = w.Result()
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

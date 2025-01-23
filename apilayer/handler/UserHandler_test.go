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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetSignInApi_Success(t *testing.T) {

	t.Skip()
	// profiles are derieved from sign in
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserCredentials)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/sigin", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	Signin(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)
}

func Test_GetSignInApi_Failure(t *testing.T) {
	t.Skip()
	draftUserCredentials := model.UserCredentials{}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserCredentials)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/sigin", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	Signin(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetSignUpApi_Success(t *testing.T) {
	t.Skip()
	id := uuid.New()

	// profiles are derieved from sign in
	draftUserCredentials := model.UserCredentials{
		Email:             fmt.Sprintf("test%s@gmail.com", id),
		Username:          id.String(),
		Role:              "TESTER",
		EncryptedPassword: id.String(),
		Birthday:          "2009-01-01",
	}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserCredentials)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}

	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/signup", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	Signup(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var user model.UserCredentials
	err = json.Unmarshal(data, &user)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Greater(t, len(data), 0)

	db.RemoveUser(config.CTO_USER, user.ID)
}

func Test_GetSignUpApi_Failure(t *testing.T) {
	t.Skip()

	draftUserCredentials := model.UserCredentials{}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserCredentials)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/signup", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	Signup(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "400 Bad Request", res.Status)
}

func Test_GetLogout(t *testing.T) {
	t.Skip()

	draftUserCredentials := model.UserCredentials{}

	// Marshal the draftEvent into JSON bytes
	requestBody, err := json.Marshal(draftUserCredentials)
	if err != nil {
		t.Errorf("failed to marshal JSON: %v", err)
	}
	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/logout", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	Logout(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "200 OK", res.Status)
}

func Test_IsValidUserEmail_Success_Valid_Email(t *testing.T) {
	t.Skip()
	draftUserEmail := model.UserResponse{
		EmailAddress: "admin@gmail.com",
	}
	reqBody, err := json.Marshal(draftUserEmail)
	if err != nil {
		t.Errorf("failed to marshall json. error: %+v", err)
	}
	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/isValidEmail", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	IsValidUserEmail(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var truthyValue bool
	err = json.Unmarshal(data, &truthyValue)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, false, truthyValue) // false because api responds false if the email is found
}

func Test_IsValidUserEmail_Success_Invalid_Email(t *testing.T) {
	t.Skip()
	draftUserEmail := model.UserResponse{
		EmailAddress: "admin23@gmail.com",
	}
	reqBody, err := json.Marshal(draftUserEmail)
	if err != nil {
		t.Errorf("failed to marshall json. error: %+v", err)
	}
	config.PreloadAllTestVariables()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/isValidEmail", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	IsValidUserEmail(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var truthyValue bool
	err = json.Unmarshal(data, &truthyValue)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, true, truthyValue) // true because api responds false if the email is found
}

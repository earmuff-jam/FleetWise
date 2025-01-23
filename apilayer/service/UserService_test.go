package service

import (
	"os"
	"testing"

	"github.com/mohit2530/communityCare/config"
	"github.com/mohit2530/communityCare/db"
	"github.com/mohit2530/communityCare/model"
	"github.com/stretchr/testify/assert"
)

func Test_FetchUser(t *testing.T) {

	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	resp, err := FetchUser(config.CTO_USER, &draftUserCredentials)

	assert.NoError(t, err)
	assert.Equal(t, resp.EmailAddress, "admin@gmail.com")
}

func Test_FetchUser_DefaultTokenValidityTime(t *testing.T) {

	os.Setenv("TOKEN_VALIDITY_TIME", "")
	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	resp, err := FetchUser(config.CTO_USER, &draftUserCredentials)

	assert.NoError(t, err)
	assert.Equal(t, resp.EmailAddress, "admin@gmail.com")
}

func Test_FetchUser_InvalidUser(t *testing.T) {

	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	_, err := FetchUser(config.CEO_USER, &draftUserCredentials)
	assert.Error(t, err)

}

func Test_RegisterUser(t *testing.T) {

	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "test_admin_user@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
		Username:          "testUsername",
	}

	resp, err := RegisterUser(config.CTO_USER, &draftUserCredentials)

	assert.NoError(t, err)
	assert.Equal(t, resp.Email, "test_admin_user@gmail.com")

	db.RemoveUser(config.CTO_USER, resp.ID)

}

func Test_RegisterUser_InvalidUser(t *testing.T) {

	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "test_admin_user@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
		Username:          "testUsername",
	}

	_, err := RegisterUser(config.CEO_USER, &draftUserCredentials)
	assert.Error(t, err)

}

func Test_ValidateCredentials(t *testing.T) {
	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	resp, err := FetchUser(config.CTO_USER, &draftUserCredentials)

	assert.NoError(t, err)
	assert.Equal(t, resp.EmailAddress, "admin@gmail.com")

}

func Test_ValidateCredentials_InvalidUser(t *testing.T) {
	config.PreloadAllTestVariables()

	// retrieve the selected profile
	draftUserCredentials := model.UserCredentials{
		Email:             "admin@gmail.com",
		Role:              "TESTER",
		EncryptedPassword: "1231231",
	}

	resp, err := FetchUser(config.CTO_USER, &draftUserCredentials)

	assert.NoError(t, err)
	assert.Equal(t, resp.EmailAddress, "admin@gmail.com")

}

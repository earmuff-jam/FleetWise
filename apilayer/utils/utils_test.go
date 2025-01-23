package utils

import (
	"testing"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/ciri-stormrider/types"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateJwtToken(t *testing.T) {

	resp, err := stormRider.CreateJWT(&types.Credentials{}, "")

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.LicenceKey), 10)
	assert.GreaterOrEqual(t, len(resp.Cookie), 20)

	err = ValidateJwtToken(resp.Cookie)
	assert.NoError(t, err)
}

func Test_ValidateJwtToken_InvalidToken(t *testing.T) {

	fakeCreds := &types.Credentials{
		Cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDB9._5Dsrto92EwgLi5kd9-gD1NqOpUu29JnhZxOwupucyc",
	}
	err := ValidateJwtToken(fakeCreds.Cookie)
	assert.Error(t, err)

}

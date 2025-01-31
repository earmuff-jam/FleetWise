package utils

import (
	"errors"
	"os"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/fleetwise/config"
)

// ValidateJwtToken ...
//
// Method is used to validate the provided JWT token
func ValidateJwtToken(token string) error {

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		secretToken = ""
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
	}

	isValid, err := stormRider.ValidateJWT(token, secretToken)
	if err != nil {
		config.Log("invalid token detected", err)
		return err
	}

	// check token validity
	if !isValid {
		config.Log("token in invalid", err)
		return errors.New("token is invalid")
	}

	return nil
}

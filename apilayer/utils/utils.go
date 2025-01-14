package utils

import (
	"errors"
	"log"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
)

// ValidateJwtToken ...
//
// Method is used to validate the provided JWT token
func ValidateJwtToken(token string) error {

	isValid, err := stormRider.ValidateJWT(token, "")
	if err != nil {
		log.Printf("invalid token detected. error :%+v", err)
		return err
	}

	// check token validity
	if !isValid {
		log.Printf("token in invalid. error: %+v", err)
		return errors.New("token is invalid")
	}

	return nil
}

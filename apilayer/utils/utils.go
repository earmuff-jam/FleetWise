package utils

import (
	"errors"
	"log"
	"os"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
)

// ValidateJwtToken ...
//
// Method is used to validate the provided JWT token
func ValidateJwtToken(token string) error {

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		log.Print("unable to retrieve secret token key. defaulting to default values")
		secretToken = ""
	}

	isValid, err := stormRider.ValidateJWT(token, secretToken)
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

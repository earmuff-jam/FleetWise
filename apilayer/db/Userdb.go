package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/ciri-stormrider/types"

	"github.com/google/uuid"
	"github.com/mohit2530/communityCare/model"
	"github.com/mohit2530/communityCare/utils"
	"golang.org/x/crypto/bcrypt"
)

// SaveUser ...
func SaveUser(user string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// generate the hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(draftUser.EncryptedPassword), 8)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	sqlStr := `
	INSERT INTO auth.users(email, username, birthdate, encrypted_password, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	var draftUserID string
	err = tx.QueryRow(
		sqlStr,
		draftUser.Email,
		draftUser.Username,
		draftUser.Birthday,
		string(hashedPassword),
		draftUser.Role,
	).Scan(&draftUserID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	draftUser.ID, err = uuid.Parse(draftUserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return draftUser, nil
}

// RetrieveUser ...
//
// Function is used to retrieve details about the selected user. The email address is the unique fieldset
// for any given selected user. JWT token is applied only after the user is verified. The ID of the selected
// user is used to prefil from the database.
func RetrieveUser(user string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {
	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to setup database. error: %+v", err)
		return nil, err
	}
	defer db.Close()

	// retrive the encrypted pwd. EMAIL must be UNIQUE field.
	sqlStr := `SELECT id, role, encrypted_password FROM auth.users WHERE email=$1;`

	result := db.QueryRow(sqlStr, draftUser.Email)

	storedCredentials := &model.UserCredentials{}
	err = result.Scan(&storedCredentials.ID, &storedCredentials.Role, &storedCredentials.EncryptedPassword)
	if err != nil {
		log.Printf("unable to retrieve user details. error: +%v", err)
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.EncryptedPassword), []byte(draftUser.EncryptedPassword)); err != nil {
		return nil, err
	}

	draftUser.ID = storedCredentials.ID
	draftUser.Role = storedCredentials.Role

	draftTime := os.Getenv("TOKEN_VALIDITY_TIME")
	if len(draftTime) <= 0 {
		draftTime = "15"
	}

	userCredsWithToken, err := stormRider.CreateJWT(&types.Credentials{}, draftTime, "")
	if err != nil {
		log.Printf("unable to create JWT token. error: %+v", err)
		return nil, err
	}
	draftUser.PreBuiltToken = userCredsWithToken.Cookie
	draftUser.LicenceKey = userCredsWithToken.LicenceKey

	updateJwtToken(user, draftUser)
	return draftUser, nil
}

// updateJwtToken ...
//
// allows to update the token schema with the proper credentials for the user
// also updates the auth.users table with the used license key to decode the jwt.
// the result however is a masked entity to preserve the users jwt
func updateJwtToken(user string, draftUser *model.UserCredentials) error {
	db, err := SetupDB(user)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = upsertLicenseKey(draftUser.Id, draftUser.LicenceKey, tx)
	if err != nil {
		log.Printf("unable to add license key. error: %+v", err)
		tx.Rollback()
		return err
	}

	err = upsertOauthToken(draftUser, tx)
	if err != nil {
		log.Printf("unable to add auth token. error: %+v", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("unable to commit selected transaction. error: %+v", err)
		return err
	}

	return nil
}

// upsertLicenseKey...
//
// set the instance_id as the license key used to encode / decode the jwt.
// save each users own key so that we can decode the token in private if needed be.
func upsertLicenseKey(userID string, licenseKey string, tx *sql.Tx) error {

	sqlStr := "UPDATE auth.users SET instance_id = $1 WHERE id = $2;"
	_, err := tx.Exec(sqlStr, licenseKey, userID)
	if err != nil {
		log.Printf("unable to add license key to signed in user. error: +%v", err)
		tx.Rollback()
		return err
	}
	return nil
}

// upsertOauthToken ...
//
// updates the oauth token table in the database
func upsertOauthToken(draftUser *model.UserCredentials, tx *sql.Tx) error {

	var maskedID string

	sqlStr := `
	INSERT INTO auth.oauth
	(token, user_id, expiration_time, user_agent)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id)
	DO UPDATE SET
		token = EXCLUDED.token,
		expiration_time = EXCLUDED.expiration_time,
		user_agent = EXCLUDED.user_agent
	RETURNING id;`

	err := tx.QueryRow(
		sqlStr,
		draftUser.PreBuiltToken,
		draftUser.ID,
		draftUser.ExpirationTime,
		draftUser.UserAgent,
	).Scan(&maskedID)

	if err != nil {
		log.Printf("unable to add token. Error: %v", err)
		tx.Rollback()
		return err
	}

	// apply the masked token
	draftUser.PreBuiltToken = maskedID
	return nil
}

// IsValidUserEmail ...
func IsValidUserEmail(user string, draftUserEmail string) (bool, error) {
	db, err := SetupDB(user)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// retrive the encrypted pwd. EMAIL must be UNIQUE field.
	sqlStr := `SELECT EXISTS(SELECT 1 FROM auth.users u WHERE u.email=$1);`

	result := db.QueryRow(sqlStr, draftUserEmail)
	exists := false
	err = result.Scan(&exists)
	if err != nil {
		log.Printf("unable to validate user email address. error: +%v", err)
		return false, err
	}
	return !exists, nil // return false if found
}

// RemoveUser ...
func RemoveUser(user string, id uuid.UUID) error {
	db, err := SetupDB(user)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM auth.users WHERE id = $1`
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		log.Printf("Error deleting user with ID %s: %v", id.String(), err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected after deleting user: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No user found with ID %s", id.String())
		// You might want to return a custom error here if needed
		return nil
	}

	return nil
}

// ValidateCredentials ...
//
// Method is used to verify if the incoming api calls have a valid jwt token.
// If the validity of the token is crossed, or if the token itself is invalid the error is propogated up the method chain.
func ValidateCredentials(user string, ID string) error {
	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to setup db. error: %+v", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to setup transaction for selected db. error: %+v", err)
		return err
	}

	var tokenFromDb string
	var expirationTime time.Time
	err = tx.QueryRow(`SELECT token, expiration_time FROM auth.oauth WHERE id=$1 LIMIT 1;`, ID).Scan(&tokenFromDb, &expirationTime)
	if err != nil {
		log.Printf("unable to retrive validated token. error: +%v", err)
		tx.Rollback()
		return err
	}

	err = utils.ValidateJwtToken(tokenFromDb)
	if err != nil {
		log.Printf("unable to validate jwt token. error: %+v", err)
		tx.Rollback()
		return err
	}

	// Check if the token is within the last 30 seconds of its expiry time
	// token is about to expire. if the user is continuing with activity, create new token
	formattedTimeToLive := time.Until(expirationTime)
	if formattedTimeToLive <= 30*time.Second && formattedTimeToLive > 0 {
		updatedToken, err := stormRider.RefreshToken("15", "")
		if err != nil {
			log.Printf("unable to refresh token. error :%+v", err)
			tx.Rollback()
			return err
		}
		parsedUserID, err := uuid.Parse(ID)
		if err != nil {
			log.Printf("unable to determine user id. error :%+v", err)
			return err
		}
		ghostUser := model.UserCredentials{
			ID:             parsedUserID,
			ExpirationTime: time.Now().Add((time.Duration(15) * time.Minute)),
			PreBuiltToken:  updatedToken,
		}
		err = revalidateOauthToken(&ghostUser, tx)
		if err != nil {
			log.Printf("unable to revalidate the user. error %+v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("unable to commit transaction. error: %+v", err)
		tx.Rollback()
		return err
	}

	return nil
}

// revalidateOauthToken ...
//
// revalidates the user token.
func revalidateOauthToken(draftUser *model.UserCredentials, tx *sql.Tx) error {

	sqlStr := `UPDATE auth.oauth SET token=$1, expiration_time=$2 WHERE id=$3;`
	_, err := tx.Exec(sqlStr, draftUser.PreBuiltToken, draftUser.ExpirationTime, draftUser.ID)

	if err != nil {
		log.Printf("unable to refresh token. Error: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}

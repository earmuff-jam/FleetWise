package service

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/ciri-stormrider/types"
	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/earmuff-jam/fleetwise/utils"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	EMAIL_SUBJECT_LINE  = "Verify your email address for FleetWise Application"
	WEB_APPLICATION_URI = "/api/v1/"
)

// FetchUser ...
//
// Function is used to retrieve user details and perform jwt maniupulation in the application
func FetchUser(user string, draftUser *model.UserCredentials) (*model.UserResponse, error) {

	draftTime := os.Getenv("TOKEN_VALIDITY_TIME")
	if len(draftTime) <= 0 {
		config.Log("unable to find token validity time. defaulting to default values", nil)
		draftTime = config.DefaultTokenValidityTime
	}

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
		secretToken = ""
	}

	draftUser, err := db.RetrieveUser(user, draftUser)
	if err != nil {
		config.Log("unable to retrieve user details", err)
		return nil, err
	}

	formattedTime, err := strconv.ParseInt(draftTime, 10, 64)
	if err != nil {
		config.Log("unable to parse provided time", err)
		return nil, err
	}

	draftCredentials := types.Credentials{
		Claims: jwt.StandardClaims{
			ExpiresAt: formattedTime,
			Subject:   draftUser.ID.String(),
		},
	}

	userCredsWithToken, err := stormRider.CreateJWT(&draftCredentials, secretToken)

	if err != nil {
		config.Log("unable to create JWT token", err)
		return nil, err
	}
	draftUser.PreBuiltToken = userCredsWithToken.Cookie
	draftUser.LicenceKey = userCredsWithToken.LicenceKey

	err = updateJwtToken(user, draftUser)
	if err != nil {
		config.Log("unable to upsert token", err)
		return nil, err
	}

	return &model.UserResponse{
		ID:           draftUser.ID.String(),
		EmailAddress: draftUser.Email,
		IsVerified:   draftUser.IsVerified,
	}, nil
}

// updateJwtToken ...
//
// allows to update the token schema with the proper credentials for the user
// also updates the auth.users table with the used license key to decode the jwt.
// the result however is a masked entity to preserve the users jwt
func updateJwtToken(user string, draftUser *model.UserCredentials) error {
	db, err := db.SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to setup transaction for db", err)
		return err
	}

	err = upsertLicenseKey(draftUser.ID.String(), draftUser.LicenceKey, tx)
	if err != nil {
		config.Log("unable to add license key", err)
		tx.Rollback()
		return err
	}

	err = upsertOauthToken(draftUser, tx)
	if err != nil {
		config.Log("unable to add auth token", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		config.Log("unable to commit selected transaction", err)
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
		config.Log("unable to add license key to signed in user", err)
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
		config.Log("unable to add token", err)
		tx.Rollback()
		return err
	}

	// apply the masked token
	draftUser.PreBuiltToken = maskedID
	return nil
}

// RegisterUser ...
//
// Performs saveUser operation and sends email service to the user to verify registration.
// Creates a random 6 digit token and adds it to the database for the selected user
// if the user is not verified. The random token is sent to the user specific email address.
// Verified users are users whos six digit cryptographic key has matched during the time of registration.
func RegisterUser(userName string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {

	resp, err := db.SaveUser(userName, draftUser)
	if err != nil {
		config.Log("unable to save user", err)
		return nil, err
	}

	PerformEmailNotificationService(userName, draftUser.Email, draftUser.ID.String())

	return resp, nil
}

// PerformEmailNotificationService ...
//
// Updates user fields in db with new token and sends email notification for email verification
// to client using Send Grid api. This function is also re-used when users attempt to re-verify the
// token if Send Grid fails to send the api. The userID is passed in the subject in the token to verify
// if the user is correct.
//
// Error handling is ignored since email notification service failures are ignored and we still want the user
// to login and perform regular operations even without verification of email.
func PerformEmailNotificationService(username string, emailAddress string, userID string) {

	formattedTime, err := strconv.ParseInt(config.DefaultTokenValidityTime, 10, 64)
	if err != nil {
		config.Log("unable to parse provided time", err)
		return
	}

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
		secretToken = ""
	}

	draftCredentials := types.Credentials{
		Claims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: formattedTime,
		},
	}
	credentials, err := stormRider.CreateJWT(&draftCredentials, secretToken)
	if err != nil {
		config.Log("unable to create email token for verification services", err)
		return
	}

	isEmailServiceEnabled := os.Getenv("_SENDGRID_EMAIL_SERVICE")
	if isEmailServiceEnabled != "true" {
		config.Log("email service feature flags are disabled. Email Service is inoperative.", nil)
		return
	}

	sendGridEmailUser := os.Getenv("SEND_GRID_USER")
	if len(sendGridEmailUser) <= 0 {
		config.Log("email service username is not configured. Unable to send email.", nil)
		return
	}

	sendGridUserEmailAddress := os.Getenv("SEND_GRID_USER_EMAIL_ADDRESS")
	if len(sendGridUserEmailAddress) <= 0 {
		config.Log("email service username is not configured. Unable to send email.", nil)
		return
	}

	from := mail.NewEmail(sendGridEmailUser, sendGridUserEmailAddress)
	to := mail.NewEmail(username, emailAddress)

	WebApplicationEndpoint := os.Getenv("REACT_APP_LOCALHOST_URL")
	if len(WebApplicationEndpoint) <= 0 {
		config.Log("unable to determine the web application endpoint", err)
		return
	}

	verificationLink := fmt.Sprintf("%s%sverify?token=%s", WebApplicationEndpoint, WEB_APPLICATION_URI, credentials.Cookie)

	plainText := fmt.Sprintf("Please use the following verification token: %s", credentials.Cookie)
	htmlContent := fmt.Sprintf(`
		<p>Click on the following link to verify your email address:</p>
		<a href="%s">%s</a>
	`, verificationLink, verificationLink)

	message := mail.NewSingleEmail(from, EMAIL_SUBJECT_LINE, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SEND_GRID_API_KEY"))

	_, err = client.Send(message)
	if err != nil {
		config.Log("unable to send email verification", err)
		return
	}

	config.Log("Email notification sent to %s on %+v", nil, emailAddress, time.Now())
}

// ValidateCredentials ...
//
// Method is used to verify if the incoming api calls have a valid jwt token.
// If the validity of the token is crossed, or if the token itself is invalid the error is propogated up the method chain.
func ValidateCredentials(user string, ID string) error {
	db, err := db.SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to setup transaction for selected db", err)
		return err
	}

	var tokenFromDb string
	var expirationTime time.Time
	err = tx.QueryRow(`SELECT token, expiration_time FROM auth.oauth WHERE id=$1 LIMIT 1;`, ID).Scan(&tokenFromDb, &expirationTime)
	if err != nil {
		config.Log("unable to retrive validated token", err)
		tx.Rollback()
		return err
	}

	err = utils.ValidateJwtToken(tokenFromDb)
	if err != nil {
		config.Log("unable to validate jwt token", err)
		tx.Rollback()
		return err
	}

	// Check if the token is within the last 30 seconds of its expiry time
	// token is about to expire. if the user is continuing with activity, create new token
	formattedTimeToLive := time.Until(expirationTime)
	if formattedTimeToLive <= 30*time.Second && formattedTimeToLive > 0 {

		formattedTime, err := strconv.ParseInt(config.DefaultTokenValidityTime, 10, 64)
		if err != nil {
			config.Log("unable to parse provided time", err)
			return err
		}

		secretToken := os.Getenv("TOKEN_SECRET_KEY")
		if len(secretToken) <= 0 {
			config.Log("unable to retrieve secret token key. defaulting to default values", nil)
			secretToken = ""
		}

		draftCredentials := &types.Credentials{
			Claims: jwt.StandardClaims{
				ExpiresAt: formattedTime,
			},
		}

		updatedToken, err := stormRider.RefreshToken(draftCredentials, secretToken)
		if err != nil {
			config.Log("unable to refresh token", err)
			tx.Rollback()
			return err
		}

		parsedUserID, err := uuid.Parse(ID)
		if err != nil {
			config.Log("unable to determine user id", err)
			return err
		}

		tokenValidityMinutes, err := strconv.Atoi(config.DefaultTokenValidityTime)
		if err != nil {
			config.Log("Invalid token validity time: %v", err)
			return err
		}

		draftUser := model.UserCredentials{
			ID:             parsedUserID,
			ExpirationTime: time.Now().Add((time.Duration(tokenValidityMinutes) * time.Minute)),
			PreBuiltToken:  updatedToken,
		}
		err = upsertOauthToken(&draftUser, tx)
		if err != nil {
			config.Log("unable to revalidate the user", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		config.Log("unable to commit transaction", err)
		tx.Rollback()
		return err
	}

	return nil
}

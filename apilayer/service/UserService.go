package service

import (
	"fmt"
	"log"
	"os"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/ciri-stormrider/types"
	"github.com/mohit2530/communityCare/db"
	"github.com/mohit2530/communityCare/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	EMAIL_SUBJECT_LINE  = "Verify your email address for FleetWise Application"
	WEB_APPLICATION_URI = "/api/v1/verify/"
)

// RegisterUser ...
//
// Performs saveUser operation and sends email service to the user to verify registration.
// Creates a random 6 digit token and adds it to the database for the selected user
// if the user is not verified. The random token is sent to the user specific email address.
// Verified users are users whos six digit cryptographic key has matched during the time of registration.
func RegisterUser(userName string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {

	resp, err := db.SaveUser(userName, draftUser)
	if err != nil {
		log.Printf("unable to save user. error: %+v", err)
		return nil, err
	}

	PerformEmailNotificationService(userName, draftUser.Email)
	return resp, nil
}

// PerformEmailNotificationService ...
//
// Updates user fields in db with new token and sends email notification for email verification
// to client using Send Grid api. This function is also re-used when users attempt to re-verify the
// token if Send Grid fails to send the api.
//
// Error handling is ignored since email notification service failures are ignored and we still want the user
// to login and perform regular operations even without verification of email.
func PerformEmailNotificationService(username string, emailAddress string) {

	credentials, err := stormRider.CreateJWT(&types.Credentials{}, "15", "")
	if err != nil {
		log.Printf("unable to create email token for verification services. error: %+v", err)
		return
	}

	isEmailServiceEnabled := os.Getenv("_SENDGRID_EMAIL_SERVICE")
	if isEmailServiceEnabled != "true" {
		log.Printf("email service feature flags are disabled. Email Service is inoperative.")
		return
	}

	sendGridEmailUser := os.Getenv("SEND_GRID_USER")
	if len(sendGridEmailUser) <= 0 {
		log.Printf("email service username is not configured. Unable to send email.")
		return
	}

	sendGridUserEmailAddress := os.Getenv("SEND_GRID_USER_EMAIL_ADDRESS")
	if len(sendGridUserEmailAddress) <= 0 {
		log.Printf("email service username is not configured. Unable to send email.")
		return
	}

	from := mail.NewEmail(sendGridEmailUser, sendGridUserEmailAddress)
	to := mail.NewEmail(username, emailAddress)

	WebApplicationEndpoint := os.Getenv("REACT_APP_LOCALHOST_URL")
	if len(WebApplicationEndpoint) <= 0 {
		log.Printf("unable to determine the web application endpoint. error: %+v", err)
		return
	}

	verificationLink := fmt.Sprintf("%s%sverify?token=%s", WebApplicationEndpoint, WEB_APPLICATION_URI, credentials.Cookie)

	plainText := fmt.Sprintf("Please use the following verification token: %s", credentials.Cookie)
	htmlContent := fmt.Sprintf(`
		<p>Click on the following link to verify your email address:</p>
		<a href="%s">%s</a>
	`, verificationLink, verificationLink)

	message := mail.NewSingleEmail(from, EMAIL_SUBJECT_LINE, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	_, err = client.Send(message)
	if err != nil {
		log.Printf("unable to send email verification. error: %+v", err)
		return
	}
}

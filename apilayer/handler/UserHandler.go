package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mohit2530/communityCare/db"
	"github.com/mohit2530/communityCare/model"
	"github.com/mohit2530/communityCare/service"
)

// Signup ...
// swagger:route POST /api/v1/signup Authentication signup
//
// # Sign up users into the database system.
//
// Parameters:
//   - +name: email
//     in: query
//     description: The email address of the current user
//     type: string
//     required: true
//   - +name: password
//     in: query
//     description: The password of the current user
//     type: string
//     required: true
//   - +name: birthday
//     in: query
//     description: The birthdate of the current user. Must be 13 years of age.
//     type: string
//     required: true
//   - +name: role
//     in: query
//     description: The user role for the application.
//     type: string
//     default: false
//     required: false
//
// Responses:
// 200: UserCredentials
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Signup(rw http.ResponseWriter, r *http.Request) {

	draftUser := &model.UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if len(draftUser.Email) <= 0 || len(draftUser.EncryptedPassword) <= 0 {
		log.Printf("unable to decode user")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	if len(draftUser.Username) <= 3 {
		log.Printf("user name is required and must be at least 4 characters in length")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	if len(draftUser.Role) <= 0 {
		draftUser.Role = "USER"
	}

	t, err := time.Parse("2006-01-02", draftUser.Birthday)
	if err != nil {
		log.Printf("Error parsing birthdate. Error: %+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	// Check if the user is at least 13 years old
	age := time.Now().Year() - +t.Year()
	if age <= 13 {
		log.Println("unable to sign up user. verification failed.")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	backendClientUsr := os.Getenv("CLIENT_USER")
	if len(backendClientUsr) == 0 {
		log.Printf("unable to retrieve user from env.")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
	}

	resp, err := service.RegisterUser(backendClientUsr, draftUser)
	if err != nil {
		log.Printf("Unable to create new user. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// Signin ...
// swagger:route POST /api/v1/signin Authentication signin
//
// # Sign in users into the database system.
//
// Parameters:
//   - +name: email
//     in: query
//     description: The email address of the current user
//     type: string
//     required: true
//   - +name: password
//     in: query
//     description: The password of the current user
//     type: string
//     required: true
//
// Responses:
// 200: UserCredentials
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Signin(rw http.ResponseWriter, r *http.Request) {

	draftUser := &model.UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		log.Printf("Unable to decode request parameters. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if len(draftUser.Email) <= 0 || len(draftUser.EncryptedPassword) <= 0 {
		log.Printf("Unable to decode user. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	draftUser.UserAgent = r.UserAgent()
	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		log.Printf("unable to retrieve user from env. Unable to sign in.")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}
	resp, err := service.FetchUser(user, draftUser)
	if err != nil {
		log.Printf("Unable to sign user in. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   draftUser.PreBuiltToken,
		Expires: draftUser.ExpirationTime,
	})

	rw.Header().Add("Role2", draftUser.Role)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp.ID)
}

// IsValidUserEmail ...
// swagger:route POST /api/v1/isValidEmail Authentication IsValidUserEmail
//
// # Returns true or false is the user is already in the system
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func IsValidUserEmail(rw http.ResponseWriter, r *http.Request) {

	draftUserEmail := &model.UserEmail{}
	err := json.NewDecoder(r.Body).Decode(draftUserEmail)
	r.Body.Close()
	if err != nil {
		log.Printf("unable to validate user email address. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		log.Printf("unable to retrieve user from env. Unable to sign in.")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}

	resp, err := db.IsValidUserEmail(user, draftUserEmail.EmailAddress)
	if err != nil {
		log.Printf("unable to verify user email address. error: %+v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// VerifyEmailAddress ...
// swagger:route GET /api/v1/verify Authentication VerifyEmailAddress
//
// # Used to verify if the user correctly verified the selected email address. If the token
// is valid, then the user was successfully verified
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func VerifyEmailAddress(rw http.ResponseWriter, r *http.Request) {

}

// ResetEmailToken ...
// swagger:route POST /api/v1/reset Authentication ResetEmailToken
//
// # Resets the token in the database and allows users to resend email in case the token
// is incorrect or failed to reach the user.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func ResetEmailToken(rw http.ResponseWriter, r *http.Request) {

	draftUserEmail := &model.UserEmail{}
	err := json.NewDecoder(r.Body).Decode(draftUserEmail)
	r.Body.Close()
	if err != nil {
		log.Printf("unable to validate user email address. error: +%v", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		log.Printf("unable to retrieve user from env. Unable to sign in.")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}

	service.PerformEmailNotificationService(user, draftUserEmail.EmailAddress)

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("200 OK")
}

// Logout ...
// swagger:route POST /api/v1/logout Authentication logout
//
// # Logs users out of the database system.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Logout(rw http.ResponseWriter, r *http.Request) {

	// immediately clear the token cookie
	http.SetCookie(rw, &http.Cookie{
		Name:     "token",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(nil)
}

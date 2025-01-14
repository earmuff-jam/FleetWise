package model

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// UserCredentials ...
// swagger:model UserCredentials
//
// UserCredentials object. Used for authentication purposes only.
type UserCredentials struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Email             string    `json:"email,omitempty"`
	Username          string    `json:"username,omitempty"`
	Birthday          string    `json:"birthday,omitempty"`
	Role              string    `json:"role"`
	UserAgent         string    `json:"user_agent,omitempty"`
	EncryptedPassword string    `json:"password,omitempty"`
	PreBuiltToken     string    `json:"pre_token,omitempty"`
	LicenceKey        string    `json:"licence_key,omitempty"`
	ExpirationTime    time.Time `json:"expiration_time,omitempty"`
	jwt.StandardClaims
}

// User ...
// swagger:model User
//
// User object. Used for authentication purposes only.
type User struct {
	InstanceID               uuid.UUID
	ID                       uuid.UUID       `json:"id,omitempty"`
	Aud                      string          `json:"aud,omitempty"`
	Role                     string          `json:"role,omitempty"`
	Email                    string          `json:"email,omitempty"`
	FirstName                string          `json:"first_name,omitempty"`
	MiddleName               string          `json:"middle_name,omitempty"`
	LastName                 string          `json:"last_name,omitempty"`
	EncryptedPassword        string          `json:"encrypted_password,omitempty"`
	EmailConfirmedAt         time.Time       `json:"email_confirmed_at,omitempty"`
	InvitedAt                time.Time       `json:"invited_at,omitempty"`
	ConfirmationToken        string          `json:"confirmation_token,omitempty"`
	ConfirmationSentAt       time.Time       `json:"confirmation_sent_at,omitempty"`
	RecoveryToken            string          `json:"recovery_token,omitempty"`
	RecoverySentAt           time.Time       `json:"recovery_sent_at,omitempty"`
	EmailChangeTokenNew      string          `json:"email_change_token_new,omitempty"`
	EmailChange              string          `json:"email_change,omitempty"`
	EmailChangeSentAt        time.Time       `json:"email_change_sent_at,omitempty"`
	LastSignInAt             time.Time       `json:"last_sign_in_at,omitempty"`
	RawAppMetaData           json.RawMessage `json:"raw_app_meta_data,omitempty"`
	RawUserMetaData          json.RawMessage `json:"raw_user_meta_data,omitempty"`
	CreatedAt                time.Time       `json:"created_at,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at,omitempty"`
	Phone                    string          `json:"phone,omitempty"`
	PhoneConfirmedAt         time.Time       `json:"phone_confirmed_at,omitempty"`
	PhoneChange              string          `json:"phone_change,omitempty"`
	PhoneChangeToken         string          `json:"phone_change_token,omitempty"`
	PhoneChangeSentAt        time.Time       `json:"phone_change_sent_at,omitempty"`
	ConfirmedAt              time.Time       `json:"confirmed_at,omitempty"`
	EmailChangeTokenCurrent  string          `json:"email_change_token_current,omitempty"`
	EmailChangeConfirmStatus int             `json:"email_change_confirm_status,omitempty"`
	BannedUntil              time.Time       `json:"banned_until,omitempty"`
	ReAuthenticationToken    string          `json:"re_authentication_token,omitempty"`
	ReAuthenticationSentAt   time.Time       `json:"re_authentication_sent_at,omitempty"`
	IsSSOUser                bool            `json:"is_sso_user,omitempty"`
	DeletedAt                time.Time       `json:"deleted_at,omitempty"`
}

// UserEmail ...
// swagger:model UserEmail
//
// UserEmail object. Used to validate if the email already exists in the db
type UserEmail struct {
	EmailAddress string `json:"email"`
}

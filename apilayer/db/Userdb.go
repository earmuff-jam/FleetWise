package db

import (
	"time"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SaveUser ...
//
// Function is used to persist the user into the database and validate the password fields.
// Errors are propogated up the system if found.
func SaveUser(user string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup database", err)
		return nil, err
	}
	defer db.Close()

	// generate the hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(draftUser.EncryptedPassword), 8)
	if err != nil {
		config.Log("unable to decode password", err)
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to setup transaction", err)
		return nil, err
	}

	sqlStr := `INSERT INTO auth.users(email, username, birthdate, encrypted_password, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	var draftUserID string

	config.Log("SqlStr: %s", nil, sqlStr)
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
		config.Log("unable to query selected row", err)
		return nil, err
	}

	draftUser.ID, err = uuid.Parse(draftUserID)
	if err != nil {
		tx.Rollback()
		config.Log("invalid id for user detected", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		config.Log("unable to commit to database", err)
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
		config.Log("unable to setup database", err)
		return nil, err
	}
	defer db.Close()

	// retrive the encrypted pwd. EMAIL must be UNIQUE field.
	sqlStr := `SELECT id, role, encrypted_password, is_verified FROM auth.users WHERE email=$1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	result := db.QueryRow(sqlStr, draftUser.Email)

	storedCredentials := &model.UserCredentials{}
	err = result.Scan(&storedCredentials.ID, &storedCredentials.Role, &storedCredentials.EncryptedPassword, &storedCredentials.IsVerified)
	if err != nil {
		config.Log("unable to retrieve user details", err)
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.EncryptedPassword), []byte(draftUser.EncryptedPassword)); err != nil {
		config.Log("unable to match password", err)
		return nil, err
	}

	draftUser.ID = storedCredentials.ID
	draftUser.Role = storedCredentials.Role
	draftUser.IsVerified = storedCredentials.IsVerified

	return draftUser, nil
}

// IsValidUserEmail ...
//
// Function is used to validate any email address if they are of the correct form
func IsValidUserEmail(user string, draftUserEmail string) (bool, error) {
	db, err := SetupDB(user)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// retrive the encrypted pwd. EMAIL must be UNIQUE field.
	sqlStr := `SELECT EXISTS(SELECT 1 FROM auth.users u WHERE u.email=$1);`

	config.Log("SqlStr: %s", nil, sqlStr)
	result := db.QueryRow(sqlStr, draftUserEmail)
	exists := false
	err = result.Scan(&exists)
	if err != nil {
		config.Log("unable to validate user email address", err)
		return false, err
	}
	return !exists, nil // return false if found
}

// VerifyUser ...
//
// Function is used to verify the user from the email login
func VerifyUser(user string, draftUserID string) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `UPDATE auth.users
	SET is_verified = $1, email_confirmed_at = $2
	WHERE id = $3;`

	config.Log("SqlStr: %s", nil, sqlStr)
	_, err = db.Exec(sqlStr, true, time.Now(), draftUserID)
	if err != nil {
		config.Log("failed to update user verification", err)
		return err
	}

	config.Log("user %s successfully verified", nil, draftUserID)
	return nil
}

// RemoveUser ...
//
// Removes the user from the database
func RemoveUser(user string, id uuid.UUID) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM auth.users WHERE id = $1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		config.Log("Error deleting user with ID %s", err, id.String())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		config.Log("unable to retrieve selected rows", err)
		return err
	}

	if rowsAffected == 0 {
		config.Log("unable to find the selected user details", nil)
		return nil
	}

	return nil
}

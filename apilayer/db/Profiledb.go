package db

import (
	"database/sql"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/mohit2530/communityCare/model"
)

// FetchAllUserProfiles ...
func FetchAllUserProfiles(user string) ([]model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `
        SELECT 
            id,
            username,
            full_name,
            COALESCE(ENCODE(avatar_url::bytea, 'base64'), '') AS base64,
            email_address,
            phone_number,
            goal,
            about_me,
            onlinestatus,
            role
        FROM community.profiles;
    `

	var profiles []model.Profile

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profile model.Profile

		var username, fullName, emailAddress, phoneNumber, avatarUrl, goal, aboutMe, role sql.NullString

		if err := rows.Scan(&profile.ID, &username, &fullName, &avatarUrl, &emailAddress, &phoneNumber, &goal, &aboutMe, &profile.OnlineStatus, &role); err != nil {
			return nil, err
		}

		if username.Valid {
			profile.Username = username.String
		}
		if fullName.Valid {
			profile.FullName = fullName.String
		}
		if avatarUrl.Valid {
			profile.AvatarUrl = avatarUrl.String
		}
		if emailAddress.Valid {
			profile.EmailAddress = emailAddress.String
		}
		if phoneNumber.Valid {
			profile.PhoneNumber = phoneNumber.String
		}
		if aboutMe.Valid {
			profile.AboutMe = aboutMe.String
		}

		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

// FetchUserProfile ...
func FetchUserProfile(user string, userID string) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `
		SELECT 
			id,
			username,
			full_name,
			CASE 
				WHEN avatar_url IS NOT NULL THEN ENCODE(avatar_url::bytea,'base64')
				ELSE ''
			END AS base64,
			email_address,
			phone_number,
			goal,
			about_me,
			onlinestatus,
			role,
			updated_at
		FROM community.profiles
		WHERE id=$1
	`

	var draftProfile model.Profile

	var profileID sql.NullString
	var userName sql.NullString
	var fullName sql.NullString
	var avatarUrl sql.NullString
	var emailAddress sql.NullString
	var phoneNumber sql.NullString
	var goal sql.NullString
	var aboutMe sql.NullString
	var onlineStatus sql.NullBool
	var role sql.NullString
	var updated_at sql.NullTime

	rows, err := db.Query(sqlStr, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&profileID, &userName, &fullName, &avatarUrl, &emailAddress, &phoneNumber, &goal, &aboutMe, &onlineStatus, &role, &updated_at); err != nil {
			return nil, err
		}
		draftProfile.ID, _ = uuid.Parse(profileID.String)
		draftProfile.Username = userName.String
		draftProfile.FullName = fullName.String
		draftProfile.AvatarUrl = avatarUrl.String
		draftProfile.EmailAddress = emailAddress.String
		draftProfile.PhoneNumber = phoneNumber.String
		draftProfile.AboutMe = aboutMe.String
		draftProfile.OnlineStatus = onlineStatus.Bool
		draftProfile.UpdatedAt = updated_at.Time
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	draftProfile.Validate()
	return &draftProfile, nil
}

// UpdateUserProfile ...
func UpdateUserProfile(user string, userID string, draftProfile model.Profile) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	sqlStr := `
		UPDATE community.profiles 
		SET username=$2, full_name=$3, email_address=$4, phone_number=$5, about_me=$6, onlinestatus=$7
		WHERE id=$1
		RETURNING id, username, full_name, avatar_url, email_address, phone_number, about_me, onlinestatus
	`

	var updatedProfile model.Profile

	var avatarUrl sql.NullString // Assuming avatar_url is a string column, not bytea

	row := tx.QueryRow(sqlStr,
		userID,
		draftProfile.Username,
		draftProfile.FullName,
		draftProfile.EmailAddress,
		draftProfile.PhoneNumber,
		draftProfile.AboutMe,
		draftProfile.OnlineStatus,
	)

	err = row.Scan(
		&updatedProfile.ID,
		&updatedProfile.Username,
		&updatedProfile.FullName,
		&avatarUrl,
		&updatedProfile.EmailAddress,
		&updatedProfile.PhoneNumber,
		&updatedProfile.AboutMe,
		&updatedProfile.OnlineStatus,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if avatarUrl.Valid {
		updatedProfile.AvatarUrl = avatarUrl.String
	} else {
		updatedProfile.AvatarUrl = ""
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updatedProfile.Validate()
	return &updatedProfile, nil
}

// UpdateProfileAvatar ...
func UpdateProfileAvatar(user string, userID string, header *multipart.FileHeader, fileBytes []byte) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	var updatedProfile model.Profile

	sqlStr := `
		UPDATE community.profiles 
		SET avatar_url = $2
		WHERE id = $1
		RETURNING id, username, full_name, avatar_url, phone_number, goal, about_me, onlinestatus, role
	`

	var avatarUrl sql.NullString // Assuming avatar_url is a string column, not bytea
	var goal sql.NullString

	row := tx.QueryRow(sqlStr, userID, fileBytes)

	err = row.Scan(
		&updatedProfile.ID,
		&updatedProfile.Username,
		&updatedProfile.FullName,
		&avatarUrl,
		&updatedProfile.PhoneNumber,
		&goal,
		&updatedProfile.AboutMe,
		&updatedProfile.OnlineStatus,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if avatarUrl.Valid {
		updatedProfile.AvatarUrl = avatarUrl.String
	} else {
		updatedProfile.AvatarUrl = ""
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updatedProfile.Validate()
	return &updatedProfile, nil
}

package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RetrieveNotes ...
func RetrieveNotes(user string, userID uuid.UUID) ([]model.Note, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
	n.id,
	n.title, 
	n.description,
	s.id,
	s.name AS status_name,
	s.description AS status_description,
	n.color,
	n.completionDate,
	n.location[0] AS lon, -- Extract longitude from POINT
	n.location[1] AS lat, -- Extract latitude from POINT
	n.created_at,
	n.created_by,
	COALESCE(cp.full_name, cp.username, cp.email_address) AS creator_name,
	n.updated_at,
	n.updated_by,
	COALESCE(up.full_name, up.username, up.email_address)  AS updater_name
	FROM community.notes n
	LEFT JOIN community.statuses s on s.id = n.status
	LEFT JOIN community.profiles cp on cp.id = n.created_by
	LEFT JOIN community.profiles up on up.id = n.updated_by
	WHERE $1::UUID = ANY(n.sharable_groups)
	ORDER BY n.updated_at DESC;`

	config.Log("SqlStr: %s", nil, sqlStr)
	rows, err := db.Query(sqlStr, userID)
	if err != nil {
		config.Log("unable to retrieve selected details", err)
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note

	for rows.Next() {
		var note model.Note
		var lon, lat sql.NullFloat64
		var statusID sql.NullString
		var statusName sql.NullString
		var statusDescription sql.NullString
		var completionDate sql.NullTime

		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &statusID, &statusName, &statusDescription, &note.Color, &completionDate, &lon, &lat, &note.CreatedAt, &note.CreatedBy, &note.Creator, &note.UpdatedAt, &note.UpdatedBy, &note.Updator); err != nil {
			config.Log("unable to scan selected values", err)
			return nil, err
		}

		if completionDate.Valid {
			note.CompletionDate = &completionDate.Time
		} else {
			note.CompletionDate = nil
		}

		if statusID.Valid {
			note.Status = statusID.String
		}

		if statusName.Valid {
			note.StatusName = statusName.String
		}

		if statusDescription.Valid {
			note.StatusDescription = statusDescription.String
		}

		if lon.Valid && lat.Valid {
			note.Location = model.Location{
				Lon: lon.Float64,
				Lat: lat.Float64,
			}
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		config.Log("unable to scan selected data", err)
		return nil, err
	}

	return notes, nil
}

// AddNewNote ...
func AddNewNote(user string, userID string, draftNote model.Note) (*model.Note, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	selectedStatusDetails, err := RetrieveStatusDetails(user, draftNote.Status)
	if err != nil {
		config.Log("unable to retrieve selected status details", err)
		return nil, err
	}
	if selectedStatusDetails == nil {
		config.Log("selected status details empty", errors.New("unable to find selected status"))
		return nil, errors.New("unable to find selected status")
	}

	sqlStr := `
		INSERT INTO community.notes (title, description, status, color, completionDate, location, created_by, updated_by, sharable_groups) 
		VALUES ($1, $2, $3, $4, $5, POINT($6, $7), $8, $9, $10) 
		RETURNING id;`

	parsedCreatorID, err := uuid.Parse(draftNote.UpdatedBy)
	if err != nil {
		config.Log("unable to parse creator id", err)
		return nil, err
	}

	var draftNoteID string
	draftNote.CreatedAt = time.Now()
	draftNote.UpdatedAt = time.Now()

	var sharableGroups = make([]uuid.UUID, 0)
	sharableGroups = append(sharableGroups, parsedCreatorID)

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	config.Log("SqlStr: %s", nil, sqlStr)
	err = tx.QueryRow(
		sqlStr,
		draftNote.Title,
		draftNote.Description,
		selectedStatusDetails.ID,
		draftNote.Color,
		draftNote.CompletionDate,
		draftNote.Location.Lon,
		draftNote.Location.Lat,
		parsedCreatorID,
		parsedCreatorID,
		pq.Array(sharableGroups),
	).Scan(&draftNoteID)

	if err != nil {
		tx.Rollback()
		config.Log("unable to scan selected values", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		config.Log("unable to commit selected transaction", err)
		return nil, err
	}

	sqlStr = `SELECT id, username, full_name from community.profiles p where p.id = $1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	row := db.QueryRow(sqlStr, userID)

	var creatorID string
	var creatorUsername sql.NullString
	var creatorFullName sql.NullString
	err = row.Scan(&creatorID, &creatorUsername, &creatorFullName)
	if err != nil {
		config.Log("creator not found", err)
		return nil, err
	}

	draftNote.ID = draftNoteID
	draftNote.StatusName = selectedStatusDetails.Name
	draftNote.StatusDescription = selectedStatusDetails.Description

	if creatorUsername.Valid {
		// creator === updator for the 1st time
		draftNote.Creator = creatorUsername.String
		draftNote.Updator = creatorUsername.String
	} else if creatorFullName.Valid {
		draftNote.Creator = creatorFullName.String
		draftNote.Updator = creatorFullName.String
	}
	return &draftNote, nil
}

// UpdateNote ...
func UpdateNote(user string, userID string, draftNote model.Note) (*model.Note, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	// retrieve selected status
	selectedStatusDetails, err := RetrieveStatusDetails(user, draftNote.Status)
	if err != nil {
		config.Log("unable to retrieve selected status details", err)
		return nil, err
	}
	if selectedStatusDetails == nil {
		config.Log("selected status details is empty", errors.New("unable to find selected status"))
		return nil, errors.New("unable to find selected status")
	}

	sqlStr := `UPDATE community.notes 
	SET
	title = $2,
	description = $3,
	status = $4,
	color = $5,
	location = POINT($6, $7),
	updated_by = $8,
	updated_at = $9
	WHERE id = $1
	RETURNING id, title, description, color, created_at, created_by, updated_at, updated_by;`

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to begin transaction", err)
		tx.Rollback()
		return nil, err
	}

	parsedCreatorID, err := uuid.Parse(draftNote.UpdatedBy)
	if err != nil {
		config.Log("unable to parse user id", err)
		tx.Rollback()
		return nil, err
	}

	var updatedNote model.Note
	config.Log("SqlStr: %s", nil, sqlStr)

	row := tx.QueryRow(sqlStr,
		draftNote.ID,
		draftNote.Title,
		draftNote.Description,
		selectedStatusDetails.ID,
		draftNote.Color,
		draftNote.Location.Lon,
		draftNote.Location.Lat,
		parsedCreatorID,
		time.Now(),
	)

	err = row.Scan(
		&updatedNote.ID,
		&updatedNote.Title,
		&updatedNote.Description,
		&updatedNote.Color,
		&updatedNote.CreatedAt,
		&updatedNote.CreatedBy,
		&updatedNote.UpdatedAt,
		&updatedNote.UpdatedBy,
	)

	if err != nil {
		config.Log("unable to scan selected details", err)
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		config.Log("unable to commit selected transaction", err)
		return nil, err
	}

	updatedNote.Status = selectedStatusDetails.ID.String()
	updatedNote.StatusName = selectedStatusDetails.Name
	updatedNote.StatusDescription = selectedStatusDetails.Description
	updatedNote.Location.Lat = draftNote.Location.Lat
	updatedNote.Location.Lon = draftNote.Location.Lon

	return &updatedNote, nil
}

// RemoveNote ...
func RemoveNote(user string, draftNoteID string) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.notes WHERE id=$1;`
	config.Log("SqlStr: %s", nil, sqlStr)

	_, err = db.Exec(sqlStr, draftNoteID)
	if err != nil {
		config.Log("unable to delete selected note", err)
		return err
	}
	return nil
}

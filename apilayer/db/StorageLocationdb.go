package db

import (
	"database/sql"
	"time"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RetrieveAllStorageLocation ...
func RetrieveAllStorageLocation(user string) ([]model.StorageLocation, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT sl.id, sl.location, sl.created_at, sl.created_by, sl.updated_at, sl.updated_by, sl.sharable_groups 
		FROM community.storage_locations sl 
		ORDER BY sl.updated_at;`

	config.Log("SqlStr: %s", nil, sqlStr)
	rows, err := db.Query(sqlStr)

	if err != nil {
		config.Log("unable to retrieve query details", err)
		return nil, err
	}
	defer rows.Close()

	var data []model.StorageLocation

	var storageLocationID sql.NullString
	var storageLocation sql.NullString
	var createdBy sql.NullString
	var createdAt sql.NullTime
	var updatedBy sql.NullString
	var updatedAt sql.NullTime
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.StorageLocation
		if err := rows.Scan(&storageLocationID, &storageLocation, &createdAt, &createdBy, &updatedAt, &updatedBy, &sharableGroups); err != nil {
			config.Log("unable to scan selected details", err)
			return nil, err
		}

		ec.ID = storageLocationID.String
		ec.Location = storageLocation.String
		ec.CreatedAt = createdAt.Time
		ec.CreatedBy = createdBy.String
		ec.UpdatedAt = updatedAt.Time
		ec.UpdatedBy = updatedBy.String
		ec.SharableGroups = sharableGroups

		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		config.Log("unable to view selected rows", err)
		return nil, err
	}

	return data, nil
}

// DeleteStorageLocation
func DeleteStorageLocation(user string, locationID string) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.storage_locations WHERE id=$1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	_, err = db.Exec(sqlStr, locationID)
	if err != nil {
		config.Log("unable to delete selected storage locations", err)
		return err
	}
	return nil
}

// addNewStorageLocation ...
//
// adds new storage location if not existing but if there was an existing storage location, we just return that ID
func addNewStorageLocation(user string, draftLocation string, created_by string, emptyLocationID *string) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	var count int
	fetchSqlStr := `SELECT count(sl.id), sl.id FROM community.storage_locations sl WHERE sl.location = $1 GROUP BY sl.id;`

	config.Log("SqlStr: %s", nil, fetchSqlStr)
	err = db.QueryRow(fetchSqlStr, draftLocation).Scan(&count, emptyLocationID)
	if err != nil {
		config.Log("found existing location for selected item. using existing location for %+v", nil, draftLocation)
	}

	// save new storage location if it does not already exists
	if count == 0 {
		config.Log("new storage location is found. adding it into the system", nil)
		sqlStr := `INSERT INTO community.storage_locations(location, created_by, updated_by, created_at, updated_at, sharable_groups) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

		var locationID string
		var sharableGroups = make([]uuid.UUID, 0)

		userID, err := uuid.Parse(created_by)
		if err != nil {
			config.Log("unable to parse selected user id", err)
			return err
		}
		sharableGroups = append(sharableGroups, userID)

		err = db.QueryRow(
			sqlStr,
			draftLocation,
			created_by,
			created_by,
			time.Now(),
			time.Now(),
			pq.Array(sharableGroups),
		).Scan(&locationID)

		if err != nil {
			config.Log("unable to query selected details", err)
			return err
		}
		*emptyLocationID = locationID
	}

	return nil
}

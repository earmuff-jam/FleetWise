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

// RetrieveAllMaintenancePlans ...
func RetrieveAllMaintenancePlans(user string, userID string, limit int) (*[]model.MaintenancePlan, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
	mp.id,
	mp.name,
	mp.description,
	ms.id,
	ms.name AS status_name,
	ms.description AS status_description,
	mp.color, 
	mp.plan_type,
	mp.location[0] AS lon,
	mp.location[1] AS lat,
	mp.created_at,
	mp.created_by,
	COALESCE(cp.full_name, cp.username, cp.email_address) AS creator_name, 
	mp.updated_at,
	mp.updated_by,
	COALESCE(up.full_name, up.username, up.email_address) AS updater_name,
	mp.sharable_groups
	FROM community.maintenance_plan mp
	LEFT JOIN community.statuses ms on ms.id = mp.status
	LEFT JOIN community.profiles cp on cp.id = mp.created_by
	LEFT JOIN community.profiles up on up.id = mp.updated_by
	WHERE $1::UUID = ANY(mp.sharable_groups)
	ORDER BY mp.updated_at DESC
	LIMIT $2;`

	config.Log("SqlStr: %s", nil, sqlStr)
	rows, err := db.Query(sqlStr, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.MaintenancePlan

	var lon, lat sql.NullFloat64
	var maintenancePlanID sql.NullString
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.MaintenancePlan
		if err := rows.Scan(&maintenancePlanID, &ec.Name, &ec.Description, &ec.Status, &ec.StatusName, &ec.StatusDescription,
			&ec.Color, &ec.PlanType, &lon, &lat, &ec.CreatedAt, &ec.CreatedBy, &ec.Creator, &ec.UpdatedAt, &ec.UpdatedBy, &ec.Updator, &sharableGroups); err != nil {
			config.Log("unable to scan selected maintenance plan", err)
			return nil, err
		}

		if lon.Valid && lat.Valid {
			ec.Location = model.Location{
				Lon: lon.Float64,
				Lat: lat.Float64,
			}
		}

		content, _, _, err := FetchImage(maintenancePlanID.String)
		if err != nil {
			if err.Error() == "NoSuchKey" {
				config.Log("cannot find the selected document", err)
			}
		}

		ec.Image = content
		ec.ID = maintenancePlanID.String
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		config.Log("unable to validate all rows", err)
		return nil, err
	}

	return &data, nil
}

// RetrieveMaintenancePlan ...
func RetrieveMaintenancePlan(user string, userID string, maintenanceID string) (model.MaintenancePlan, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return model.MaintenancePlan{}, err
	}
	defer db.Close()

	sqlStr := `SELECT 
	mp.id,
	mp.name,
	mp.description,
	ms.name AS status_name,
	ms.description AS status_description,
	mp.color, 
	mp.location[0] AS lon,
	mp.location[1] AS lat,
	mp.plan_type,
	mp.created_at,
	mp.created_by,
	COALESCE(cp.full_name, cp.username, cp.email_address) AS creator_name, 
	mp.updated_at,
	mp.updated_by,
	COALESCE(up.full_name, up.username, up.email_address) AS updater_name,
	mp.sharable_groups
	FROM community.maintenance_plan mp
	LEFT JOIN community.statuses ms on ms.id = mp.status
	LEFT JOIN community.profiles cp on cp.id = mp.created_by
	LEFT JOIN community.profiles up on up.id = mp.updated_by
	WHERE $1::UUID = ANY(mp.sharable_groups) AND mp.id = $2;`

	config.Log("SqlStr: %s", nil, sqlStr)
	row := db.QueryRow(sqlStr, userID, maintenanceID)
	selectedMaintenancePlan := model.MaintenancePlan{}

	var lon, lat sql.NullFloat64

	err = row.Scan(
		&selectedMaintenancePlan.ID,
		&selectedMaintenancePlan.Name,
		&selectedMaintenancePlan.Description,
		&selectedMaintenancePlan.StatusName,
		&selectedMaintenancePlan.StatusDescription,
		&selectedMaintenancePlan.Color,
		&lon,
		&lat,
		&selectedMaintenancePlan.PlanType,
		&selectedMaintenancePlan.CreatedAt,
		&selectedMaintenancePlan.CreatedBy,
		&selectedMaintenancePlan.Creator,
		&selectedMaintenancePlan.UpdatedAt,
		&selectedMaintenancePlan.UpdatedBy,
		&selectedMaintenancePlan.Updator,
		pq.Array(&selectedMaintenancePlan.SharableGroups),
	)

	if lon.Valid && lat.Valid {
		selectedMaintenancePlan.Location = model.Location{
			Lon: lon.Float64,
			Lat: lat.Float64,
		}
	}

	if err != nil {
		config.Log("unable to find the selected maintenance plan", err)
		return model.MaintenancePlan{}, err
	}

	return selectedMaintenancePlan, nil
}

// RetrieveAllMaintenancePlanItems ...
func RetrieveAllMaintenancePlanItems(user string, userID string, maintenancePlanID string, limit int) ([]model.MaintenanceItemResponse, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
		mi.id,
		mi.maintenance_plan_id,
		mi.item_id,
		i.name,
		i.description,
		i.price,
		i.quantity,
		i.location,
		mi.created_by,
		COALESCE(cp.username, cp.full_name, cp.email_address, 'Anonymous') as creator,
		mi.created_at,
		mi.updated_by,
		COALESCE(up.username, up.full_name, cp.email_address, 'Anonymous') as updator,
		mi.updated_at,
		mi.sharable_groups
	FROM community.maintenance_item mi
	LEFT JOIN community.inventory i ON mi.item_id = i.id
	LEFT JOIN community.profiles cp ON mi.created_by = cp.id
	LEFT JOIN community.profiles up ON mi.updated_by = up.id
	WHERE $1::UUID = ANY(mi.sharable_groups) AND mi.maintenance_plan_id = $2
	ORDER BY mi.updated_at DESC FETCH FIRST $3 ROWS ONLY;`

	config.Log("SqlStr: %s", nil, sqlStr)
	rows, err := db.Query(sqlStr, userID, maintenancePlanID, limit)
	if err != nil {
		config.Log("unable to retrieve maintenance items", err)
		return nil, err
	}
	defer rows.Close()

	var data []model.MaintenanceItemResponse
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.MaintenanceItemResponse
		if err := rows.Scan(&ec.ID, &ec.MaintenancePlanID, &ec.ItemID, &ec.Name, &ec.Description, &ec.Price, &ec.Quantity, &ec.Location, &ec.CreatedBy, &ec.Creator, &ec.CreatedAt, &ec.UpdatedBy, &ec.Updator, &ec.UpdatedAt, &sharableGroups); err != nil {
			config.Log("unable to retrieve maintenance items", err)
			return nil, err
		}
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		config.Log("unable to retrieve maintenance items", err)
		return nil, err
	}

	return data, nil
}

// CreateMaintenancePlan ...
func CreateMaintenancePlan(user string, draftMaintenancePlan *model.MaintenancePlan) (*model.MaintenancePlan, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("failed to connect to the database", err)
		return nil, err
	}
	defer db.Close()

	// Retrieve selected status
	selectedStatusDetails, err := RetrieveStatusDetails(user, draftMaintenancePlan.Status)
	if err != nil {
		config.Log("error retrieving status details", err)
		return nil, err
	}
	if selectedStatusDetails == nil {
		config.Log("unable to find selected status", errors.New("unable to find selected status"))
		return nil, errors.New("unable to find selected status")
	}

	parsedCreatorID, err := uuid.Parse(draftMaintenancePlan.CreatedBy)
	if err != nil {
		config.Log("Error parsing creator id", err)
		return nil, err
	}

	sqlStr := `INSERT INTO community.maintenance_plan(name, description, color, status, plan_type, plan_due, location, created_by, created_at, updated_by, updated_at, sharable_groups)
	VALUES ($1, $2, $3, $4, $5, $6, POINT($7, $8), $9, $10, $11, $12, $13)
	RETURNING id, name, description, color, status, plan_type, plan_due, created_at, created_by, updated_at, updated_by, sharable_groups;`

	tx, err := db.Begin()
	if err != nil {
		config.Log("error starting transaction", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	config.Log("SqlStr: %s", nil, sqlStr)
	row := tx.QueryRow(
		sqlStr,
		draftMaintenancePlan.Name,
		draftMaintenancePlan.Description,
		draftMaintenancePlan.Color,
		selectedStatusDetails.ID,
		draftMaintenancePlan.PlanType,
		draftMaintenancePlan.PlanDue,
		draftMaintenancePlan.Location.Lon,
		draftMaintenancePlan.Location.Lat,
		parsedCreatorID,
		time.Now(),
		parsedCreatorID,
		time.Now(),
		pq.Array([]uuid.UUID{parsedCreatorID}),
	)

	err = row.Scan(
		&draftMaintenancePlan.ID,
		&draftMaintenancePlan.Name,
		&draftMaintenancePlan.Description,
		&draftMaintenancePlan.Color,
		&draftMaintenancePlan.Status,
		&draftMaintenancePlan.PlanType,
		&draftMaintenancePlan.PlanDue,
		&draftMaintenancePlan.CreatedAt,
		&draftMaintenancePlan.CreatedBy,
		&draftMaintenancePlan.UpdatedAt,
		&draftMaintenancePlan.UpdatedBy,
		pq.Array(&draftMaintenancePlan.SharableGroups),
	)
	if err != nil {
		config.Log("unable to scan result", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		config.Log("unable to commit transaction", err)
		return nil, err
	}

	draftMaintenancePlan.StatusName = selectedStatusDetails.Name
	draftMaintenancePlan.StatusDescription = selectedStatusDetails.Description

	return draftMaintenancePlan, nil
}

// UpdateMaintenancePlan ...
func UpdateMaintenancePlan(user string, draftMaintenancePlan *model.MaintenancePlan) (*model.MaintenancePlan, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	// retrieve selected status
	selectedStatusDetails, err := RetrieveStatusDetails(user, draftMaintenancePlan.Status)
	if err != nil {
		config.Log("unable to view selected status details", err)
		return nil, err
	}
	if selectedStatusDetails == nil {
		config.Log("unable to find selected status", errors.New("unable to find selected status"))
		return nil, errors.New("unable to find selected status")
	}

	sqlStr := `UPDATE community.maintenance_plan 
    SET 
    name = $2,
    description = $3,
	color = $4,
	status = $5,
	plan_type = $6,
	plan_due = $7,
	location = POINT($8, $9),
    updated_by = $10,
    updated_at = $11,
	sharable_groups = $12
    WHERE id = $1
    RETURNING id, name, description, color, status, plan_type, plan_due, created_at, created_by, updated_at, updated_by, sharable_groups;
`
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		config.Log("unable to begin transaction", err)
		return nil, err
	}

	parsedUpdatorID, err := uuid.Parse(draftMaintenancePlan.UpdatedBy)
	if err != nil {
		tx.Rollback()
		config.Log("unable to parse selected user id", err)
		return nil, err
	}

	var updatedMaintenancePlan model.MaintenancePlan

	config.Log("SqlStr: %s", nil, sqlStr)
	row := tx.QueryRow(sqlStr,
		draftMaintenancePlan.ID,
		draftMaintenancePlan.Name,
		draftMaintenancePlan.Description,
		draftMaintenancePlan.Color,
		selectedStatusDetails.ID,
		draftMaintenancePlan.PlanType,
		draftMaintenancePlan.PlanDue,
		draftMaintenancePlan.Location.Lon,
		draftMaintenancePlan.Location.Lat,
		parsedUpdatorID,
		time.Now(),
		pq.Array(draftMaintenancePlan.SharableGroups),
	)

	err = row.Scan(
		&updatedMaintenancePlan.ID,
		&updatedMaintenancePlan.Name,
		&updatedMaintenancePlan.Description,
		&updatedMaintenancePlan.Color,
		&updatedMaintenancePlan.Status,
		&updatedMaintenancePlan.PlanType,
		&updatedMaintenancePlan.PlanDue,
		&updatedMaintenancePlan.CreatedAt,
		&updatedMaintenancePlan.CreatedBy,
		&updatedMaintenancePlan.UpdatedAt,
		&updatedMaintenancePlan.UpdatedBy,
		pq.Array(&updatedMaintenancePlan.SharableGroups),
	)

	if err != nil {
		tx.Rollback()
		config.Log("unable to scan and parse selected maintenance details", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		config.Log("unable to commit selected transaction", err)
		return nil, err
	}

	updatedMaintenancePlan.Status = selectedStatusDetails.ID.String()
	updatedMaintenancePlan.StatusName = selectedStatusDetails.Name
	updatedMaintenancePlan.StatusDescription = selectedStatusDetails.Description
	updatedMaintenancePlan.Location.Lat = draftMaintenancePlan.Location.Lat
	updatedMaintenancePlan.Location.Lon = draftMaintenancePlan.Location.Lon

	return &updatedMaintenancePlan, nil
}

// RemoveMaintenancePlan ...
func RemoveMaintenancePlan(user string, planID string) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.maintenance_plan WHERE id=$1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	_, err = db.Exec(sqlStr, planID)
	if err != nil {
		config.Log("unable to delete selected maintenance_plan", err)
		return err
	}
	return nil
}

// AddAssetToMaintenancePlan ...
func AddAssetToMaintenancePlan(user string, draftMaintenanceItemRequest *model.MaintenanceItemRequest) ([]model.MaintenanceItemResponse, error) {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `INSERT INTO community.maintenance_item(maintenance_plan_id, item_id, created_by, created_at, updated_by, updated_at, sharable_groups)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	tx, err := db.Begin()
	if err != nil {
		config.Log("error starting transaction", err)
		return nil, err
	}

	config.Log("SqlStr: %s", nil, sqlStr)
	currentTime := time.Now()
	for _, assetID := range draftMaintenanceItemRequest.AssetIDs {
		_, err := tx.Exec(
			sqlStr,
			draftMaintenanceItemRequest.ID,
			assetID,
			draftMaintenanceItemRequest.UserID,
			currentTime,
			draftMaintenanceItemRequest.UserID,
			currentTime,
			pq.Array(draftMaintenanceItemRequest.Collaborators),
		)
		if err != nil {
			tx.Rollback()
			config.Log("Error executing query", err)
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		config.Log("error committing transaction", err)
		return nil, err
	}

	sqlStr = `SELECT 
		mi.id,
		mi.maintenance_plan_id,
		mi.item_id,
		i.name,
		i.description,
		i.price,
		i.quantity,
		i.location,
		mi.created_by,
		COALESCE(cp.username, cp.full_name, cp.email_address, 'Anonymous') as creator,
		mi.created_at,
		mi.updated_by,
		COALESCE(up.username, up.full_name, cp.email_address, 'Anonymous') as updator,
		mi.updated_at,
		mi.sharable_groups
	FROM community.maintenance_item mi
	LEFT JOIN community.inventory i ON mi.item_id = i.id
	LEFT JOIN community.profiles cp ON mi.created_by = cp.id
	LEFT JOIN community.profiles up ON mi.updated_by = up.id
	WHERE $1::UUID = ANY(mi.sharable_groups) AND mi.maintenance_plan_id = $2
	ORDER BY mi.updated_at DESC;`

	config.Log("SqlStr: %s", nil, sqlStr)
	rows, err := db.Query(sqlStr, draftMaintenanceItemRequest.UserID, draftMaintenanceItemRequest.ID)
	if err != nil {
		config.Log("unable to retrieve maintenance items", err)
		return nil, err
	}
	defer rows.Close()

	var data []model.MaintenanceItemResponse
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.MaintenanceItemResponse
		if err := rows.Scan(&ec.ID, &ec.MaintenancePlanID, &ec.ItemID, &ec.Name, &ec.Description, &ec.Price, &ec.Quantity, &ec.Location, &ec.CreatedBy, &ec.Creator, &ec.CreatedAt, &ec.UpdatedBy, &ec.Updator, &ec.UpdatedAt, &sharableGroups); err != nil {
			config.Log("unable to retrieve maintenance items", err)
			return nil, err
		}
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		config.Log("unable to retrieve maintenance items", err)
		return nil, err
	}

	return data, nil
}

// RemoveAssetAssociationFromMaintenancePlan ...
func RemoveAssetAssociationFromMaintenancePlan(user string, draftMaintenancePlan *model.MaintenanceItemRequest) error {
	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.maintenance_item
               WHERE maintenance_plan_id = $1 AND id = ANY($2);`

	tx, err := db.Begin()
	if err != nil {
		config.Log("Error starting transaction", err)
		return err
	}

	config.Log("SqlStr: %s", nil, sqlStr)
	_, err = tx.Exec(sqlStr, draftMaintenancePlan.ID, pq.Array(draftMaintenancePlan.AssetIDs))
	if err != nil {
		tx.Rollback()
		config.Log("Error executing delete query", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		config.Log("Error committing transaction", err)
		return err
	}

	return nil
}

// UpdateMaintenancePlanImage ...
func UpdateMaintenancePlanImage(user string, userID string, maintenanceID string, imageURL string) (bool, error) {

	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return false, err
	}
	defer db.Close()

	sqlStr := `UPDATE community.maintenance_plan mp
	SET associated_image_url = $1,
		updated_at = $4,
		updated_by = $2
		WHERE $2::UUID = ANY(mp.sharable_groups) 
		AND mp.id = $3
	RETURNING mp.id;`

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to start trasanction with selected db pool", err)
		return false, err
	}

	config.Log("SqlStr: %s", nil, sqlStr)
	var updatedMaintenancePlanID string
	err = tx.QueryRow(sqlStr, imageURL, userID, maintenanceID, time.Now()).Scan(&updatedMaintenancePlanID)
	if err != nil {
		config.Log("unable to update maintenance plan id", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		config.Log("unable to commit", err)
		return false, err
	}

	return true, nil
}

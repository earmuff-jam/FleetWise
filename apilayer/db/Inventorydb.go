package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mohit2530/communityCare/model"
)

// RetrieveAllInventoriesForUser ...
func RetrieveAllInventoriesForUser(user string, userID string) ([]model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	data, err := retrieveAllInventoryDetailsForUser(tx, userID)
	if err != nil {
		log.Printf("unable to retrieve all inventories details for user. error: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return make([]model.Inventory, 0), nil
	}

	return data, nil
}

func retrieveAllInventoryDetailsForUser(tx *sql.Tx, userID string) ([]model.Inventory, error) {
	sqlStr := `SELECT
    inv.id,
    inv.name,
    inv.description,
    inv.price,
    inv.status,
    inv.barcode,
    inv.sku,
    inv.quantity,
	inv.bought_at,
    inv.location,
    inv.storage_location_id,
	inv.is_returnable,
	inv.return_location,
	inv.max_weight,
	inv.min_weight,
	inv.max_height,
	inv.min_height,
	inv.associated_image_url,
    inv.created_by,
    COALESCE(cp.username, cp.full_name, cp.email_address) AS creator_name,
    inv.created_at,
    inv.updated_by,
    COALESCE(up.username, up.full_name, up.email_address) AS updater_name,
    inv.updated_at,
	inv.sharable_groups
FROM
    community.inventory inv
LEFT JOIN community.profiles cp ON inv.created_by = cp.id
LEFT JOIN community.profiles up ON inv.updated_by = up.id
WHERE
   inv.created_by = $1
ORDER BY
   inv.updated_at  DESC;
	`
	rows, err := tx.Query(sqlStr, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.Inventory

	for rows.Next() {
		var inventory model.Inventory

		var returnLocation sql.NullString
		var maxWeight sql.NullString
		var minWeight sql.NullString
		var maxHeight sql.NullString
		var minHeight sql.NullString
		var associatedImageURL sql.NullString

		if err := rows.Scan(
			&inventory.ID,
			&inventory.Name,
			&inventory.Description,
			&inventory.Price,
			&inventory.Status,
			&inventory.Barcode,
			&inventory.SKU,
			&inventory.Quantity,
			&inventory.BoughtAt,
			&inventory.Location,
			&inventory.StorageLocationID,
			&inventory.IsReturnable,
			&returnLocation,
			&maxWeight,
			&minWeight,
			&maxHeight,
			&minHeight,
			&associatedImageURL,
			&inventory.CreatedBy,
			&inventory.CreatorName,
			&inventory.CreatedAt,
			&inventory.UpdatedBy,
			&inventory.UpdaterName,
			&inventory.UpdatedAt,
			pq.Array(&inventory.SharableGroups),
		); err != nil {
			return nil, err
		}

		if returnLocation.Valid {
			inventory.ReturnLocation = returnLocation.String
		}

		if maxWeight.Valid {
			inventory.MaxWeight = maxWeight.String
		}

		if minWeight.Valid {
			inventory.MinWeight = minWeight.String
		}

		if maxHeight.Valid {
			inventory.MaxHeight = maxHeight.String
		}

		if minHeight.Valid {
			inventory.MinHeight = minHeight.String
		}

		if associatedImageURL.Valid {
			inventory.AssociatedImageURL = associatedImageURL.String
		}

		data = append(data, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

// RetrieveSelectedInv ...
func RetrieveSelectedInv(user string, userID string, invID string) (*model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	data, err := retrieveSelectedInv(tx, userID, invID)
	if err != nil {
		log.Printf("unable to retrieve all inventories details for user. error: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateAsset ...
func UpdateAsset(user string, userID string, draftUpdateAssetCols model.UpdateAssetColumn) (*model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	columnName := draftUpdateAssetCols.ColumnName
	if !isValidColumnName(columnName) {
		tx.Rollback()
		return nil, fmt.Errorf("invalid column name: %s", columnName)
	}

	sqlStr := fmt.Sprintf(`UPDATE 
	community.inventory inv
		SET %s = $1,
			updated_at = $4,
			updated_by = $2
			WHERE $2::UUID = ANY(inv.sharable_groups) 
			AND inv.id = $3
		RETURNING inv.id;`, columnName)

	var updatedInvID string
	err = tx.QueryRow(sqlStr, draftUpdateAssetCols.InputColumn, userID, draftUpdateAssetCols.AssetID, time.Now()).Scan(&updatedInvID)
	if err != nil {
		log.Printf("unable to update asset id. error: %+v", err)
		return nil, err
	}

	data, err := retrieveSelectedInv(tx, userID, updatedInvID)
	if err != nil {
		log.Printf("unable to retrieve asset details. error: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateAssetImage ...
func UpdateAssetImage(user string, userID string, assetID string, assetImageURL string) (bool, error) {

	db, err := SetupDB(user)
	if err != nil {
		return false, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return false, err
	}

	sqlStr := `UPDATE community.inventory inv
		SET associated_image_url = $1,
			updated_at = $4,
			updated_by = $2
			WHERE $2::UUID = ANY(inv.sharable_groups) 
			AND inv.id = $3
		RETURNING inv.id;`

	var updatedInvID string
	err = tx.QueryRow(sqlStr, assetImageURL, userID, assetID, time.Now()).Scan(&updatedInvID)
	if err != nil {
		log.Printf("unable to update asset id. error: %+v", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("unable to commit. error: %+v", err)
		return false, err
	}
	return true, nil
}

func isValidColumnName(columnName string) bool {
	validColumns := map[string]bool{
		"price":    true,
		"quantity": true,
	}
	return validColumns[columnName]
}

// retrieveSelectedInv...
func retrieveSelectedInv(tx *sql.Tx, userID string, invID string) (*model.Inventory, error) {
	sqlStr := `SELECT
    inv.id,
    inv.name,
    inv.description,
    inv.price,
    inv.status,
    inv.barcode,
    inv.sku,
    inv.quantity,
	inv.bought_at,
    inv.location,
    inv.storage_location_id,
	inv.is_returnable,
	inv.return_location,
	inv.return_datetime,
	inv.return_notes,
	inv.max_weight,
	inv.min_weight,
	inv.max_height,
	inv.min_height,
    inv.created_by,
    COALESCE(cp.username, cp.full_name, cp.email_address) AS creator_name,
    inv.created_at,
    inv.updated_by,
    COALESCE(up.username, up.full_name, up.email_address) AS updater_name,
    inv.updated_at
FROM community.inventory inv
LEFT JOIN community.profiles cp ON inv.created_by = cp.id
LEFT JOIN community.profiles up ON inv.updated_by = up.id
WHERE $1::UUID = ANY(inv.sharable_groups) AND inv.id = $2
ORDER BY inv.updated_at DESC;
	`
	row := tx.QueryRow(sqlStr, userID, invID)

	var inventory model.Inventory
	var returnLocation sql.NullString
	var returnNotes sql.NullString
	var maxWeight sql.NullString
	var minWeight sql.NullString
	var maxHeight sql.NullString
	var minHeight sql.NullString

	err := row.Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Description,
		&inventory.Price,
		&inventory.Status,
		&inventory.Barcode,
		&inventory.SKU,
		&inventory.Quantity,
		&inventory.BoughtAt,
		&inventory.Location,
		&inventory.StorageLocationID,
		&inventory.IsReturnable,
		&returnLocation,
		&inventory.ReturnDateTime,
		&returnNotes,
		&maxWeight,
		&minWeight,
		&maxHeight,
		&minHeight,
		&inventory.CreatedBy,
		&inventory.CreatorName,
		&inventory.CreatedAt,
		&inventory.UpdatedBy,
		&inventory.UpdaterName,
		&inventory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if returnLocation.Valid {
		inventory.ReturnLocation = returnLocation.String
	}
	if returnNotes.Valid {
		inventory.ReturnNotes = returnNotes.String
	}
	if maxWeight.Valid {
		inventory.MaxWeight = maxWeight.String
	}
	if minWeight.Valid {
		inventory.MinWeight = minWeight.String
	}
	if maxHeight.Valid {
		inventory.MaxHeight = maxHeight.String
	}
	if minHeight.Valid {
		inventory.MinHeight = minHeight.String
	}

	return &inventory, nil
}

// AddInventoryInBulk ...
func AddInventoryInBulk(user string, userID string, draftInventoryList model.InventoryListRequest) ([]model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable setup database connection. error: %+v", err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	for _, v := range draftInventoryList.InventoryList {

		// storage location is unique key in the database.
		// storage location can be shared across inventories and items that are stored in events.
		parsedStorageLocationID, err := uuid.Parse(v.Location)
		if err != nil {
			// if the location is not a uuid type, then it should resemble a new storage location
			emptyLocationID := ""
			err := addNewStorageLocation(user, v.Location, userID, &emptyLocationID)
			if err != nil {
				log.Printf("unable to retrieve selected location id. error: %+v", err)
				tx.Rollback()
				return nil, err
			}
			parsedStorageLocationID, err = uuid.Parse(emptyLocationID)
			if err != nil {
				log.Printf("unable to parse the found location id. error: %+v", err)
				tx.Rollback()
				return nil, err
			}
			v.StorageLocationID = emptyLocationID
		}

		parsedCreatedByUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("unable to parse the creator id. error: %+v", err)
			tx.Rollback()
			return nil, err
		}

		sqlStr := `INSERT INTO community.inventory (
			name, 
			description, 
			price, 
			status, 
			barcode, 
			sku, 
			quantity, 
			bought_at, 
			location, 
			storage_location_id, 
			created_by, 
			created_at, 
			updated_by, 
			updated_at,
			sharable_groups
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`

		_, err = tx.Exec(
			sqlStr,
			v.Name,
			v.Description,
			v.Price,
			v.Status,
			v.Barcode,
			v.SKU,
			v.Quantity,
			v.BoughtAt,
			v.Location,
			parsedStorageLocationID,
			parsedCreatedByUUID,
			time.Now(),
			parsedCreatedByUUID,
			time.Now(),
			pq.Array(v.SharableGroups),
		)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

	}

	if err := tx.Commit(); err != nil {
		log.Printf("unable to process trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	tx, err = db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	resp, err := retrieveAllInventoryDetailsForUser(tx, userID)
	if err != nil {
		log.Printf("unable to retrieve all inventories for selected user. error: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("unable to process trasanction with selected db pool. error: %+v", err)
		return nil, err
	}

	draftInventoryList.InventoryList = resp
	return draftInventoryList.InventoryList, nil
}

// AddInventory ...
func AddInventory(user string, userID string, draftInventory model.Inventory) (*model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to start the db. error: %+v", err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start tx. error: %+v", err)
		return nil, err
	}

	// if UUID is not present, add new storage location
	parsedStorageLocationID, err := uuid.Parse(draftInventory.Location)
	if err != nil {
		emptyLocationID := ""
		err := addNewStorageLocation(user, draftInventory.Location, draftInventory.CreatedBy, &emptyLocationID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		parsedStorageLocationID, err = uuid.Parse(emptyLocationID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		draftInventory.StorageLocationID = emptyLocationID
	}

	// if UUID is present, retrieve selected storage location
	sqlStr := `SELECT location FROM community.storage_locations sl WHERE sl.id=$1;`
	err = tx.QueryRow(sqlStr, parsedStorageLocationID).Scan(&draftInventory.Location)
	if err != nil {
		log.Printf("unable to retrieve selected location from storage location id. error: %+v", err)
		tx.Rollback()
		return nil, err
	}
	parsedCreatedByUUID, err := uuid.Parse(draftInventory.CreatedBy)
	if err != nil {
		log.Printf("unable to parse creator userID. error: %+v", err)
		tx.Rollback()
		return nil, err
	}

	currentTimestamp := time.Now()
	draftInventory.CreatedAt = currentTimestamp
	draftInventory.UpdatedAt = currentTimestamp

	sqlStr = `INSERT INTO community.inventory (name,
		description,
		price,
		status,
		barcode,
		sku,
		quantity,
		bought_at,
		location,
		storage_location_id,
		is_returnable,
		return_location,
		return_datetime,
		return_notes,
		max_weight,
		min_weight,
		max_height,
		min_height,
		created_by,
		created_at,
		updated_by,
		updated_at,
		sharable_groups)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
RETURNING id;`

	err = tx.QueryRow(
		sqlStr,
		draftInventory.Name,
		draftInventory.Description,
		draftInventory.Price,
		draftInventory.Status,
		draftInventory.Barcode,
		draftInventory.SKU,
		draftInventory.Quantity,
		draftInventory.BoughtAt,
		draftInventory.Location,
		parsedStorageLocationID,
		draftInventory.IsReturnable,
		draftInventory.ReturnLocation,
		draftInventory.ReturnDateTime,
		draftInventory.ReturnNotes,
		draftInventory.MaxWeight,
		draftInventory.MinWeight,
		draftInventory.MaxHeight,
		draftInventory.MinHeight,
		parsedCreatedByUUID,
		draftInventory.CreatedAt,
		parsedCreatedByUUID,
		draftInventory.UpdatedAt,
		pq.Array([]uuid.UUID{parsedCreatedByUUID}),
	).Scan(&draftInventory.ID)

	if err != nil {
		log.Printf("unable to add selected inventory. error: %+v", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	sqlGetUpdatedInventory := `
	SELECT
		inv.id,
		inv.name,
		inv.description,
		inv.price,
		inv.status,
		inv.barcode,
		inv.sku,
		inv.quantity,
		inv.bought_at,
		inv.location,
		inv.storage_location_id,
		inv.is_returnable,
		inv.return_location,
		inv.return_datetime,
		inv.return_notes,
		inv.max_weight,
		inv.min_weight,
		inv.max_height,
		inv.min_height,
		inv.created_at,
		inv.created_by,
		coalesce (cp.full_name, cp.username, cp.email_address) as creator_name,
		inv.updated_at,
		inv.updated_by,
		coalesce (up.full_name, up.username, up.email_address)  as updater_name
	FROM
		community.inventory inv
	LEFT JOIN community.storage_locations sl on sl.id = inv.storage_location_id 
	LEFT JOIN community.profiles cp on cp.id  = inv.created_by
	LEFT JOIN community.profiles up on up.id  = inv.updated_by
	WHERE inv.id = $1;
`

	row := db.QueryRow(sqlGetUpdatedInventory, draftInventory.ID)

	updatedInventory := model.Inventory{}
	var returnNotes sql.NullString

	err = row.Scan(
		&updatedInventory.ID,
		&updatedInventory.Name,
		&updatedInventory.Description,
		&updatedInventory.Price,
		&updatedInventory.Status,
		&updatedInventory.Barcode,
		&updatedInventory.SKU,
		&updatedInventory.Quantity,
		&updatedInventory.BoughtAt,
		&updatedInventory.Location,
		&updatedInventory.StorageLocationID,
		&updatedInventory.IsReturnable,
		&updatedInventory.ReturnLocation,
		&updatedInventory.ReturnDateTime,
		&returnNotes,
		&updatedInventory.MaxWeight,
		&updatedInventory.MinWeight,
		&updatedInventory.MaxHeight,
		&updatedInventory.MinHeight,
		&updatedInventory.CreatedAt,
		&updatedInventory.CreatedBy,
		&updatedInventory.CreatorName,
		&updatedInventory.UpdatedAt,
		&updatedInventory.UpdatedBy,
		&updatedInventory.UpdaterName,
	)
	if err != nil {
		return nil, err
	}

	if returnNotes.Valid {
		updatedInventory.ReturnNotes = returnNotes.String
	}

	return &updatedInventory, nil
}

// UpdateInventory ...
func UpdateInventory(user string, userID string, draftInventory model.Inventory) (*model.Inventory, error) {

	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to start the db. error: %+v", err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start tx. error: %+v", err)
		return nil, err
	}

	// if UUID is not present, add new storage location
	parsedStorageLocationID, err := uuid.Parse(draftInventory.Location)
	if err != nil {
		emptyLocationID := ""
		err := addNewStorageLocation(user, draftInventory.Location, draftInventory.CreatedBy, &emptyLocationID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		parsedStorageLocationID, err = uuid.Parse(emptyLocationID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		draftInventory.StorageLocationID = emptyLocationID
	}

	// if UUID is present, retrieve selected storage location
	sqlStr := `SELECT location FROM community.storage_locations sl WHERE sl.id=$1;`
	err = tx.QueryRow(sqlStr, parsedStorageLocationID).Scan(&draftInventory.Location)

	if err != nil {
		log.Printf("unable to retrieve selected location from storage location id. error: %+v", err)
		tx.Rollback()
		return nil, err
	}
	parsedCreatedByUUID, err := uuid.Parse(draftInventory.CreatedBy)
	if err != nil {
		log.Printf("unable to parse creator userID. error: %+v", err)
		tx.Rollback()
		return nil, err
	}

	if !draftInventory.IsReturnable {
		draftInventory.ReturnLocation = ""
		draftInventory.ReturnDateTime = nil
	}

	sqlStr = `UPDATE community.inventory inv
	SET name = $2,
		description = $3,
		price = $4,
		status = $5,
		barcode = $6,
		sku = $7,
		quantity = $8,
		bought_at = $9,
		location = $10,
		storage_location_id = $11,
		is_returnable = $12,
		return_location = $13,
		return_datetime = $14,
		return_notes = $15,
		max_weight = $16,
		min_weight = $17,
		max_height = $18,
		min_height = $19,
		created_by = $20,
		created_at = $21,
		updated_by = $22,
		updated_at = $23
	WHERE inv.id = $1
	RETURNING id;`

	err = tx.QueryRow(
		sqlStr,
		draftInventory.ID,
		draftInventory.Name,
		draftInventory.Description,
		draftInventory.Price,
		draftInventory.Status,
		draftInventory.Barcode,
		draftInventory.SKU,
		draftInventory.Quantity,
		draftInventory.BoughtAt,
		draftInventory.Location,
		parsedStorageLocationID,
		draftInventory.IsReturnable,
		draftInventory.ReturnLocation,
		draftInventory.ReturnDateTime,
		draftInventory.ReturnNotes,
		draftInventory.MaxWeight,
		draftInventory.MinWeight,
		draftInventory.MaxHeight,
		draftInventory.MinHeight,
		parsedCreatedByUUID,
		draftInventory.CreatedAt,
		parsedCreatedByUUID,
		time.Now(),
	).Scan(&draftInventory.ID)

	if err != nil {
		log.Printf("unable to update selected inventory. error: %+v", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// new tx to bring fresh values
	tx, err = db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sqlGetUpdatedInventory := `
	SELECT
		inv.id,
		inv.name,
		inv.description,
		inv.price,
		inv.status,
		inv.barcode,
		inv.sku,
		inv.quantity,
		inv.bought_at,
		inv.location,
		inv.storage_location_id,
		inv.is_returnable,
		inv.return_location,
		inv.return_datetime,
		inv.return_notes,
		inv.max_weight,
		inv.min_weight,
		inv.max_height,
		inv.min_height,
		inv.created_at,
		inv.created_by,
		coalesce (cp.full_name, cp.username, cp.email_address) as creator_name,
		inv.updated_at,
		inv.updated_by,
		coalesce (up.full_name, up.username, up.email_address)  as updater_name
	FROM
		community.inventory inv
	LEFT JOIN community.storage_locations sl on sl.id = inv.storage_location_id 
	LEFT JOIN community.profiles cp on cp.id  = inv.created_by
	LEFT JOIN community.profiles up on up.id  = inv.updated_by
	WHERE inv.id = $1;
`

	row := tx.QueryRow(sqlGetUpdatedInventory, draftInventory.ID)

	updatedInventory := model.Inventory{}
	var returnNotes sql.NullString

	err = row.Scan(
		&updatedInventory.ID,
		&updatedInventory.Name,
		&updatedInventory.Description,
		&updatedInventory.Price,
		&updatedInventory.Status,
		&updatedInventory.Barcode,
		&updatedInventory.SKU,
		&updatedInventory.Quantity,
		&updatedInventory.BoughtAt,
		&updatedInventory.Location,
		&updatedInventory.StorageLocationID,
		&updatedInventory.IsReturnable,
		&updatedInventory.ReturnLocation,
		&updatedInventory.ReturnDateTime,
		&returnNotes,
		&updatedInventory.MaxWeight,
		&updatedInventory.MinWeight,
		&updatedInventory.MaxHeight,
		&updatedInventory.MinHeight,
		&updatedInventory.CreatedAt,
		&updatedInventory.CreatedBy,
		&updatedInventory.CreatorName,
		&updatedInventory.UpdatedAt,
		&updatedInventory.UpdatedBy,
		&updatedInventory.UpdaterName,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if returnNotes.Valid {
		updatedInventory.ReturnNotes = returnNotes.String
	}

	// Return the updated inventory object
	return &updatedInventory, nil
}

// DeleteInventory ...
func DeleteInventory(user string, userID string, pruneInventoriesIDs []string) ([]string, error) {

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.inventory WHERE id = ANY($1)`
	_, err = db.Exec(sqlStr, pq.Array(pruneInventoriesIDs))
	if err != nil {
		log.Printf("unable to delete selected inventories: %v", pruneInventoriesIDs)
		return nil, err
	}
	return pruneInventoriesIDs, nil
}

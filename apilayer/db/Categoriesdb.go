package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RetrieveAllCategories ...
func RetrieveAllCategories(user string, userID string, limit int) ([]model.Category, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
	c.id,
	c.name,
	c.description,
	s.id,
	s.name AS status_name,
	s.description AS status_description,
	c.color, 
	c.location[0] AS lon,
	c.location[1] AS lat,
	c.created_at,
	c.created_by,
	COALESCE(cp.full_name, cp.username, cp.email_address) AS creator_name, 
	c.updated_at,
	c.updated_by,
	COALESCE(up.full_name, up.username, up.email_address)  AS updater_name,
	c.sharable_groups
	FROM community.category c
	LEFT JOIN community.statuses s on s.id = c.status
	LEFT JOIN community.profiles cp on cp.id = c.created_by
	LEFT JOIN community.profiles up on up.id = c.updated_by
	WHERE $1::UUID = ANY(c.sharable_groups)
	ORDER BY c.updated_at DESC
	LIMIT $2;`

	rows, err := db.Query(sqlStr, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.Category

	var lon, lat sql.NullFloat64
	var categoryID sql.NullString
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.Category
		if err := rows.Scan(&categoryID, &ec.Name, &ec.Description, &ec.Status, &ec.StatusName, &ec.StatusDescription,
			&ec.Color, &lon, &lat, &ec.CreatedAt, &ec.CreatedBy, &ec.Creator, &ec.UpdatedAt, &ec.UpdatedBy, &ec.Updator, &sharableGroups); err != nil {
			return nil, err
		}

		if lon.Valid && lat.Valid {
			ec.Location = model.Location{
				Lon: lon.Float64,
				Lat: lat.Float64,
			}
		}

		content, _, _, err := FetchImage(categoryID.String)
		if err != nil {
			if err.Error() == "NoSuchKey" {
				log.Printf("cannot find the selected document. error: %+v", err)
			}
		}

		ec.Image = content

		ec.ID = categoryID.String
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// RetrieveAllCategoryItems ...
func RetrieveAllCategoryItems(user string, userID string, categoryID string, limit int) ([]model.CategoryItemResponse, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
		ci.id,
		ci.category_id,
		ci.item_id,
		i.name,
		i.description,
		i.price,
		i.quantity,
		i.location,
		ci.created_by,
		COALESCE(cp.username, cp.full_name, cp.email_address, 'Anonymous') as creator,
		ci.created_at,
		ci.updated_by,
		COALESCE(up.username, up.full_name, cp.email_address, 'Anonymous') as updator,
		ci.updated_at,
		ci.sharable_groups
	FROM community.category_item ci
	LEFT JOIN community.inventory i ON ci.item_id = i.id
	LEFT JOIN community.profiles cp ON ci.created_by = cp.id
	LEFT JOIN community.profiles up ON ci.updated_by = up.id
	WHERE $1::UUID = ANY(ci.sharable_groups) AND ci.category_id = $2
	ORDER BY ci.updated_at DESC FETCH FIRST $3 ROWS ONLY;`

	rows, err := db.Query(sqlStr, userID, categoryID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.CategoryItemResponse
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.CategoryItemResponse
		if err := rows.Scan(&ec.ID, &ec.CategoryID, &ec.ItemID, &ec.Name, &ec.Description, &ec.Price, &ec.Quantity, &ec.Location, &ec.CreatedBy, &ec.Creator, &ec.CreatedAt, &ec.UpdatedBy, &ec.Updator, &ec.UpdatedAt, &sharableGroups); err != nil {
			return nil, err
		}
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// RetrieveCategory ...
func RetrieveCategory(user string, userID string, categoryID string) (model.Category, error) {
	category, err := retrieveCategoryByID(user, userID, categoryID)
	if err != nil {
		log.Printf("unable to retrieve selected category. Error: %+v", err)
		return model.Category{}, err
	}
	return category, nil
}

// CreateCategory ...
func CreateCategory(user string, draftCategory *model.Category) (*model.Category, error) {
	db, err := SetupDB(user)
	if err != nil {
		log.Printf("Failed to connect to the database: %+v", err)
		return nil, err
	}
	defer db.Close()

	// retrieve selected status
	selectedStatusDetails, err := RetrieveStatusDetails(user, draftCategory.Status)
	if err != nil {
		log.Printf("error retrieving status details: %+v", err)
		return nil, err
	}
	if selectedStatusDetails == nil {
		return nil, errors.New("unable to find selected status")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %+v", err)
		return nil, err
	}

	sqlStr := `
	INSERT INTO community.category(name, description, color, status, location, created_by, created_at, updated_by, updated_at, sharable_groups
	) VALUES($1, $2, $3, $4, POINT($5, $6), $7, $8, $9, $10, $11)
	RETURNING id;`

	row := tx.QueryRow(
		sqlStr,
		draftCategory.Name,
		draftCategory.Description,
		draftCategory.Color,
		selectedStatusDetails.ID,
		draftCategory.Location.Lon,
		draftCategory.Location.Lat,
		draftCategory.CreatedBy,
		time.Now(),
		draftCategory.UpdatedBy,
		time.Now(),
		pq.Array(draftCategory.SharableGroups),
	)

	var selectedCategoryID string
	err = row.Scan(
		&selectedCategoryID,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Error scanning result: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %+v", err)
		return nil, err
	}

	selectedCategory, err := retrieveCategoryByID(user, draftCategory.CreatedBy, selectedCategoryID)
	if err != nil {
		log.Printf("unable to retrieve selected category. Error: %+v", err)
		return nil, err
	}

	selectedCategory.StatusName = selectedStatusDetails.Name
	selectedCategory.StatusDescription = selectedStatusDetails.Description

	return &selectedCategory, nil
}

// UpdateCategory ...
func UpdateCategory(user string, draftCategory *model.Category) (*model.Category, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectedStatusDetails, err := RetrieveStatusDetails(user, draftCategory.Status)
	if err != nil {
		return nil, err
	}
	if selectedStatusDetails == nil {
		return nil, errors.New("unable to find selected status")
	}

	sqlStr := `UPDATE community.category 
    SET 
    name = $2,
    description = $3,
	color = $4,
	status = $5,
	location = POINT($6, $7),
    updated_by = $8,
    updated_at = $9,
	sharable_groups = $10
    WHERE id = $1
    RETURNING id;`

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	parsedUpdatorID, err := uuid.Parse(draftCategory.UpdatedBy)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	row := tx.QueryRow(sqlStr,
		draftCategory.ID,
		draftCategory.Name,
		draftCategory.Description,
		draftCategory.Color,
		selectedStatusDetails.ID,
		draftCategory.Location.Lon,
		draftCategory.Location.Lat,
		parsedUpdatorID,
		time.Now(),
		pq.Array(draftCategory.SharableGroups),
	)

	var selectedCategoryID string
	err = row.Scan(
		&selectedCategoryID,
	)

	if err != nil {
		tx.Rollback()
		log.Printf("unable to retrieve selected category. Error: %+v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	selectedCategory, err := retrieveCategoryByID(user, draftCategory.UpdatedBy, selectedCategoryID)

	if err != nil {
		log.Printf("unable to retrieve selected category. Error: %+v", err)
		return nil, err
	}

	selectedCategory.Status = selectedStatusDetails.ID.String()
	selectedCategory.StatusName = selectedStatusDetails.Name
	selectedCategory.StatusDescription = selectedStatusDetails.Description
	selectedCategory.Location.Lat = draftCategory.Location.Lat
	selectedCategory.Location.Lon = draftCategory.Location.Lon

	return &selectedCategory, nil
}

// RemoveCategory ...
func RemoveCategory(user string, categoryID string) error {
	db, err := SetupDB(user)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.category WHERE id=$1;`
	_, err = db.Exec(sqlStr, categoryID)
	if err != nil {
		log.Printf("unable to delete selected category. error: %+v", err)
		return err
	}
	return nil
}

// retrieves the selected category by ID
func retrieveCategoryByID(user string, userID string, categoryID string) (model.Category, error) {
	db, err := SetupDB(user)
	if err != nil {
		return model.Category{}, err
	}
	defer db.Close()

	sqlStr := `SELECT 
	c.id,
	c.name,
	c.description,
	s.id,
	s.name AS status_name,
	s.description AS status_description,
	c.color,
	c.location[0] AS lon,
	c.location[1] AS lat,
	c.created_at,
	c.created_by,
	COALESCE(cp.full_name, cp.username, cp.email_address) AS creator, 
	c.updated_at,
	c.updated_by,
	COALESCE(up.full_name, up.username, up.email_address) AS updator,
	c.sharable_groups
	FROM community.category c
	LEFT JOIN community.statuses s on s.id = c.status
	LEFT JOIN community.profiles cp on cp.id = c.created_by
	LEFT JOIN community.profiles up on up.id = c.updated_by
	WHERE c.id = $2 AND $1::UUID = ANY(c.sharable_groups)
	ORDER BY c.updated_at DESC;`

	row := db.QueryRow(sqlStr, userID, categoryID)
	selectedCategory := model.Category{}

	var lon, lat sql.NullFloat64

	err = row.Scan(
		&selectedCategory.ID,
		&selectedCategory.Name,
		&selectedCategory.Description,
		&selectedCategory.Status,
		&selectedCategory.StatusName,
		&selectedCategory.StatusDescription,
		&selectedCategory.Color,
		&lon,
		&lat,
		&selectedCategory.CreatedAt,
		&selectedCategory.CreatedBy,
		&selectedCategory.Creator,
		&selectedCategory.UpdatedAt,
		&selectedCategory.UpdatedBy,
		&selectedCategory.Updator,
		pq.Array(&selectedCategory.SharableGroups),
	)

	if lon.Valid && lat.Valid {
		selectedCategory.Location = model.Location{
			Lon: lon.Float64,
			Lat: lat.Float64,
		}
	}

	if err != nil {
		return model.Category{}, err
	}

	return selectedCategory, nil
}

// AddAssetToCategory ...
func AddAssetToCategory(user string, draftCategory *model.CategoryItemRequest) ([]model.CategoryItemResponse, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `INSERT INTO community.category_item(category_id, item_id, created_by, created_at, updated_by, updated_at, sharable_groups)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error starting transaction: %+v", err)
		return nil, err
	}

	currentTime := time.Now()
	for _, assetID := range draftCategory.AssetIDs {
		_, err := tx.Exec(
			sqlStr,
			draftCategory.ID,
			assetID,
			draftCategory.UserID,
			currentTime,
			draftCategory.UserID,
			currentTime,
			pq.Array(draftCategory.Collaborators),
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error executing query: %+v", err)
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error committing transaction: %+v", err)
		return nil, err
	}

	sqlStr = `SELECT 
		ci.id,
		ci.category_id,
		ci.item_id,
		i.name,
		i.description,
		i.price,
		i.quantity,
		i.location,
		ci.created_by,
		COALESCE(cp.username, cp.full_name, cp.email_address, 'Anonymous') as creator,
		ci.created_at,
		ci.updated_by,
		COALESCE(up.username, up.full_name, cp.email_address, 'Anonymous') as updator,
		ci.updated_at,
		ci.sharable_groups
	FROM community.category_item ci
	LEFT JOIN community.inventory i ON ci.item_id = i.id
	LEFT JOIN community.profiles cp ON ci.created_by = cp.id
	LEFT JOIN community.profiles up ON ci.updated_by = up.id
	WHERE $1::UUID = ANY(ci.sharable_groups) AND ci.category_id = $2
	ORDER BY ci.updated_at DESC;`

	rows, err := db.Query(sqlStr, draftCategory.UserID, draftCategory.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.CategoryItemResponse
	var sharableGroups pq.StringArray

	for rows.Next() {
		var ec model.CategoryItemResponse
		if err := rows.Scan(&ec.ID, &ec.CategoryID, &ec.ItemID, &ec.Name, &ec.Description, &ec.Price, &ec.Quantity, &ec.Location, &ec.CreatedBy, &ec.Creator, &ec.CreatedAt, &ec.UpdatedBy, &ec.Updator, &ec.UpdatedAt, &sharableGroups); err != nil {
			return nil, err
		}
		ec.SharableGroups = sharableGroups
		data = append(data, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// RemoveAssetAssociationFromCategory ...
func RemoveAssetAssociationFromCategory(user string, draftCategory *model.CategoryItemRequest) error {
	db, err := SetupDB(user)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr := `DELETE FROM community.category_item
               WHERE category_id = $1 AND id = ANY($2);`

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %+v", err)
		return err
	}

	_, err = tx.Exec(sqlStr, draftCategory.ID, pq.Array(draftCategory.AssetIDs))
	if err != nil {
		tx.Rollback()
		log.Printf("Error executing delete query: %+v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %+v", err)
		return err
	}

	return nil
}

// UpdateCategoryImage ...
func UpdateCategoryImage(user string, userID string, categoryID string, imageURL string) (bool, error) {

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
	sqlStr := `UPDATE community.category c
		SET associated_image_url = $1,
			updated_at = $4,
			updated_by = $2
			WHERE $2::UUID = ANY(c.sharable_groups) 
			AND c.id = $3
		RETURNING c.id;`

	var updatedCategoryID string
	err = tx.QueryRow(sqlStr, imageURL, userID, categoryID, time.Now()).Scan(&updatedCategoryID)
	if err != nil {
		log.Printf("unable to update category id. error: %+v", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("unable to commit. error: %+v", err)
		return false, err
	}

	return true, nil
}

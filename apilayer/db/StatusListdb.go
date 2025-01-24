package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
)

// RetrieveStatusDetails ...
func RetrieveStatusDetails(user string, statusID string) (*model.StatusList, error) {

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT id, name, description FROM community.statuses s WHERE s.name=$1;`
	row := db.QueryRow(sqlStr, statusID)

	var selectedStatusID, selectedStatusName, selectedStatusDescription string
	err = row.Scan(&selectedStatusID, &selectedStatusName, &selectedStatusDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("unable to find selected status")
		}
		log.Printf("invalid status selected. error: %+v", err)
		return nil, err
	}

	parsedSelectedStatusID, err := uuid.Parse(selectedStatusID)
	if err != nil {
		log.Printf("error in parsing selected status. error: %+v", err)
		return nil, err
	}

	selectedStatus := model.StatusList{
		ID:          parsedSelectedStatusID,
		Name:        selectedStatusName,
		Description: selectedStatusDescription,
	}

	return &selectedStatus, nil
}

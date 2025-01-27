package db

import (
	"database/sql"
	"errors"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
)

// RetrieveStatusDetails ...
func RetrieveStatusDetails(user string, statusID string) (*model.StatusList, error) {

	db, err := SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT id, name, description FROM community.statuses s WHERE s.name=$1;`

	config.Log("SqlStr: %s", nil, sqlStr)
	row := db.QueryRow(sqlStr, statusID)

	var selectedStatusID, selectedStatusName, selectedStatusDescription string
	err = row.Scan(&selectedStatusID, &selectedStatusName, &selectedStatusDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			config.Log("no rows during scan", errors.New("unable to find selected status"))
			return nil, errors.New("unable to find selected status")
		}
		config.Log("invalid status selected", err)
		return nil, err
	}

	parsedSelectedStatusID, err := uuid.Parse(selectedStatusID)
	if err != nil {
		config.Log("error in parsing selected status", err)
		return nil, err
	}

	selectedStatus := model.StatusList{
		ID:          parsedSelectedStatusID,
		Name:        selectedStatusName,
		Description: selectedStatusDescription,
	}

	return &selectedStatus, nil
}

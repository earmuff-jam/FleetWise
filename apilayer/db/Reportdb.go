package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
)

// RetrieveReports ...
func RetrieveReports(user string, userID uuid.UUID, sinceDateTime string, includeOverdueAssets string) ([]model.Report, error) {

	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to setup db. error: %+v", err)
		return nil, err
	}
	defer db.Close()

	var additionalWhereClause string

	if includeOverdueAssets == "true" {
		additionalWhereClause = "OR inv.return_datetime >= $2::TIMESTAMP WITH TIME ZONE"
	}

	draftSqlStr := `
	WITH filtered_inventory AS (
			SELECT 
			inv.id,
			inv.price
			FROM community.inventory inv
			WHERE 
				(inv.updated_at >= $2::TIMESTAMP WITH TIME ZONE %s)
				AND $1::UUID = ANY(inv.sharable_groups)
		)
	SELECT 
		(SELECT SUM(price) FROM filtered_inventory) AS total_cost,
		(SELECT SUM(inv.price) 
			FROM filtered_inventory inv 
			LEFT JOIN ( SELECT DISTINCT item_id FROM community.category_item ) ci ON ci.item_id = inv.id 
			WHERE ci.item_id IS NOT NULL
		) AS total_category_items_cost;`

	parsedSqlStr := fmt.Sprintf(draftSqlStr, additionalWhereClause)
	var reports []model.Report
	rows, err := db.Query(parsedSqlStr, userID, sinceDateTime)
	if err != nil {
		log.Printf("unable to retrieve report details. error: %+v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var draftReport model.Report
		var totalValuationDraft sql.NullFloat64
		var totalCategoryItemsCostDraft sql.NullFloat64

		if err := rows.Scan(&totalValuationDraft, &totalCategoryItemsCostDraft); err != nil {
			log.Printf("unable to scan reports. error: %+v", err)
			return nil, err
		}
		if totalValuationDraft.Valid {
			draftReport.ItemValuation = totalValuationDraft.Float64
		}
		if totalCategoryItemsCostDraft.Valid {
			draftReport.TotalCategoryItemsCost = totalCategoryItemsCostDraft.Float64
		}

		parsedTime, _ := time.Parse(time.RFC3339, sinceDateTime)
		draftReport.SelectedTimeRange = parsedTime
		reports = append(reports, draftReport)
	}

	if err := rows.Err(); err != nil {
		log.Printf("unable to retrieve report. error: %+v", err)
		return nil, err
	}
	return reports, nil
}

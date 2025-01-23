package db

import (
	"log"

	"github.com/earmuff-jam/fleetwise/model"
	"github.com/lib/pq"
)

// RetrieveAssetsAndSummary...
func RetrieveAssetsAndSummary(user string, userID string) (model.AssetsAndSummaryResponse, error) {
	db, err := SetupDB(user)
	if err != nil {
		log.Printf("unable to setup db. error: %+v", err)
		return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
	}
	defer db.Close()

	sqlStr := `SELECT 
    id, 
    name, 
    type, 
    returnTime, 
    price, 
	items,
    created_at,
    updated_at, 
    sharable_groups 
		FROM (
			SELECT
				c.id,
				c.name,
				'C' as type,
				NULL::timestamp with time zone as returnTime,
				0 as price,
				array_agg(COALESCE(i.name, '')) AS items,
				c.created_at,
				c.updated_at,
				c.sharable_groups
			FROM community.category c
			LEFT JOIN community.category_item ci ON c.id = ci.category_id
			LEFT JOIN community.inventory i ON i.id = ci.item_id
			GROUP BY c.id
			UNION
			SELECT
				mp.id,
				mp.name,
				'M' as type,
				NULL::timestamp with time zone as returnTime,
				0 as price,
				array_agg(COALESCE(i.name, '')) AS items,
				mp.created_at,
				mp.updated_at,
				mp.sharable_groups
			FROM community.maintenance_plan mp
			LEFT JOIN community.maintenance_item mi ON mp.id = mi.maintenance_plan_id
			LEFT JOIN community.inventory i ON i.id = mi.item_id
			GROUP BY mp.id
			UNION
			SELECT 
				i.id,
				i.name,
				'A' as type, 
				i.return_datetime AS returnTime,
				i.price,
				ARRAY[]::TEXT[] AS items,
				i.created_at,
				i.updated_at,
				i.sharable_groups
			FROM community.inventory i
			WHERE i.is_returnable = TRUE AND i.return_datetime IS NOT NULL AND i.return_datetime < CURRENT_TIMESTAMP
		) AS combined
		WHERE $1::UUID = ANY(sharable_groups)
		ORDER BY type, updated_at DESC;`

	rows, err := db.Query(sqlStr, userID)
	if err != nil {
		log.Printf("unable to retrieve asset summary details. error: %+v", err)
		return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
	}
	defer rows.Close()

	var response model.AssetsAndSummaryResponse
	var data []model.AssetSummary

	for rows.Next() {
		var as model.AssetSummary
		var returnDateTime pq.NullTime
		var items pq.StringArray
		var sharableGroups pq.StringArray
		if err := rows.Scan(&as.ID, &as.Name, &as.Type, &returnDateTime, &as.Price, &items, &as.CreatedAt, &as.UpdatedAt, &sharableGroups); err != nil {
			log.Printf("unable to retrieve asset summary details. error: %+v", err)
			return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
		}

		if returnDateTime.Valid {
			as.ReturnTime = returnDateTime.Time
		}

		as.Items = items

		as.SharableGroups = sharableGroups
		data = append(data, as)
	}

	if err := rows.Err(); err != nil {
		log.Printf("unable to retrieve asset summary details. error: %+v", err)
		return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
	}

	sqlStr = `SELECT 
	i.name, 
	i.description, 
	i.price, 
	i.quantity, 
	i.bought_at
		FROM community.inventory i 
	WHERE 
		$1::UUID = ANY(sharable_groups) ORDER BY updated_at DESC;`

	var inventories []model.Inventory

	rows, err = db.Query(sqlStr, userID)
	if err != nil {
		log.Printf("unable to retrieve asset details in summary. error: %+v", err)
		return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
	}

	for rows.Next() {
		var i model.Inventory
		if err := rows.Scan(&i.Name, &i.Description, &i.Price, &i.Quantity, &i.BoughtAt); err != nil {
			log.Printf("unable to retrieve asset details in summary. error: %+v", err)
			return model.AssetsAndSummaryResponse{AssetSummaryList: []model.AssetSummary{}, AssetList: []model.Inventory{}}, err
		}
		inventories = append(inventories, i)
	}

	response.AssetSummaryList = data
	response.AssetList = inventories

	return response, nil
}

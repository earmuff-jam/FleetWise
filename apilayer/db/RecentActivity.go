package db

import (
	"time"

	"github.com/earmuff-jam/fleetwise/model"
	"github.com/lib/pq"
)

// RetrieveRecentActivities ...
func RetrieveRecentActivities(user string, userID string, limit int, untilDate string) ([]model.RecentActivity, error) {

	retrieveRecentActivitiesBaseSQL := `SELECT ra.id,
       ra.activity_id,
       ra.type,
	   ra.title,
       ra.custom_action,
       ra.created_by,
       COALESCE(cp.username, cp.full_name, cp.email_address) as creator,
       ra.created_at,
       ra.updated_by,
       COALESCE(up.username, up.full_name, up.email_address) as updator,
       ra.updated_at,
       ra.sharable_groups
	FROM community.recent_activities ra
	LEFT JOIN community.profiles cp on ra.created_by = cp.id
	LEFT JOIN community.profiles up on ra.updated_by = up.id
	WHERE $1::UUID = ANY(ra.sharable_groups)`

	var endDate time.Time
	var err error

	if untilDate != "" {
		endDate, err = time.Parse(time.RFC3339, untilDate)
		if err != nil {
			return nil, err
		}
		retrieveRecentActivitiesBaseSQL += " AND ra.updated_at >= $3"
	}

	retrieveRecentActivitiesBaseSQL += " ORDER BY ra.updated_at DESC LIMIT $2;"

	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	args := []interface{}{userID, limit}
	if untilDate != "" {
		args = append(args, endDate)
	}

	rows, err := db.Query(retrieveRecentActivitiesBaseSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recentActivities []model.RecentActivity

	for rows.Next() {
		var activity model.RecentActivity

		if err := rows.Scan(&activity.ID, &activity.ActivityID, &activity.Type, &activity.Title, &activity.CustomAction, &activity.CreatedBy, &activity.Creator, &activity.CreatedAt, &activity.UpdatedBy, &activity.Updator, &activity.UpdatedAt, pq.Array(&activity.SharableGroups)); err != nil {
			return nil, err
		}

		recentActivities = append(recentActivities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recentActivities, nil
}

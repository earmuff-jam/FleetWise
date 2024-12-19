package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Profile ...
// swagger:model Profile
//
// Profile object of the user.
type Profile struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	AvatarUrl    string    `json:"avatar_url,omitempty"`
	EmailAddress string    `json:"email_address"`
	PhoneNumber  string    `json:"phone_number"`
	AboutMe      string    `json:"about_me"`
	OnlineStatus bool      `json:"online_status"`
	Appearance   bool      `json:"appearance"`
	GridView     bool      `json:"grid_view"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty"`
	Creator      string    `json:"creator,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	UpdatedBy    string    `json:"updated_by,omitempty"`
	Updator      string    `json:"updator,omitempty"`
}

// Validate ...
//
// used to validate the profile fields
func (p *Profile) Validate() error {
	if len(p.Username) == 0 || len(p.FullName) == 0 || len(p.EmailAddress) == 0 || len(p.PhoneNumber) == 0 {
		return errors.New("empty required field")
	}
	return nil
}

// ProfileStats ...
// swagger:model ProfileStats
//
// ProfileStats object that returns the count of assets, categories and maintenance plans
type ProfileStats struct {
	TotalCategories       int `json:"totalCategories"`
	TotalMaintenancePlans int `json:"totalMaintenancePlans"`
	TotalAssets           int `json:"totalAssets"`
}

// MaintenanceAlertNotifications ...
// swagger:model MaintenanceAlertNotifications
//
// MaintenanceAlertNotifications object that returns the notifications alert for the maintenance plans that are within 7 days of being due
type MaintenanceAlertNotifications struct {
	ID             string    `json:"maintenance_plan_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	PlanDue        time.Time `json:"plan_due"`
	IsRead         bool      `json:"is_read"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	UpdatedBy      string    `json:"updated_by,omitempty"`
	Updator        string    `json:"updator,omitempty"`
	SharableGroups []string  `json:"sharable_groups"`
}

// MaintenanceAlertNotificationRequest ...
// swagger:model MaintenanceAlertNotificationRequest
//
// MaintenanceAlertNotificationRequest object that is used to update the db based on user selection.
type MaintenanceAlertNotificationRequest struct {
	ID        string    `json:"maintenance_plan_id"`
	IsRead    bool      `json:"is_read"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}

// Location ...
// swagger:model Location
//
// Location object. Used to store the lat and lon details
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Notes ...
// swagger:model Note
//
// Notes object. Used to store the note details
type Note struct {
	ID                string      `json:"noteID"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Status            string      `json:"status"`
	StatusName        string      `json:"status_name"`
	StatusDescription string      `json:"status_description"`
	Color             string      `json:"color"`
	CompletionDate    *time.Time  `json:"completionDate,omitempty"`
	Location          Location    `json:"location,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	CreatedBy         string      `json:"created_by"`
	Creator           string      `json:"creator,omitempty"`
	UpdatedAt         time.Time   `json:"updated_at"`
	UpdatedBy         string      `json:"updated_by"`
	Updator           string      `json:"updator,omitempty"`
	SharableGroups    []uuid.UUID `json:"sharable_groups"`
}

// FavouriteItem ...
// swagger:model FavouriteItem
//
// FavouriteItem object. Used to store the list of categories or maintenance plans marked as favourite by the selected user. This does not contain the audit columns because they are not used directly. Created user should be the only reference to populate this model.
type FavouriteItem struct {
	ID                    string      `json:"id"`
	CategoryID            string      `json:"category_id,omitempty"`
	CategoryName          string      `json:"category_name,omitempty"`
	CategoryStatus        string      `json:"category_status,omitempty"`
	MaintenancePlanID     string      `json:"maintenance_plan_id,omitempty"`
	MaintenancePlanName   string      `json:"maintenance_plan_name,omitempty"`
	MaintenancePlanStatus string      `json:"maintenance_plan_status,omitempty"`
	CreatedAt             time.Time   `json:"created_at"`
	CreatedBy             string      `json:"created_by"`
	Creator               string      `json:"creator"`
	UpdatedAt             time.Time   `json:"updated_at"`
	UpdatedBy             string      `json:"updated_by"`
	Updator               string      `json:"updator"`
	SharableGroups        []uuid.UUID `json:"sharable_groups"`
}

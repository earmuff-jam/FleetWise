package model

import (
	"time"
)

// Inventory ...
// swagger:model Inventory
//
// Inventory is the selected row for each inventory
type Inventory struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Price              float64    `json:"price"`
	Status             string     `json:"status"`
	Barcode            string     `json:"barcode"`
	SKU                string     `json:"sku"`
	Color              string     `json:"color,omitempty"`
	Quantity           int        `json:"quantity"`
	Location           string     `json:"location"`
	StorageLocationID  string     `json:"storage_location_id"`
	IsReturnable       bool       `json:"is_returnable"`
	ReturnLocation     string     `json:"return_location"`
	ReturnDateTime     *time.Time `json:"return_datetime,omitempty"`
	ReturnNotes        string     `json:"return_notes,omitempty"`
	MaxWeight          int        `json:"max_weight,omitempty"`
	MinWeight          int        `json:"min_weight,omitempty"`
	MaxHeight          int        `json:"max_height,omitempty"`
	MinHeight          int        `json:"min_height,omitempty"`
	AssociatedImageURL string     `json:"associated_image_url"`
	Image              []byte     `json:"image,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          string     `json:"created_by"`
	CreatorName        string     `json:"creator_name"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          string     `json:"updated_by"`
	UpdaterName        string     `json:"updator"`
	BoughtAt           string     `json:"bought_at"`
	SharableGroups     []string   `json:"sharable_groups"`
}

// RawInventory ...
// swagger:model RawInventory
//
// RawInventory is used to derieve the single row from bulk uploaded excel file
type RawInventory struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	Quantity         int64   `json:"quantity"`
	StorageLocation  string  `json:"Storage Location"`
	Color            string  `json:"color"`
	SKU              string  `json:"sku"`
	Barcode          string  `json:"barcode"`
	PurchaseLocation string  `json:"Purchase Location"`
	MaximumWeight    int     `json:"Maximum Weight"`
	MinimumWeight    int     `json:"Minimum Weight"`
	MaximumHeight    int     `json:"Maximum Height"`
	MinimumHeight    int     `json:"Minimum Height"`
}

// InventoryListRequest ...
// swagger:model InventoryListRequest
//
// InventoryListRequest is used to formulate the bulk download for selected inventory
type InventoryListRequest struct {
	InventoryList []Inventory
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

// InventoryItemToUpdate ...
// swagger:model InventoryItemToUpdate
//
// InventoryItemToUpdate is the object that needs to be updated when the client is updating the inventory list of their own personal account. A user can update a certain limit of columns
type InventoryItemToUpdate struct {
	Column string `json:"column"`
	Value  string `json:"value"`
	ID     string `json:"id"`
	UserID string `json:"userID"`
}

// UpdateAssetColumn
// swagger:model UpdateAssetColumn
//
// UpdateAssetColumn struct is used to update a specific inventory item.
type UpdateAssetColumn struct {
	AssetID     string `json:"assetID"`
	ColumnName  string `json:"columnName"`
	InputColumn string `json:"inputColumn"`
}

// StorageLocation ...
// swagger:model StorageLocation
//
// Storage Location is the location where the item has been stored
type StorageLocation struct {
	ID             string    `json:"id"`
	Location       string    `json:"location"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedBy      string    `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	SharableGroups []string  `json:"sharable_groups"`
}

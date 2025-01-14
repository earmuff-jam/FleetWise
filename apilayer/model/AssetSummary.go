package model

import "time"

// AssetSummary ...
// swagger:model AssetSummary
//
// AssetSummary instance to resemble a single assetSummary
type AssetSummary struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	ReturnTime     time.Time `json:"returntime"`
	Price          float64   `json:"price"`
	Items          []string  `json:"items"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	SharableGroups []string  `json:"sharable_groups"`
}

// AssetsAndSummaryResponse ...
// swagger:model AssetsAndSummaryResponse
//
// AssetsAndSummaryResponse struct that contains the Assets and Summary Responses
type AssetsAndSummaryResponse struct {
	AssetSummaryList []AssetSummary
	AssetList        []Inventory
}

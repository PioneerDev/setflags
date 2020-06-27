package models

import (
	"set-flags/schemas"
	"time"

	"github.com/fox-one/mixin-sdk"
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Asset entity
type Asset struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Symbol   string     `json:"symbol"`
	PriceUSD float64    `json:"price_usd"`
	Balance  float64    `json:"balance"`
	PaidAt   *time.Time `json:"paid_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FindAssetByID find asset by id
func FindAssetByID(assetID uuid.UUID) (asset Asset) {
	db.Where("id = ?", assetID.String()).First(&asset)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (a *Asset) BeforeCreate(scope *gorm.Scope) error {
	// uuid_, _ := uuid.NewV4()
	// scope.SetColumn("ID", uuid_)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

// BeforeUpdate set field updateAt
func (a *Asset) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// FindAsset find asset by id
func FindAsset(assetID uuid.UUID) *Asset {
	var assets []Asset
	db.Find(&assets)
	for _, a := range assets {
		if a.ID == assetID {
			return &a
		}
	}
	return nil
}

// UpsertAsset UpsertAsset
func UpsertAsset(asset *mixin.Asset) {

	assetID, _ := uuid.FromString(asset.AssetID)
	priceUSD, _ := asset.PriceUsd.Float64()
	balance, _ := asset.Balance.Float64()

	dbAsset := Asset{
		ID:       assetID,
		Symbol:   asset.Symbol,
		PriceUSD: priceUSD,
		Balance:  balance,
	}

	if err := db.Where("id = ?", asset.AssetID).First(&dbAsset).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Create(&dbAsset)
		}
	} else {
		db.Model(&Asset{}).Where("id = ?", asset.AssetID).Updates(map[string]interface{}{
			"symbol":    asset.Symbol,
			"price_usd": priceUSD,
			"balance":   balance,
		})
	}
}

// ReadAssets ReadAssets
func ReadAssets(pageSize, currentPage int) (assets []schemas.AssetSchema, count int) {
	var dbAssets []Asset
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Order("created_at desc").Find(&dbAssets)
	db.Model(&Asset{}).Count(&count)

	for _, a := range dbAssets {

		assets = append(assets, schemas.AssetSchema{
			AssetID:  a.ID.String(),
			Symbol:   a.Symbol,
			PriceUSD: a.PriceUSD,
		})
	}
	return
}

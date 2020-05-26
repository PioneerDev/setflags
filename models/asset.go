package models

import (
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Asset struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Symbol   string     `json:"symbol"`
	PriceUSD float64    `json:"price_usd"`
	Balance  float64    `json:"balance"`
	PaidAt   *time.Time `json:"paid_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FindAssetByID(assetId string) (asset Asset) {
	db.Where("id = ?", assetId).First(&asset)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (a *Asset) BeforeCreate(scope *gorm.Scope) error {
	uuid_, _ := uuid.NewV4()
	scope.SetColumn("ID", uuid_)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (a *Asset) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func FindAsset(assetId uuid.UUID) *Asset {
	var assets []Asset
	db.Find(&assets)
	for _, a := range assets {
		if a.ID == assetId {
			return &a
		}
	}
	return nil
}

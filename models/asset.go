package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Asset struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Symbol   string     `json:"symbol"`
	PriceUSD float64    `json:"price_usd"`
	Balance  float64    `json:"balance"`
	PaidAt   *time.Time `json:"paid_at"`

	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FindAssetsByID(assetId string) (assets []Asset) {
	db.Where("id = ?", assetId).Find(&assets)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (a *Asset) BeforeCreate(scope *gorm.Scope) error {
	uuid_ := uuid.NewV4()
	scope.SetColumn("ID", uuid_)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (a *Asset) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

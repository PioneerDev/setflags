package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Flag struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	PayerId       string  `json:"payer_id"`
	Task          string  `json:"task"`
	Days          int     `json:"days"`
	MaxWitness    int     `json:"max_witness"`
	AssetId       string  `json:"asset_id"`
	Amount        float64 `json:"amount"`
	TimesAchieved int     `json:"times_achieved"`
	Status        string  `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateFlag(data map[string]interface{}) bool {
	db.Create(&Flag{
		PayerId:       data["payer_id"].(string),
		Task:          data["task"].(string),
		Days:          int(data["days"].(float64)),
		MaxWitness:    int(data["max_witness"].(float64)),
		AssetId:       data["asset_id"].(string),
		Amount:        data["amount"].(float64),
		TimesAchieved: int(data["times_achieved"].(float64)),
		Status:        data["status"].(string),
	})

	return true
}

func GetAllFlags() (flags []Flag) {
	db.Find(&flags)
	return
}

func FindFlagsByUserID(userId string) (flags []Flag) {
	db.Where("payer_id = ?", userId).Find(&flags)
	return
}

func FLagExists(flagId string) bool {
	var count int

	db.Model(&Flag{}).Where("id = ?", flagId).Count(&count)

	return count == 1
}


// BeforeCreate will set a UUID rather than numeric ID.
func (flag *Flag) BeforeCreate(scope *gorm.Scope) error {
	uuid_ := uuid.NewV4()
	scope.SetColumn("ID", uuid_)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (flag *Flag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

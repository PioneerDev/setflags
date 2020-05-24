package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Flag struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	PayerId         uuid.UUID `json:"payer_id"`
	Task            string    `json:"task"`
	Days            int       `json:"days"`
	MaxWitness      int       `json:"max_witness"`
	AssetId         uuid.UUID `json:"asset_id"`
	Amount          float64   `json:"amount"`
	TimesAchieved   int       `json:"times_achieved"`
	Status          string    `json:"status"`
	RemainingDays   int       `json:"remaining_days"`
	RemainingAmount float64   `json:"remaining_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func CreateFlag(data map[string]interface{}) bool {
	db.Create(&Flag{
		PayerId:    data["payer_id"].(uuid.UUID),
		Task:       data["task"].(string),
		Days:       int(data["days"].(int)),
		MaxWitness: int(data["max_witness"].(int)),
		AssetId:    data["asset_id"].(uuid.UUID),
		Amount:     data["amount"].(float64),
		Status:     data["status"].(string),
		// below are derived
		RemainingAmount: data["amount"].(float64),
		RemainingDays:   int(data["days"].(int)),
		TimesAchieved:   0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})

	return true
}

func GetAllFlags() (flags []Flag) {
	db.Order("created_at desc").Find(&flags)
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

func FindFlagByID(flagId string) (flag Flag) {
	db.Where("id = ?", flagId).First(&flag)
	return
}

func UpdateFlagStatus(flagId, status string) bool {
	db.Model(&Flag{}).Where("id = ?", flagId).Update("status", status)
	return true
}

// BeforeCreate will set a UUID rather than numeric ID.
func (flag *Flag) BeforeCreate(scope *gorm.Scope) error {
	uuid_, _ := uuid.NewV4()
	scope.SetColumn("ID", uuid_)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (flag *Flag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (flag *Flag) Witnesses() []*Witness {
	var witnesses []*Witness
	db.Where("flag_id = ?", flag.ID).Find(&witnesses)
	return witnesses
}

func ListActiveFlags(paid bool) []*Flag {
	var flags []*Flag
	if paid {
		db.Where("days > 1 and date_part('day', now() - created_at::date) < days and status='PAID'").Find(&flags)
	} else {
		db.Where("days > 1 and date_part('day', now() - created_at::date) < days and status!='PAID'").Find(&flags)
	}
	return flags
}

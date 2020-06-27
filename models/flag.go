package models

import (
	"set-flags/schemas"
	"time"

	"github.com/fox-one/mixin-sdk"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Flag entity
type Flag struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	PayerID         uuid.UUID `json:"payer_id"`
	PayerName       string    `json:"payer_name"`
	PayerAvatarURL  string    `json:"payer_avatar_url"`
	Task            string    `json:"task"`
	Days            int       `json:"days"`
	MaxWitness      int       `json:"max_witness"`
	AssetID         uuid.UUID `json:"asset_id"`
	Amount          float64   `json:"amount"`
	TimesAchieved   int       `json:"times_achieved"`
	Status          string    `json:"status"`
	RemainingDays   int       `json:"remaining_days"`
	RemainingAmount float64   `json:"remaining_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateFlag create flag
func CreateFlag(flagJSON *schemas.Flag, user *UserSchema) bool {
	db.Create(&Flag{
		PayerID:        flagJSON.PayerID,
		PayerName:      user.FullName,
		PayerAvatarURL: user.AvatarURL,
		Task:           flagJSON.Task,
		Days:           flagJSON.Days,
		MaxWitness:     flagJSON.MaxWitness,
		AssetID:        flagJSON.AssetID,
		Amount:         flagJSON.Amount,
		Status:         "unverified",
		// below are derived
		RemainingAmount: flagJSON.Amount,
		RemainingDays:   flagJSON.Days,
		TimesAchieved:   0,
	})

	return true
}

// GetAllFlags fetch all flags
func GetAllFlags(pageSize, currentPage int) (flags []Flag, count int) {
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Order("created_at desc").Find(&flags)
	db.Model(&Flag{}).Count(&count)
	return
}

// FindFlagsByUserID find current user's flags
func FindFlagsByUserID(userID uuid.UUID, currentPage, pageSize int) (flags []Flag, total int) {
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("payer_id = ?", userID.String()).Find(&flags)
	db.Model(&Flag{}).Where("payer_id = ?", userID.String()).Count(&total)
	return
}

// FlagExists check flag exist
func FlagExists(flagID uuid.UUID) bool {
	var count int

	db.Model(&Flag{}).Where("id = ?", flagID.String()).Count(&count)

	return count == 1
}

// FindFlagByID find flag by it's id
func FindFlagByID(flagID uuid.UUID) (flag Flag) {
	db.Where("id = ?", flagID.String()).First(&flag)
	return
}

// UpdateFlagStatus update flag's status
func UpdateFlagStatus(flagID uuid.UUID, status string) bool {
	db.Model(&Flag{}).Where("id = ?", flagID.String()).Update("status", status)
	return true
}

// UpdateFlagUserInfo update flag's user info
func UpdateFlagUserInfo(user *mixin.Profile) bool {
	db.Model(&Flag{}).Where("payer_id = ?", user.UserID).
		Updates(map[string]interface{}{
			"payer_name":       user.FullName,
			"payer_avatar_url": user.AvatarURL,
		})
	return true
}

// BeforeCreate will set a UUID rather than numeric ID.
func (flag *Flag) BeforeCreate(scope *gorm.Scope) error {
	uuid, _ := uuid.NewV4()
	scope.SetColumn("ID", uuid)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

// BeforeUpdate will set field udpate time.
func (flag *Flag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// Witnesses fetch flag's witness.
func (flag *Flag) Witnesses() []*Witness {
	var witnesses []*Witness
	db.Where("flag_id = ?", flag.ID).Find(&witnesses)
	return witnesses
}

// GetWitnesses fetch the witness of the flag by its ID and page number.
func GetWitnesses(flagID uuid.UUID, pageSize, currentPage int) []*Witness {
	var witnesses []*Witness
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ?", flagID).Find(&witnesses)
	return witnesses
}

// ListActiveFlags fetch active flags
func ListActiveFlags(paid bool) []*Flag {
	var flags []*Flag
	if paid {
		db.Where("days > 1 and date_part('day', now() - created_at::date) < days and status='PAID'").Find(&flags)
	} else {
		db.Where("days > 1 and date_part('day', now() - created_at::date) < days and status!='PAID'").Find(&flags)
	}
	return flags
}

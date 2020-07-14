package models

import (
	"set-flags/schemas"
	"strings"
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
	Symbol          string    `json:"symbol"`
	Amount          float64   `json:"amount"`
	TimesAchieved   int       `json:"times_achieved"`
	Status          string    `json:"status"`
	RemainingDays   int       `json:"remaining_days"`
	RemainingAmount float64   `json:"remaining_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateFlag create flag
func CreateFlag(flagJSON *schemas.FlagSchema, user *UserSchema) uuid.UUID {
	flag := &Flag{
		PayerID:        flagJSON.PayerID,
		PayerName:      user.FullName,
		PayerAvatarURL: user.AvatarURL,
		Task:           flagJSON.Task,
		Days:           flagJSON.Days,
		MaxWitness:     flagJSON.MaxWitness,
		AssetID:        flagJSON.AssetID,
		Symbol:         flagJSON.Symbol,
		Amount:         flagJSON.Amount,
		Status:         strings.ToUpper("unverified"),
		// below are derived
		RemainingAmount: flagJSON.Amount,
		RemainingDays:   flagJSON.Days,
		TimesAchieved:   0,
	}
	db.Create(flag)

	return flag.ID
}

// GetAllFlags fetch all flags
func GetAllFlags(pageSize, currentPage int) (flags []Flag, count int) {
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Order("updated_at desc").Find(&flags)
	db.Model(&Flag{}).Count(&count)
	return
}

// GetFlagsWithVerified update flag status according to witness
func GetFlagsWithVerified(pageSize, currentPage int, userID uuid.UUID) (flagSchemas []schemas.FlagSchema, count int) {
	skip := (currentPage - 1) * pageSize

	var flags []Flag

	db.Model(&Flag{}).Where("status = ?", "PAID").Count(&count)

	// first fetch flags
	db.Offset(skip).Limit(pageSize).Where("status = ?", "PAID").Order("updated_at desc").Find(&flags)

	// then fetch witness according to userID and flagID
	flagIDs := make([]uuid.UUID, len(flags))
	for i := 0; i < len(flags); i++ {
		flagIDs = append(flagIDs, flags[i].ID)
	}

	var witnesses []Witness
	db.Where("flag_id IN (?) and payee_id = ?", flagIDs, userID).Find(&witnesses)

	for _, flag := range flags {
		verified := 0
		for _, w := range witnesses {
			if w.FlagID != w.FlagID {
				continue
			}
			verified = w.Verified
			break
		}

		flagSchemas = append(flagSchemas, schemas.FlagSchema{
			ID:              flag.ID,
			PayerID:         flag.PayerID,
			PayerName:       flag.PayerName,
			PayerAvatarURL:  flag.PayerAvatarURL,
			Task:            flag.Task,
			Days:            flag.Days,
			MaxWitness:      flag.MaxWitness,
			AssetID:         flag.AssetID,
			Symbol:          flag.Symbol,
			Amount:          flag.Amount,
			TimesAchieved:   flag.TimesAchieved,
			Status:          flag.Status,
			RemainingAmount: flag.RemainingAmount,
			RemainingDays:   flag.RemainingDays,
			Verified:        verified,
		})
	}

	return
}

// FindFlagsByUserID find current user's flags
func FindFlagsByUserID(userID uuid.UUID, currentPage, pageSize int) (flagSchemas []schemas.FlagSchema, total int) {
	skip := (currentPage - 1) * pageSize
	var flags []Flag
	db.Offset(skip).Limit(pageSize).Where("payer_id = ?", userID.String()).Order("updated_at desc").Find(&flags)
	db.Model(&Flag{}).Where("payer_id = ?", userID.String()).Count(&total)

	// prepare flagSchema
	for _, flag := range flags {
		flagSchemas = append(flagSchemas, schemas.FlagSchema{
			ID:              flag.ID,
			PayerID:         flag.PayerID,
			PayerName:       flag.PayerName,
			PayerAvatarURL:  flag.PayerAvatarURL,
			Task:            flag.Task,
			Days:            flag.Days,
			MaxWitness:      flag.MaxWitness,
			AssetID:         flag.AssetID,
			Symbol:          flag.Symbol,
			Amount:          flag.Amount,
			TimesAchieved:   flag.TimesAchieved,
			Status:          flag.Status,
			RemainingAmount: flag.RemainingAmount,
			RemainingDays:   flag.RemainingDays,
		})
	}
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
	db.Model(&Flag{}).Where("id = ?", flagID).Update("status", strings.ToUpper(status))
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

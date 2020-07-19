package models

import (
	"set-flags/schemas"
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Witness entity
type Witness struct {
	FlagID        uuid.UUID `json:"flag_id"`
	PayeeID       uuid.UUID `json:"payee_id"`
	Verified      string    `json:"verified"`
	WitnessedTime time.Time `json:"witnessed_time"`
	Period        int       `json:"period"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetWitnessByFlagIDAndPayeeID GetWitnessByFlagIDAndPayeeID
func GetWitnessByFlagIDAndPayeeID(flagID, payeeID uuid.UUID, period int) (w Witness) {
	db.Where("flag_id = ? and payee_id = ? and period = ?", flagID, payeeID, period).Find(&w)
	return
}

// UpsertWitness UpsertWitness
func UpsertWitness(flagID, payeeID uuid.UUID, op string, period int) {

	verified := strings.ToUpper(op)
	witness := Witness{
		FlagID:   flagID,
		PayeeID:  payeeID,
		Verified: verified,
		Period:   period,
	}

	// no found witness, insert witness
	// found, update witness
	db.Where(Witness{FlagID: flagID, PayeeID: payeeID}).Assign(Witness{Verified: verified}).FirstOrCreate(&witness)
}

// GetAllWitnessByFlagID GetAllWitnessByFlagID
func GetAllWitnessByFlagID(flagID uuid.UUID, pageSize, currentPage int) ([]schemas.WitnessSchema, int) {

	var count int

	result := make([]schemas.WitnessSchema, 0, 0)

	var witnesses []*Witness
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ?", flagID).Order("period desc, updated_at desc").Find(&witnesses)

	db.Model(&Witness{}).Where("flag_id = ?", flagID).Count(&count)

	for _, witness := range witnesses {
		var dbUser User
		db.Where("user_id = ?", witness.PayeeID.String()).First(&dbUser)
		result = append(result, schemas.WitnessSchema{
			FlagID:         flagID,
			PayeeID:        witness.PayeeID,
			PayeeName:      dbUser.FullName,
			PayeeAvatarURL: dbUser.AvatarURL,
			// todo
			Symbol:        "BTC",
			Amount:        1.0,
			Verified:      witness.Verified,
			WitnessedTime: witness.WitnessedTime,
			Period:        witness.Period,
		})
	}
	return result, count
}

// GetWitnessWithPeriod GetWitnessWithPeriod
func GetWitnessWithPeriod(flagID uuid.UUID, pageSize, currentPage, period int) ([]schemas.WitnessSchema, int) {

	var count int

	result := make([]schemas.WitnessSchema, 0, 0)

	var witnesses []*Witness
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ? and period = ?", flagID, period).
		Order("period desc, updated_at desc").Find(&witnesses)

	db.Model(&Witness{}).Where("flag_id = ? and period = ?", flagID, period).Count(&count)

	for _, witness := range witnesses {
		var dbUser User
		db.Where("user_id = ?", witness.PayeeID.String()).First(&dbUser)
		result = append(result, schemas.WitnessSchema{
			FlagID:         flagID,
			PayeeID:        witness.PayeeID,
			PayeeName:      dbUser.FullName,
			PayeeAvatarURL: dbUser.AvatarURL,
			Symbol:         "BTC",
			Amount:         1.0,
			Verified:       witness.Verified,
			WitnessedTime:  witness.WitnessedTime,
			Period:         witness.Period,
		})
	}
	return result, count
}

// BeforeCreate will set field CreatedAt.
func (w *Witness) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("WitnessedTime", time.Now())
	return nil
}

// BeforeUpdate will set field update time.
func (w *Witness) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	scope.SetColumn("WitnessedTime", time.Now())
	return nil
}

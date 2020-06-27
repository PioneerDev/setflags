package models

import (
	"set-flags/schemas"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Witness entity
type Witness struct {
	FlagID   uuid.UUID `json:"flag_id"`
	PayeeID  uuid.UUID `json:"payee_id"`
	Verified int       `json:"verified"`
}

// UpsertWitness UpsertWitness
func UpsertWitness(flagID, payeeID uuid.UUID) {

	witness := Witness{
		FlagID:   flagID,
		PayeeID:  payeeID,
		Verified: 1,
	}

	if err := db.Where("flag_id = ? and payee_id = ?", flagID.String(), payeeID.String()).First(&witness).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Create(&witness)
		}
	} else {
		db.Model(&Asset{}).
			Where("flag_id = ? and payee_id = ?", flagID.String(), payeeID.String()).
			UpdateColumn("verified", gorm.Expr("verified + ?", 1))
	}
}

// GetWitnessSchema GetWitnessSchema
func GetWitnessSchema(flagID uuid.UUID, pageSize, currentPage int) ([]schemas.WitnessSchema, int) {

	var count int

	result := make([]schemas.WitnessSchema, 0, 0)

	var witnesses []*Witness
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ?", flagID).Find(&witnesses)

	db.Model(&Witness{}).Where("flag_id = ?", flagID).Count(&count)

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
		})
	}
	return result, count
}

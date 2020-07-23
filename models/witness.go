package models

import (
	"fmt"
	"set-flags/schemas"
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Witness entity
type Witness struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	FlagID        uuid.UUID `json:"flag_id"`
	PayeeID       uuid.UUID `json:"payee_id"`
	AssetID       uuid.UUID `json:"asset_id"`
	Verified      string    `json:"verified"`
	Period        int       `json:"period"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	Symbol        string    `json:"symbol"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	WitnessedTime time.Time `json:"witnessed_time"`
}

// GetWitnessByFlagIDAndPayeeID GetWitnessByFlagIDAndPayeeID
func GetWitnessByFlagIDAndPayeeID(flagID, payeeID uuid.UUID, period int) (w Witness) {
	db.Where("flag_id = ? and payee_id = ? and period = ?", flagID, payeeID, period).Find(&w)
	return
}

// UpsertWitness UpsertWitness
func UpsertWitness(flagID, payeeID, assetID uuid.UUID, op, symbol string, period, maxWitness int) error {
	return upsertWitness(db, flagID, payeeID, assetID, op, symbol, period, maxWitness)
}

func upsertWitness(db *gorm.DB, flagID, payeeID, assetID uuid.UUID, op, symbol string, period, maxWitness int) error {
	verified := strings.ToUpper(op)
	witness := &Witness{
		FlagID:   flagID,
		PayeeID:  payeeID,
		Verified: verified,
		Period:   period,
		AssetID:  assetID,
		Status:   strings.ToUpper("pending"),
		Symbol:   strings.ToUpper(symbol),
		Amount:   0.0,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	var exist, count int

	// check if exist
	if err := tx.Model(&Witness{}).
		Where("flag_id = ? and payee_id = ? and period = ?", flagID, payeeID, period).
		Count(&exist).Error; err != nil {
		tx.Rollback()
		return err
	}

	if exist == 1 {
		tx.Model(&Witness{}).
			Where("flag_id = ? and payee_id = ? and period = ?", flagID, payeeID, period).
			Update("status", strings.ToUpper(verified))
	} else if exist == 0 {
		if err := tx.Model(&Witness{}).Where("flag_id = ? and period = ?", flagID, period).Count(&count).Error; err != nil {
			tx.Rollback()
			return err
		}

		if count >= maxWitness {
			return fmt.Errorf("max witness")
		}

		if err := tx.Create(witness).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
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
			Symbol:         witness.Symbol,
			Amount:         witness.Amount,
			Verified:       witness.Verified,
			WitnessedTime:  witness.WitnessedTime,
			Period:         witness.Period,
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
			Symbol:         witness.Symbol,
			Amount:         witness.Amount,
			Verified:       witness.Verified,
			WitnessedTime:  witness.WitnessedTime,
			Period:         witness.Period,
		})
	}
	return result, count
}

// GetWitnessByFlagIDAndPeriod GetWitnessByFlagIDAndPeriod
func GetWitnessByFlagIDAndPeriod(flagID uuid.UUID, period int, status string) (witnesses []Witness) {
	db.Where("flag_id = ? and period = ? and status = ?", flagID, period, strings.ToUpper(status)).Find(&witnesses)
	return
}

// GetErrorWitnessByFlagID GetErrorWitnessByFlagID
func GetErrorWitnessByFlagID(flagID uuid.UUID, status string) (witnesses []Witness) {
	db.Where("flag_id = ? and status = ?", flagID, strings.ToUpper(status)).Find(&witnesses)
	return
}

// UpdateWitnessStatus update witness's status
func UpdateWitnessStatus(witnessID uuid.UUID, status string, amount float64) {
	// db.Model(&Witness{}).Where("id = ?", witnessID).Update("status", strings.ToUpper(status))
	db.Model(&Witness{}).Where("id = ?", witnessID).Updates(map[string]interface{}{
		"status": strings.ToUpper(status),
		"amount": amount,
	})
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

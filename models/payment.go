package models

import (
	"set-flags/schemas"
	"strings"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Payment entity
type Payment struct {
	TraceID    uuid.UUID `gorm:"type:uuid;primary_key;" json:"trace_id"`
	FlagID     uuid.UUID `json:"flag_id"`
	AssetID    string    `json:"asset_id,omitempty"`
	OpponentID string    `json:"opponent_id,omitempty"`
	Amount     string    `json:"amount,omitempty"`
	Memo       string    `json:"memo,omitempty"`
	Status     string    `json:"status,omitempty"`
	// AddressID  string    `json:"address_id,omitempty"`
	// OpponentKey string    `json:"opponent_key,omitempty"`
}

// CreatePayment create payment
func CreatePayment(payment Payment) {
	payment.Status = strings.ToUpper("pending")
	db.Create(&payment)
}

// ListNoPaidPayment ListNoPaidPayment
func ListNoPaidPayment() (payments []*Payment) {
	db.Where("status = ?", strings.ToUpper("pending")).Find(&payments)
	return
}

// UpdatePaymentStatus UpdatePaymentStatus
func UpdatePaymentStatus(traceID uuid.UUID, status string) {
	db.Model(&Payment{}).Where("trace_id = ?", traceID).Update("status", strings.ToUpper(status))
}

// UpdatePaymentAndFlag update payment and flag status
func UpdatePaymentAndFlag(db *gorm.DB, snapshot schemas.AccountSnapshot) error {
	// 请注意，事务一旦开始，你就应该使用 tx 作为数据库句柄
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&Payment{}).Where("trace_id = ?", snapshot.TraceID).Updates(map[string]interface{}{
		"amount": snapshot.Amount,
		"memo":   snapshot.Memo,
		"status": strings.ToUpper("PAID"),
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	var payment Payment
	db.Where("trace_id = ?", snapshot.TraceID).Select("flag_id").First(&payment)

	err = tx.Model(&Flag{}).Where("id = ?", payment.FlagID).Updates(map[string]interface{}{
		"amount": snapshot.Amount,
		"period": 1,
		"status": strings.ToUpper("PAID"),
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

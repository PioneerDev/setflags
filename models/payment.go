package models

import (
	"strings"

	uuid "github.com/gofrs/uuid"
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

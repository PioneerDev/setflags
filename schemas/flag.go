package schemas

import "github.com/gofrs/uuid"

// FlagSchema validate flag json
type FlagSchema struct {
	ID              uuid.UUID `json:"id"`
	PayerID         uuid.UUID `json:"payer_id"`
	PayerName       string    `json:"payer_name"`
	PayerAvatarURL  string    `json:"payer_avatar_url"`
	Task            string    `json:"task" binding:"required"`
	Days            int       `json:"days" binding:"required"`
	MaxWitness      int       `json:"max_witness" binding:"required"`
	AssetID         uuid.UUID `json:"asset_id" binding:"required"`
	Symbol          string    `json:"symbol" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	TimesAchieved   int       `json:"times_achieved"`
	Status          string    `json:"status"`
	Verified        int       `json:"verified"`
	RemainingDays   int       `json:"remaining_days"`
	RemainingAmount float64   `json:"remaining_amount"`
}

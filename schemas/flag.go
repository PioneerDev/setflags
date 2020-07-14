package schemas

import "github.com/gofrs/uuid"

// FlagSchema validate flag json
type FlagSchema struct {
	ID              uuid.UUID `json:"id,omitempty"`
	PayerID         uuid.UUID `json:"payer_id,omitempty"`
	PayerName       string    `json:"payer_name,omitempty"`
	PayerAvatarURL  string    `json:"payer_avatar_url,omitempty"`
	Task            string    `json:"task,omitempty" binding:"required"`
	Days            int       `json:"days,omitempty" binding:"required"`
	MaxWitness      int       `json:"max_witness,omitempty" binding:"required"`
	AssetID         uuid.UUID `json:"asset_id,omitempty" binding:"required"`
	Symbol          string    `json:"symbol,omitempty" binding:"required"`
	Amount          float64   `json:"amount,omitempty" binding:"required"`
	TimesAchieved   int       `json:"times_achieved,omitempty"`
	Status          string    `json:"status,omitempty"`
	Verified        int       `json:"verified"`
	RemainingDays   int       `json:"remaining_days,omitempty"`
	RemainingAmount float64   `json:"remaining_amount,omitempty"`
}

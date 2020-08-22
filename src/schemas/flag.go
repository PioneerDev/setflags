package schemas

import (
	"github.com/gofrs/uuid"
)

// FlagSchema validate flag json
type FlagSchema struct {
	ID              uuid.UUID `json:"id,omitempty"`
	PayerID         uuid.UUID `json:"payer_id,omitempty"`
	PayerName       string    `json:"payer_name,omitempty"`
	PayerAvatarURL  string    `json:"payer_avatar_url,omitempty"`
	Days            int       `json:"days,omitempty"`
	Task            string    `json:"task,omitempty" binding:"required"`
	MaxWitness      int       `json:"max_witness,omitempty" binding:"required"`
	AssetID         uuid.UUID `json:"asset_id,omitempty" binding:"required"`
	Symbol          string    `json:"symbol,omitempty" binding:"required"`
	Amount          float64   `json:"amount,omitempty" binding:"required"`
	DaysPerPeriod   int       `json:"days_per_period" binding:"required"`
	TotalPeriod     int       `json:"total_period" binding:"required"`
	TimesAchieved   int       `json:"times_achieved"`
	Period          int       `json:"period"`
	Status          string    `json:"status,omitempty"`
	PeriodStatus    string    `json:"period_status"`
	Verified        string    `json:"verified"`
	RemainingDays   int       `json:"remaining_days,omitempty"`
	RemainingAmount float64   `json:"remaining_amount,omitempty"`
}

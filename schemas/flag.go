package schemas

import "github.com/gofrs/uuid"

// Flag validate flag json
type Flag struct {
	PayerID        uuid.UUID `json:"payer_id"`
	PayerName      string    `json:"payer_name"`
	PayerAvatarURL string    `json:"payer_avatar_url"`
	Task           string    `json:"task" binding:"required"`
	Days           int       `json:"days" binding:"required"`
	MaxWitness     int       `json:"max_witness" binding:"required"`
	AssetID        uuid.UUID `json:"asset_id" binding:"required"`
	Amount         float64   `json:"amount" binding:"required"`
	TimesAchieved  int       `json:"times_achieved" binding:"required"`
	Status         string    `json:"status"`
}

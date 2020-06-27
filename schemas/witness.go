package schemas

import (
	"time"

	"github.com/gofrs/uuid"
)

// WitnessSchema WitnessSchema
type WitnessSchema struct {
	FlagID         uuid.UUID `json:"flag_id"`
	PayeeID        uuid.UUID `json:"payee_id"`
	PayeeName      string    `json:"payee_name"`
	PayeeAvatarURL string    `json:"payee_avatar_url"`
	WitnessedTime  time.Time `json:"witnessed_time"`
	Amount         float64   `json:"amount"`
	Symbol         string    `json:"symbol"`
}

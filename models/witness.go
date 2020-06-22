package models

import (
	uuid "github.com/gofrs/uuid"
)

// Witness entity
type Witness struct {
	FlagID   uuid.UUID `json:"flag_id"`
	PayeeID  uuid.UUID `json:"payee_id"`
	Verified int       `json:"verified"`
}

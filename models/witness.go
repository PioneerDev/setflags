package models

import (
	uuid "github.com/gofrs/uuid"
)

type Witness struct {
	FlagId   uuid.UUID `json:"flag_id"`
	PayeeId  uuid.UUID `json:"payee_id"`
	Verified int       `json:"verified"`
}

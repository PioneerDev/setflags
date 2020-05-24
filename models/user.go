package models

import (
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	ID             uuid.UUID `json:"id"`
	IdentityNumber int       `json:"identity_number"`
	FullName       string    `json:"full_name"`
	AvatarUrl      string    `json:"avatar_url"`
}

func FindUser(userId uuid.UUID) *User {
	var users []User
	db.Find(&users)
	for _, u := range users {
		if u.ID == userId {
			return &u
		}
	}
	return nil
}

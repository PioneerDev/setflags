package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
}

type Payer struct {
	gorm.Model

	User User `json:"user"`
	Paid bool `json:"paid"`
}

type Payee struct {
	gorm.Model
	User     User `json:"user"`
	Verified bool `json:"verified"`
}

package models

import (
	"time"

	"github.com/fox-one/mixin-sdk"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// User entity
type User struct {
	UserID         uuid.UUID `json:"id" gorm:"Column:user_id"`
	IdentityNumber string    `json:"identity_number"`
	FullName       string    `json:"full_name"`
	AvatarURL      string    `json:"avatar_url"`
	AccessToken    string    `json:"access_token"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserSchema return to front end
type UserSchema struct {
	UserID         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
}

// FindUser find user by id
func FindUser(userID uuid.UUID) *User {
	var users []User
	db.Find(&users)
	for _, u := range users {
		if u.UserID == userID {
			return &u
		}
	}
	return nil
}

// FindUserByID find user by id
// no return access token to front end
func FindUserByID(userID uuid.UUID) *UserSchema {
	var dbUser User
	db.Where("user_id = ?", userID.String()).First(&dbUser)
	var user UserSchema
	user.UserID = dbUser.UserID.String()
	user.AvatarURL = dbUser.AvatarURL
	user.FullName = dbUser.FullName
	user.IdentityNumber = dbUser.IdentityNumber
	return &user
}

// CreateUser create user
func CreateUser(userProfile *mixin.Profile, accessToken string) bool {

	userID, _ := uuid.FromString(userProfile.UserID)
	db.Create(&User{
		UserID:         userID,
		IdentityNumber: userProfile.IdentityNumber,
		FullName:       userProfile.FullName,
		AvatarURL:      userProfile.AvatarURL,
		AccessToken:    accessToken,
	})

	return true
}

// FindUserToken find user's access token
func FindUserToken(userID string) (string, error) {
	var user User
	db.Where("user_id = ?", userID).First(&user)
	return user.AccessToken, nil
}

// UserExist check user exist.
func UserExist(userID uuid.UUID) bool {
	var count int

	db.Model(&User{}).Where("user_id = ?", userID.String()).Count(&count)

	return count == 1
}

// UpdateUser update user's access token.
func UpdateUser(userProfile *mixin.Profile, accessToken string) {
	db.Model(&User{}).Where("user_id = ?", userProfile.UserID).
		Updates(map[string]interface{}{
			"full_name":    userProfile.FullName,
			"avatar_url":   userProfile.AvatarURL,
			"access_token": accessToken,
		})
}

// BeforeCreate will set field CreatedAt.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

// BeforeUpdate will set field update time.
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

package utils

type UserInfo struct {
	Type           string `json:"type"`
	UserId         string `json:"user_id"`
	Name           string `json:"name"`
	IdentityNumber string `json:"identity_number"`
}

package schemas

// AuthToken token
type AuthToken struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

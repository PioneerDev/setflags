package schemas

// CheckReward user_id flag_id
type CheckReward struct {
	UserID string `uri:"user_id" binding:"required,uuid"`
	FlagID string `uri:"flag_id" binding:"required,uuid"`
}

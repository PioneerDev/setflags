package schemas

import "time"

// AccountSnapshot AccountSnapshot
type AccountSnapshot struct {
	Amount        float64   `json:"amount,omitempty"`
	AssetID       string    `json:"asset_id,omitempty"`
	CounterUserID string    `json:"counter_user_id"`
	AddressID     string    `json:"address_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Memo          string    `json:"memo,omitempty"`
	OpponentID    string    `json:"opponent_id,omitempty"`
	TraceID       string    `json:"trace_id,omitempty"`
	SnapshotID    string    `json:"snapshot_id"`
	Type          string    `json:"type"`
}

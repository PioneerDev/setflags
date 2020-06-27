package schemas

// AssetSchema AssetSchema
type AssetSchema struct {
	AssetID  string  `json:"asset_id"`
	Symbol   string  `json:"symbol"`
	PriceUSD float64 `json:"price_usd"`
	Balance  float64 `json:"balance"`
}

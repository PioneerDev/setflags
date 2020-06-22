package schemas

// Pagination pagination params
type Pagination struct {
	CurrentPage int `form:"current_page"`
	PageSize    int `form:"page_size"`
}

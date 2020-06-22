package schemas

// Header header params
type Header struct {
	XUSERID string `header:"x-user-id" binding:"required,uuid"`
}

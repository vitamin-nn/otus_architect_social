package form

type Feed struct {
	Body string `json:"body" binding:"required"`
}

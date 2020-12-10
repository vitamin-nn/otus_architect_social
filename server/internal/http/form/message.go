package form

type Message struct {
	ToUserID int    `json:"to_user_id" binding:"required"`
	Text     string `json:"body" binding:"required"`
}

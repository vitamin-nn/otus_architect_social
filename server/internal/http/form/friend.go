package form

type Friend struct {
	FriendID int `json:"friend_id" binding:"required"`
}

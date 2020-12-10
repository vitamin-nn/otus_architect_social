package repository

import (
	"context"
	"time"
)

type Message struct {
	ID         int       `json:"id"`
	DialogID   int       `json:"dialog_id"`
	SenderID   int       `json:"sender_user_id"`
	ReceiverID int       `json:"receiver_user_id"`
	SentAt     time.Time `json:"sent_at"`
	Text       string    `json:"body"`
	IsRead     bool      `json:"is_read"`
}

type MessengerRepo interface {
	Create(ctx context.Context, message *Message) (*Message, error)
	GetDialogMessageList(ctx context.Context, user1ID, user2ID, limit, offset int) ([]*Message, error)
}

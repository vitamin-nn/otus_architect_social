package repository

import (
	"context"
	"time"
)

type Feed struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	CreateAt time.Time `json:"sent_at"`
	Body     string    `json:"body"`
}

type FeedRepo interface {
	Create(ctx context.Context, message *Feed) (*Feed, error)
	GetFeedByUserID(ctx context.Context, userID, limit, offset int) ([]*Feed, error)
}

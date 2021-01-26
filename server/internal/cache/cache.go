package cache

import (
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

type Feed interface {
	Add(userID int, f *repository.Feed) error
	GetUserFeed(userID int) ([]*repository.Feed, error)
}

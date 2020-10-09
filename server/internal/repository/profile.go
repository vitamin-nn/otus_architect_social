package repository

import (
	"context"
	"time"
)

type Profile struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Birth        time.Time `json:"birth_date"`
	Sex          string    `json:"sex"`
	Interest     string    `json:"interest"`
	City         string    `json:"city"`
}

type ProfileRepo interface {
	CreateProfile(ctx context.Context, profile *Profile) (*Profile, error)
	UpdateProfile(ctx context.Context, profile *Profile) (*Profile, error)
	GetProfileList(ctx context.Context, limit, offset int) ([]*Profile, error)
	GetFriendsProfileList(ctx context.Context, profileID, limit, offset int) ([]*Profile, error)
	GetProfileByID(ctx context.Context, profileID int) (*Profile, error)
	AddFriend(ctx context.Context, profileID1, profileID2 int) error
	RemoveFriend(ctx context.Context, profileID1, profileID2 int) error
	GetProfileByEmail(ctx context.Context, email string) (*Profile, error)
}

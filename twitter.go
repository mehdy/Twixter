package twixter

import (
	"time"

	"github.com/google/uuid"
)

// TwitterProfile represents a user's profile on Twitter.
type TwitterProfile struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	//Twitter Data
	TwitterID           string
	Name                string
	Username            string
	Location            string
	Bio                 string
	URL                 string
	Email               string
	ProfileBannerURL    string
	ProfileImageURL     string
	Verified            bool
	Protected           bool
	DefaultProfile      bool
	DefaultProfileImage bool
	FollowersCount      int
	FollowingsCount     int
	FavouritesCount     int
	ListedCount         int
	TweetsCount         int
	Entities            map[string]interface{}
	JoinedAt            time.Time
	FollowingsIDs       []string // TwitterID references
	FollowersIDs        []string // TwitterID references
}

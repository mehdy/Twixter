package twixter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// nolint: gosec // no credentials
	twitterTokenURL = "https://api.twitter.com/oauth2/token"
	maxFetchCount   = 200
)

// TwitterProfile represents a user's profile on Twitter.
type TwitterProfile struct {
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

type Twitter struct {
	api *twitter.Client
}

func NewTwitter(config *viper.Viper) *Twitter {
	creds := &clientcredentials.Config{
		ClientID:     config.GetString("twitter.consumer.key"),
		ClientSecret: config.GetString("twitter.consumer.secret"),
		TokenURL:     twitterTokenURL,
	}

	return &Twitter{
		api: twitter.NewClient(creds.Client(context.TODO())),
	}
}

func (t *Twitter) GetProfile(username string) (*TwitterProfile, error) {
	user, resp, err := t.api.Users.Show(&twitter.UserShowParams{ScreenName: username})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile from twitter API: %w", err)
	}
	defer resp.Body.Close()

	return t.toTwitterProfile(*user), nil
}

func (t *Twitter) GetFollowings(username string) []*TwitterProfile {
	skipStatus := true
	includeUserEntities := true

	var cursor int64 = -1
	profiles := []*TwitterProfile{}

	for cursor != 0 {
		following, resp, err := t.api.Friends.List(&twitter.FriendListParams{
			ScreenName:          username,
			Count:               maxFetchCount,
			Cursor:              cursor,
			SkipStatus:          &skipStatus,
			IncludeUserEntities: &includeUserEntities,
		})
		if err != nil {
			fmt.Printf("Failed to fetch followings of %q from twitter API: %s", username, err)
			continue
		}
		defer resp.Body.Close()

		cursor = following.NextCursor

		for _, u := range following.Users {
			profiles = append(profiles, t.toTwitterProfile(u))
		}
	}

	return profiles
}

func (t *Twitter) GetFollowers(username string) []*TwitterProfile {
	skipStatus := true
	includeUserEntities := true

	var cursor int64 = -1
	profiles := []*TwitterProfile{}

	for cursor != 0 {
		follower, resp, err := t.api.Followers.List(&twitter.FollowerListParams{
			ScreenName:          username,
			Count:               maxFetchCount,
			Cursor:              cursor,
			SkipStatus:          &skipStatus,
			IncludeUserEntities: &includeUserEntities,
		})
		if err != nil {
			fmt.Printf("Failed to fetch followers of %q from twitter API: %s", username, err)
			continue
		}
		defer resp.Body.Close()

		cursor = follower.NextCursor

		for _, u := range follower.Users {
			profiles = append(profiles, t.toTwitterProfile(u))
		}
	}

	return profiles
}

func (t *Twitter) toTwitterProfile(user twitter.User) *TwitterProfile {
	joinedAt, err := time.Parse(time.RubyDate, user.CreatedAt)
	if err != nil {
		fmt.Printf("failed to parse joinedAt(%v) time for %q: %s", user.CreatedAt, user.ScreenName, err)
	}

	jsonEntites, err := json.Marshal(user.Entities)
	if err != nil {
		fmt.Printf("failed to marshal user.Entities(%v) to JSON for %q: %s", user.Entities, user.ScreenName, err)
	}

	var ent map[string]interface{}

	err = json.Unmarshal(jsonEntites, &ent)
	if err != nil {
		fmt.Printf("failed to unmarshall user.Entities(%v) from JSON for %q: %s", ent, user.ScreenName, err)
	}

	return &TwitterProfile{
		TwitterID:           user.IDStr,
		Name:                user.Name,
		Username:            user.ScreenName,
		Location:            user.Location,
		Bio:                 user.Description,
		URL:                 user.URL,
		Email:               user.Email,
		ProfileBannerURL:    user.ProfileBannerURL,
		ProfileImageURL:     user.ProfileImageURLHttps,
		Verified:            user.Verified,
		Protected:           user.Protected,
		DefaultProfile:      user.DefaultProfile,
		DefaultProfileImage: user.DefaultProfileImage,
		FollowersCount:      user.FollowersCount,
		FollowingsCount:     user.FriendsCount,
		FavouritesCount:     user.FavouritesCount,
		ListedCount:         user.ListedCount,
		TweetsCount:         user.StatusesCount,
		Entities:            ent,
		JoinedAt:            joinedAt,
	}
}

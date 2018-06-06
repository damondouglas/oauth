package storage

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	api "google.golang.org/api/oauth2/v2"
)

// Storage abstracts User storage and retrieval by userID, email, authorization code or refresh token.
type Storage interface {
	Save(context.Context, *User) error
	GetByRefreshToken(context.Context, string) (*User, error)
	GetByIDToken(context.Context, string) (*User, error)
}

// User holds data about the OAuth2 user.
type User struct {
	Email string
	Token *oauth2.Token
}

// GetUserFromToken gets User from call to profile user info using oauth2.Config.
func GetUserFromToken(ctx context.Context, config *oauth2.Config, tok *oauth2.Token) (*User, error) {
	u := new(User)
	client := config.Client(ctx, tok)
	srv, err := api.New(client)
	if err != nil {
		return nil, err
	}
	info, err := srv.Tokeninfo().AccessToken(tok.AccessToken).Do()
	if err != nil {
		return nil, err
	}

	u.Email = info.Email
	u.Token = tok

	return u, nil
}

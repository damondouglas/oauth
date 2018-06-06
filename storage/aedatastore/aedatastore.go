package aedatastore

import (
	"errors"
	"time"

	"github.com/damondouglas/oauth/storage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/datastore"
)

const userKindToken = "User"
const keyOnlyFilterToken = "__key__ >"

// Datastore represents appengine datastore interface
type Datastore struct {
}

type datastoreUser struct {
	Email             string
	ID                string
	AuthorizationCode string
	AccessToken       string
	RefreshToken      string
	Expiry            time.Time
}

func from(user *storage.User) *datastoreUser {
	du := new(datastoreUser)
	du.Email = user.Email
	du.ID = user.ID
	du.AuthorizationCode = user.AuthorizationCode
	if user.Token != nil {
		du.AccessToken = user.Token.AccessToken
		du.RefreshToken = user.Token.RefreshToken
		du.Expiry = user.Token.Expiry
	}

	return du
}

func (du *datastoreUser) export() *storage.User {
	u := new(storage.User)
	u.Email = du.Email
	u.ID = du.ID
	u.AuthorizationCode = du.AuthorizationCode
	u.Token = new(oauth2.Token)
	u.Token.RefreshToken = du.RefreshToken
	u.Token.AccessToken = du.AccessToken
	u.Token.Expiry = du.Expiry

	return u
}

func (d *Datastore) newKey(ctx context.Context, u *storage.User) *datastore.Key {
	return datastore.NewKey(ctx, userKindToken, u.Email, 0, nil)
}

func (d *Datastore) query() *datastore.Query {
	return datastore.NewQuery(userKindToken)
}

func (d *Datastore) getByKey(ctx context.Context, key *datastore.Key) (*storage.User, error) {
	var u *storage.User
	err := datastore.Get(ctx, key, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *Datastore) exists(ctx context.Context, u *storage.User) bool {
	k := d.newKey(ctx, u)
	u, err := d.getByKey(ctx, k)
	if err != nil {
		return false
	}
	return u != nil
}

// Save puts storage.User data into appengine datastore.
func (d *Datastore) Save(ctx context.Context, u *storage.User) error {
	key := d.newKey(ctx, u)
	du := from(u)
	storedKey, err := datastore.Put(ctx, key, du)
	if err != nil {
		return err
	}
	if !key.Equal(storedKey) {
		return errors.New("key does not equal storedKey")
	}
	return nil
}

// GetByEmail retreives storage.User from email address.
func (d *Datastore) GetByEmail(ctx context.Context, email string) (*storage.User, error) {
	var du []datastoreUser
	u := &storage.User{Email: email}
	k := d.newKey(ctx, u)
	q := datastore.NewQuery(userKindToken).Filter(keyOnlyFilterToken, k)
	if _, err := q.GetAll(ctx, &du); err != nil {
		return nil, err
	}
	if len(du) == 0 {
		return nil, errors.New("query by email resulted empty")
	}
	return du[0].export(), nil
}

func (d *Datastore) GetByRefreshToken(ctx context.Context, refreshToken string) (*storage.User, error) {
	var du []datastoreUser
	q := datastore.NewQuery(userKindToken).Filter("RefreshToken =", refreshToken)
	if _, err := q.GetAll(ctx, &du); err != nil {
		return nil, err
	}
	if len(du) == 0 {
		return nil, errors.New("query by refresh token resulted empty")
	}
	return du[0].export(), nil
}

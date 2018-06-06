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
const iDTokenKey = "id_token"

// Datastore represents appengine datastore interface
type Datastore struct {
}

type datastoreUser struct {
	Email        string
	IDToken      string
	IDTokenPartA string
	IDTokenPartB string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

func splitToken(tokenValue string) []string {
	n := len(tokenValue)
	var a int
	if n%2 == 1 {
		a = (n - 1) / 2
	} else {
		a = n / 2
	}
	return []string{tokenValue[:a], tokenValue[a:]}
}

func from(user *storage.User) *datastoreUser {
	du := new(datastoreUser)
	du.Email = user.Email

	if user.Token != nil {
		du.AccessToken = user.Token.AccessToken
		du.RefreshToken = user.Token.RefreshToken
		du.Expiry = user.Token.Expiry
		temp := user.Token.Extra(iDTokenKey)
		IDToken, ok := temp.(string)
		if ok {
			du.IDToken = IDToken
			tokens := splitToken(IDToken)
			du.IDTokenPartA = tokens[0]
			du.IDTokenPartB = tokens[1]
		}
	}

	return du
}

func (du *datastoreUser) export() *storage.User {
	u := new(storage.User)
	u.Email = du.Email
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

// GetByRefreshToken loads storage.User from refreshToken.
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

// GetByIDToken loads storage.User from IDToken.
func (d *Datastore) GetByIDToken(ctx context.Context, IDToken string) (*storage.User, error) {
	var du []datastoreUser
	tokens := splitToken(IDToken)
	q := datastore.NewQuery(userKindToken).Filter("IDTokenPartA =", tokens[0]).Filter("IDTokenPartB =", tokens[1])
	if _, err := q.GetAll(ctx, &du); err != nil {
		return nil, err
	}
	if len(du) == 0 {
		return nil, errors.New("query by IDtoken resulted empty")
	}
	return du[0].export(), nil

}

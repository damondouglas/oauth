package mock

import (
	"os/user"
	"path"
	"testing"
	"time"

	"github.com/damondouglas/assert"
	"github.com/damondouglas/oauth/storage"
	"golang.org/x/oauth2"
)

func userHelper() (*storage.User, error) {
	u := new(storage.User)
	u.ID = "foo"
	u.Email = "foo@example.com"
	u.AuthorizationCode = "foo_authcode"
	token := new(oauth2.Token)
	token.AccessToken = "foo_accesstoken"
	exp := time.Now()
	d, err := time.ParseDuration("1000h")
	if err != nil {
		return nil, err
	}
	token.Expiry = exp.Add(d)
	token.RefreshToken = "foo_refreshtoken"
	u.Token = token

	return u, nil
}

func TestSave(t *testing.T) {
	a := assert.New(t)
	usr, err := user.Current()
	a.HandleError(err)
	root := path.Join(usr.HomeDir, "oauthmock")
	u, err := userHelper()
	a.HandleError(err)
	m, err := New(root)
	a.HandleError(err)
	err = m.Save(u)
	a.HandleError(err)
}

func TestGetByUserID(t *testing.T) {
	a := assert.New(t)
	usr, err := user.Current()
	a.HandleError(err)
	root := path.Join(usr.HomeDir, "oauthmock")
	m, err := New(root)
	a.HandleError(err)

	u, err := m.GetByUserID("foo")
	a.Equals(u.Email, "foo@example.com")
	a.Equals(u.AuthorizationCode, "foo_authcode")
	a.Equals(u.ID, "foo")
	a.Equals(u.Token.RefreshToken, "foo_refreshtoken")
	a.Equals(u.Token.AccessToken, "foo_accesstoken")
}

func TestGetByEmail(t *testing.T) {
	a := assert.New(t)
	usr, err := user.Current()
	a.HandleError(err)
	root := path.Join(usr.HomeDir, "oauthmock")
	m, err := New(root)
	a.HandleError(err)
	u, err := m.GetByEmail("foo@example.com")
	a.HandleError(err)

	a.Equals(u.Email, "foo@example.com")
	a.Equals(u.AuthorizationCode, "foo_authcode")
	a.Equals(u.ID, "foo")
	a.Equals(u.Token.RefreshToken, "foo_refreshtoken")
	a.Equals(u.Token.AccessToken, "foo_accesstoken")
}

func TestGetByAuthorizationCode(t *testing.T) {
	a := assert.New(t)
	usr, err := user.Current()
	a.HandleError(err)
	root := path.Join(usr.HomeDir, "oauthmock")
	m, err := New(root)
	a.HandleError(err)
	u, err := m.GetByAuthorizationCode("foo_authcode")
	a.HandleError(err)

	a.Equals(u.Email, "foo@example.com")
	a.Equals(u.AuthorizationCode, "foo_authcode")
	a.Equals(u.ID, "foo")
	a.Equals(u.Token.RefreshToken, "foo_refreshtoken")
	a.Equals(u.Token.AccessToken, "foo_accesstoken")
}

func TestGetByRefreshToken(t *testing.T) {
	a := assert.New(t)
	usr, err := user.Current()
	a.HandleError(err)
	root := path.Join(usr.HomeDir, "oauthmock")
	m, err := New(root)
	a.HandleError(err)
	u, err := m.GetByRefreshToken("foo_refreshtoken")
	a.HandleError(err)

	a.Equals(u.Email, "foo@example.com")
	a.Equals(u.AuthorizationCode, "foo_authcode")
	a.Equals(u.ID, "foo")
	a.Equals(u.Token.RefreshToken, "foo_refreshtoken")
	a.Equals(u.Token.AccessToken, "foo_accesstoken")
}

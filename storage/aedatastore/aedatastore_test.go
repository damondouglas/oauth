// TODO: http://localhost:59717: dial tcp [::1]:59717: getsockopt: connection refused
package aedatastore

import (
	"log"
	"testing"
	"time"

	"github.com/damondouglas/assert"
	"github.com/damondouglas/oauth/storage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/aetest"
)

var (
	ctx context.Context
	d   *Datastore
)

func init() {
	var err error
	ctx, err = appengineContext()
	if err != nil {
		log.Panic(err)
	}
	d = New(ctx, "Auth")
}

func appengineContext() (context.Context, error) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		return nil, err
	}
	defer done()
	return ctx, nil
}

func TestNewKey(t *testing.T) {
	a := assert.New(t)
	key := d.newKey(&storage.User{Email: "foo@example.com"})
	a.Equals(key.StringID(), "foo@example.com")
}

func TestExists(t *testing.T) {
	a := assert.New(t)
	u := &storage.User{Email: "notinstore@example.com"}
	e := d.exists(u)
	a.Equals(e, false)
}

func TestSave(t *testing.T) {
	a := assert.New(t)

	u := &storage.User{Email: "foo@example.com"}
	u.Token = new(oauth2.Token)
	u.Token.RefreshToken = "foo_refreshtoken"
	u.Token.AccessToken = "foo_accesstoken"
	u.Token.Expiry = time.Now()

	err := d.Save(u)
	a.HandleError(err)
}

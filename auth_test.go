package oauth

import (
	"net/url"
	"strings"
	"testing"

	"github.com/damondouglas/assert"
)

func TestBuildRedirectURL(t *testing.T) {
	var auth *Auth
	var redirectURL *RedirectURL
	var err error
	var u *url.URL

	a := assert.New(t)
	auth, err = New("./mock/client_secret.json", "http://localhost:8080/oauthcallback", nil)
	a.Equals(err, nil)
	a.NotEquals(auth, nil)

	redirectURL = auth.BuildRedirectURL(false, false)
	a.NotEquals(redirectURL.SessionID, nil)
	a.Equals(strings.HasPrefix(redirectURL.Value, "https://accounts.google.com/o/oauth2"), true)
	u, err = url.Parse(redirectURL.Value)
	a.Equals(err, nil)
	a.NotEquals(u, nil)
	q := u.Query()
	a.Equals(q.Get("client_id"), "some_client_id.apps.googleusercontent.com")
	a.Equals(q.Get("access_type"), "online")
	a.Equals(q.Get("approval_prompt"), "")
	a.Equals(q.Get("state"), redirectURL.SessionID)

	redirectURL = auth.BuildRedirectURL(true, false)
	t.Logf("value: %v", redirectURL)
	a.NotEquals(redirectURL.SessionID, nil)
	a.Equals(strings.HasPrefix(redirectURL.Value, "https://accounts.google.com/o/oauth2"), true)
	u, err = url.Parse(redirectURL.Value)
	a.Equals(err, nil)
	a.NotEquals(u, nil)
	q = u.Query()
	a.Equals(q.Get("client_id"), "some_client_id.apps.googleusercontent.com")
	a.Equals(q.Get("access_type"), "online")
	a.Equals(q.Get("approval_prompt"), "force")
	a.Equals(q.Get("state"), redirectURL.SessionID)

	redirectURL = auth.BuildRedirectURL(false, true)
	t.Logf("value: %v", redirectURL)
	a.NotEquals(redirectURL.SessionID, nil)
	a.Equals(strings.HasPrefix(redirectURL.Value, "https://accounts.google.com/o/oauth2"), true)
	u, err = url.Parse(redirectURL.Value)
	a.Equals(err, nil)
	a.NotEquals(u, nil)
	q = u.Query()
	a.Equals(q.Get("client_id"), "some_client_id.apps.googleusercontent.com")
	a.Equals(q.Get("access_type"), "offline")
	a.Equals(q.Get("approval_prompt"), "")
	a.Equals(q.Get("state"), redirectURL.SessionID)

	redirectURL = auth.BuildRedirectURL(true, true)
	t.Logf("value: %v", redirectURL)
	a.NotEquals(redirectURL.SessionID, nil)
	a.Equals(strings.HasPrefix(redirectURL.Value, "https://accounts.google.com/o/oauth2"), true)
	u, err = url.Parse(redirectURL.Value)
	a.Equals(err, nil)
	a.NotEquals(u, nil)
	q = u.Query()
	a.Equals(q.Get("client_id"), "some_client_id.apps.googleusercontent.com")
	a.Equals(q.Get("access_type"), "offline")
	a.Equals(q.Get("approval_prompt"), "force")
	a.Equals(q.Get("state"), redirectURL.SessionID)
}

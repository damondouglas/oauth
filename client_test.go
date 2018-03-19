package oauth

import (
	"testing"

	"github.com/damondouglas/assert"
)

func TestMergeScopesFromBase(t *testing.T) {
	a := assert.New(t)
	scopes := mergeScopesFromBase(nil)
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 2)

	scopes = mergeScopesFromBase([]string{})
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 2)

	scopes = mergeScopesFromBase([]string{"foo"})
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 3)
	a.Equals(scopes[0], "email")
	a.Equals(scopes[1], "foo")
	a.Equals(scopes[2], "profile")

	scopes = mergeScopesFromBase([]string{"email"})
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 2)
	a.Equals(scopes[0], "email")
	a.Equals(scopes[1], "profile")

	scopes = mergeScopesFromBase([]string{"profile"})
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 2)
	a.Equals(scopes[0], "email")
	a.Equals(scopes[1], "profile")

	scopes = mergeScopesFromBase([]string{"email", "profile", "foo"})
	a.NotEquals(scopes, nil)
	a.Equals(len(scopes), 3)
	a.Equals(scopes[0], "email")
	a.Equals(scopes[1], "foo")
	a.Equals(scopes[2], "profile")
}

func TestBuildOauthClientFromCredentialsPath(t *testing.T) {
	a := assert.New(t)
	path := "./mock/client_secret.json"
	redirectURL := "http://localhost:8080/oauth2callback"
	config, err := buildOauthClientFromCredentialsPath(path, redirectURL, nil)
	a.Equals(err, nil)
	a.NotEquals(config, nil)
	a.Equals(config.ClientID, "some_client_id.apps.googleusercontent.com")
	a.Equals(config.ClientSecret, "some_client_secret")
	a.Equals(config.Scopes[0], "email")
	a.Equals(config.Scopes[1], "profile")
}

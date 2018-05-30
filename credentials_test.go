package oauth

import (
	"testing"

	"github.com/damondouglas/assert"
)

func TestCredentialsFromJSON(t *testing.T) {
	a := assert.New(t)
	jsonData := []byte(`{"web": {"client_id": "foo", "client_secret": "bar"}}`)
	cred := credentialsFromJSON(jsonData)
	a.Equals(cred.ClientID, "foo")
	a.Equals(cred.ClientSecret, "bar")
}

func TestConfigureOauthFromFilePath(t *testing.T) {
	a := assert.New(t)
	cred, err := credentialsFromPath("./mock/client_secret.json")
	a.HandleError(err)
	a.Equals(cred.ClientID, "some_client_id.apps.googleusercontent.com")
	a.Equals(cred.ClientSecret, "some_client_secret")
}
